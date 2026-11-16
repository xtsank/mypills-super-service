package lab3

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"time"
)

type operationLatency struct {
	Operation string        `json:"operation"`
	Samples   int           `json:"samples"`
	P50       time.Duration `json:"p50_ns"`
	P75       time.Duration `json:"p75_ns"`
	P90       time.Duration `json:"p90_ns"`
	P95       time.Duration `json:"p95_ns"`
	P99       time.Duration `json:"p99_ns"`
}

func saveResults(dir string, report summary, samples [][]sample) error {
	for i, batch := range samples {
		report.Series[i].Latency = latencyByOperation(batch)
	}
	return errors.Join(
		writeSamples(filepath.Join(dir, "samples.csv"), report, samples),
		writeSummary(filepath.Join(dir, "summary.json"), report),
	)
}

func writeSummary(path string, report summary) error {
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(data, '\n'), 0644)
}

func writeSamples(path string, report summary, batches [][]sample) (err error) {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() { err = errors.Join(err, file.Close()) }()
	writer := csv.NewWriter(file)
	if err := writer.Write([]string{
		"run_id", "repeat", "variant", "scenario", "operation", "phase",
		"started_at", "offset_ns", "duration_ns", "status", "error",
	}); err != nil {
		return err
	}

	for i, batch := range batches {
		series := report.Series[i]
		for _, sample := range batch {
			status := "ok"
			if sample.err != "" {
				status = "error"
			}
			if err := writer.Write([]string{
				report.RunID, strconv.Itoa(series.Repeat), series.Variant, series.Scenario, sample.operation, sample.phase,
				sample.start.UTC().Format(time.RFC3339Nano),
				strconv.FormatInt(sample.start.Sub(series.started).Nanoseconds(), 10),
				strconv.FormatInt(sample.elapsed.Nanoseconds(), 10), status, sample.err,
			}); err != nil {
				return err
			}
		}
	}
	writer.Flush()
	return writer.Error()
}

func latencyByOperation(samples []sample) []operationLatency {
	durations := make(map[string][]time.Duration)
	for _, sample := range samples {
		if sample.phase == "measure" && sample.err == "" {
			durations[sample.operation] = append(durations[sample.operation], sample.elapsed)
		}
	}
	result := make([]operationLatency, 0, len(durations))
	for _, operation := range slices.Sorted(maps.Keys(durations)) {
		values := durations[operation]
		slices.Sort(values)
		at := func(percent int) time.Duration {
			// Nearest rank is ceil(N * percentile), using one-based ranks.
			return values[(len(values)*percent+99)/100-1]
		}
		result = append(result, operationLatency{
			Operation: operation, Samples: len(values),
			P50: at(50), P75: at(75), P90: at(90), P95: at(95), P99: at(99),
		})
	}
	return result
}
