package lab3

import (
	"cmp"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"runtime"
	"slices"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
)

type Config struct {
	DSN, RunID, Profile, ResultsDir, CPU string
}

type benchmarkProfile struct {
	repeats, warmup, measurements int
}

type phaseResult struct {
	Operations int           `json:"operations"`
	Duration   time.Duration `json:"duration_ns"`
}

type seriesResult struct {
	Repeat   int                `json:"repeat"`
	Variant  string             `json:"variant"`
	Scenario string             `json:"scenario"`
	Warmup   phaseResult        `json:"warmup"`
	Measure  phaseResult        `json:"measure"`
	Teardown phaseResult        `json:"teardown"`
	Overlap  bool               `json:"overlap,omitempty"`
	Latency  []operationLatency `json:"latency,omitempty"`
	Error    string             `json:"error,omitempty"`
	started  time.Time
}

type summary struct {
	RunID           string         `json:"run_id"`
	Profile         string         `json:"profile"`
	CPU             string         `json:"cpu"`
	Go              string         `json:"go_version"`
	Postgres        string         `json:"postgres_version"`
	Started         time.Time      `json:"started_at"`
	Duration        time.Duration  `json:"duration_ns"`
	Overhead        time.Duration  `json:"timer_overhead_median_ns"`
	HarnessOverhead time.Duration  `json:"empty_loop_ns_per_op"`
	Series          []seriesResult `json:"series"`
	Error           string         `json:"error,omitempty"`
}

type sample struct {
	operation string
	phase     string
	start     time.Time
	elapsed   time.Duration
	err       string
}

func Run(ctx context.Context, config Config) (err error) {
	profile := benchmarkProfile{3, 50, 1000}
	switch config.Profile {
	case "full":
	case "smoke":
		profile = benchmarkProfile{1, 5, 20}
	default:
		return fmt.Errorf("unknown benchmark profile %q", config.Profile)
	}

	report := summary{
		RunID: config.RunID, Profile: config.Profile, CPU: config.CPU,
		Started: time.Now(), Go: runtime.Version(), Series: []seriesResult{},
	}
	workload := scenarios(profile.warmup + profile.measurements)
	results := make([][]sample, 0, profile.repeats*2*len(workload))
	defer func() {
		report.Duration = time.Since(report.Started)
		if err != nil {
			report.Error = err.Error()
		}
		err = errors.Join(err, saveResults(config.ResultsDir, report, results))
	}()

	seed, err := pgx.ParseConfig(config.DSN)
	if err != nil {
		return err
	}
	if seed.Database != "mypills_lab3" {
		return errors.New("benchmark requires the seeded mypills_lab3 database")
	}
	adminConfig := seed.Copy()
	adminConfig.Database = "postgres"
	admin := stdlib.OpenDB(*adminConfig)
	defer admin.Close()
	if err := admin.QueryRowContext(ctx, "SHOW server_version").Scan(&report.Postgres); err != nil {
		return err
	}

	if _, err := freshDatabase(ctx, admin, seed, func(pool *sql.DB) error {
		orm, err := NewGORMCatalog(pool)
		if err != nil {
			return err
		}
		return verifyCatalogs(ctx, NewSQLCatalog(pool), orm)
	}); err != nil {
		return fmt.Errorf("preflight: %w", err)
	}

	overheadStart := time.Now()
	overhead, err := measure(ctx, nil, "measure", operation{run: func(context.Context, catalog, int) error { return nil }}, 0, 10000)
	if err != nil {
		return err
	}
	report.HarnessOverhead = time.Since(overheadStart) / time.Duration(len(overhead))
	slices.SortFunc(overhead, func(a, b sample) int { return cmp.Compare(a.elapsed, b.elapsed) })
	report.Overhead = overhead[len(overhead)/2].elapsed

	variants := []string{"sql", "gorm"}
	for repeat := 1; repeat <= profile.repeats; repeat++ {
		if repeat > 1 {
			slices.Reverse(variants)
		}
		for _, variant := range variants {
			for _, scenario := range workload {
				result := seriesResult{Repeat: repeat, Variant: variant, Scenario: scenario.name}
				var samples []sample
				teardown, err := freshDatabase(ctx, admin, seed, func(pool *sql.DB) error {
					var c catalog = NewSQLCatalog(pool)
					if variant == "gorm" {
						orm, err := NewGORMCatalog(pool)
						if err != nil {
							return err
						}
						c = orm
					}
					var measureErr error
					samples, measureErr = measureSeries(ctx, c, scenario, profile, &result)
					return measureErr
				})
				result.Teardown.Duration = teardown
				if err != nil {
					result.Error = err.Error()
				}
				report.Series = append(report.Series, result)
				results = append(results, samples)
				if err != nil {
					return fmt.Errorf("repeat %d, %s %s: %w", repeat, variant, result.Scenario, err)
				}
				fmt.Printf("repeat %d: %s %s (%d measured)\n", repeat, variant, result.Scenario, result.Measure.Operations)
			}
		}
	}
	return nil
}

func freshDatabase(ctx context.Context, admin *sql.DB, seed *pgx.ConnConfig, run func(*sql.DB) error) (teardown time.Duration, err error) {
	if _, err := admin.ExecContext(ctx, "CREATE DATABASE mypills_lab3_series TEMPLATE mypills_lab3"); err != nil {
		return 0, err
	}
	config := seed.Copy()
	config.Database = "mypills_lab3_series"
	pool := stdlib.OpenDB(*config)
	pool.SetMaxOpenConns(2)
	defer func() {
		started := time.Now()
		closeErr := pool.Close()
		cleanup, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		_, dropErr := admin.ExecContext(cleanup, "DROP DATABASE mypills_lab3_series")
		err = errors.Join(err, closeErr, dropErr)
		teardown = time.Since(started)
	}()
	return 0, run(pool)
}

func measureSeries(ctx context.Context, c catalog, scenario scenario, profile benchmarkProfile, result *seriesResult) ([]sample, error) {
	result.started = time.Now()
	warmup, duration, err := measurePhase(ctx, c, scenario, "warmup", 0, profile.warmup)
	result.Warmup = phaseResult{len(warmup), duration}
	if err != nil {
		return warmup, err
	}
	measured, duration, err := measurePhase(ctx, c, scenario, "measure", profile.warmup, profile.measurements)
	result.Measure = phaseResult{len(measured), duration}
	if err == nil && len(scenario.operations) == 2 {
		result.Overlap = overlap(measured[:profile.measurements], measured[profile.measurements:])
		if !result.Overlap {
			err = errors.New("mixed read and write did not overlap")
		}
	}
	return append(warmup, measured...), err
}

func measurePhase(ctx context.Context, c catalog, scenario scenario, phase string, offset, count int) ([]sample, time.Duration, error) {
	start := time.Now()
	first := scenario.operations[0]
	if len(scenario.operations) == 1 {
		samples, err := measure(ctx, c, phase, first, offset, count)
		return samples, time.Since(start), err
	}
	var reads, writes []sample
	err := RunMixed(ctx, func(ctx context.Context) error {
		var err error
		reads, err = measure(ctx, c, phase, first, offset, count)
		return err
	}, func(ctx context.Context) error {
		var err error
		writes, err = measure(ctx, c, phase, scenario.operations[1], offset, count)
		return err
	})
	return append(reads, writes...), time.Since(start), err
}

func measure(ctx context.Context, c catalog, phase string, operation operation, offset, count int) ([]sample, error) {
	samples := make([]sample, 0, count)
	for i := range count {
		if err := ctx.Err(); err != nil {
			return samples, err
		}
		start := time.Now()
		err := operation.run(ctx, c, offset+i)
		entry := sample{operation: operation.name, phase: phase, start: start, elapsed: time.Since(start)}
		if err != nil {
			entry.err = err.Error()
			return append(samples, entry), err
		}
		samples = append(samples, entry)
	}
	return samples, nil
}

func overlap(reads, writes []sample) bool {
	for len(reads) > 0 && len(writes) > 0 {
		read, write := reads[0], writes[0]
		switch {
		case !read.start.Add(read.elapsed).After(write.start):
			reads = reads[1:]
		case !write.start.Add(write.elapsed).After(read.start):
			writes = writes[1:]
		default:
			return true
		}
	}
	return false
}
