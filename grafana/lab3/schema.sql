CREATE TABLE lab3_runs (
    run_id text PRIMARY KEY,
    summary jsonb NOT NULL CHECK ((summary ->> 'run_id') IS NOT DISTINCT FROM run_id)
);

CREATE TABLE lab3_samples (
    run_id text NOT NULL REFERENCES lab3_runs(run_id) ON DELETE CASCADE,
    repeat integer NOT NULL CHECK (repeat > 0),
    variant text NOT NULL CHECK (variant IN ('sql', 'gorm')),
    scenario text NOT NULL,
    operation text NOT NULL,
    phase text NOT NULL CHECK (phase IN ('warmup', 'measure')),
    started_at timestamptz NOT NULL,
    offset_ns bigint NOT NULL CHECK (offset_ns >= 0),
    duration_ns bigint NOT NULL CHECK (duration_ns >= 0),
    status text NOT NULL CHECK (status IN ('ok', 'error')),
    error text,
    PRIMARY KEY (run_id, scenario, repeat, variant, operation, phase, offset_ns)
);

CREATE ROLE lab3_grafana LOGIN PASSWORD 'lab3_grafana';
GRANT SELECT ON lab3_runs, lab3_samples TO lab3_grafana;
