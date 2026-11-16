\getenv run_id LAB3_RUN_ID
\set summary `cat /input/summary.json`

BEGIN;
SELECT pg_advisory_xact_lock(hashtextextended(:'run_id', 0));

CREATE TEMP TABLE imported_samples (LIKE lab3_samples INCLUDING CONSTRAINTS);
ALTER TABLE imported_samples ADD CHECK (run_id = :'run_id');
\copy imported_samples FROM '/input/samples.csv' WITH (FORMAT csv, HEADER MATCH)

CREATE TEMP TABLE imported_run (LIKE lab3_runs INCLUDING CONSTRAINTS);
INSERT INTO imported_run VALUES (:'run_id', :'summary'::jsonb);

DO $$
BEGIN
    IF EXISTS (
        WITH expected AS (
            SELECT (series->>'repeat')::int AS repeat, series->>'variant' AS variant,
                   series->>'scenario' AS scenario,
                   (series->'warmup'->>'operations')::bigint AS warmup,
                   (series->'measure'->>'operations')::bigint AS measure
            FROM imported_run
            CROSS JOIN LATERAL jsonb_array_elements(summary->'series') AS series
        ), actual AS (
            SELECT repeat, variant, scenario,
                   count(*) FILTER (WHERE phase = 'warmup') AS warmup,
                   count(*) FILTER (WHERE phase = 'measure') AS measure
            FROM imported_samples
            GROUP BY repeat, variant, scenario
        )
        SELECT 1 FROM expected
        FULL JOIN actual USING (repeat, variant, scenario)
        WHERE coalesce(expected.warmup, 0) <> coalesce(actual.warmup, 0)
           OR coalesce(expected.measure, 0) <> coalesce(actual.measure, 0)
    ) THEN
        RAISE EXCEPTION 'operation counts in samples.csv and summary.json differ';
    END IF;
END $$;

DELETE FROM lab3_runs WHERE run_id = :'run_id';
INSERT INTO lab3_runs SELECT * FROM imported_run;
INSERT INTO lab3_samples SELECT * FROM imported_samples;
COMMIT;
