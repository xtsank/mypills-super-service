-- Run on every WebLab#5 start: the primary volume may already contain WebLab#4 data.
SELECT format('CREATE ROLE %I LOGIN', 'lab5_reader')
WHERE NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'lab5_reader') \gexec
SELECT format('ALTER ROLE %I LOGIN NOSUPERUSER NOCREATEDB NOCREATEROLE NOREPLICATION PASSWORD %L',
              'lab5_reader', :'read_password') \gexec
ALTER ROLE lab5_reader SET default_transaction_read_only = on;

SELECT format('GRANT CONNECT ON DATABASE %I TO lab5_reader', :'db_name') \gexec
GRANT USAGE ON SCHEMA public TO lab5_reader;
REVOKE ALL PRIVILEGES ON ALL TABLES IN SCHEMA public FROM lab5_reader;
GRANT SELECT ON ALL TABLES IN SCHEMA public TO lab5_reader;
SELECT format('ALTER DEFAULT PRIVILEGES FOR ROLE %I IN SCHEMA public GRANT SELECT ON TABLES TO lab5_reader',
              :'db_user') \gexec

SELECT format('CREATE ROLE %I LOGIN REPLICATION', 'lab5_replication')
WHERE NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'lab5_replication') \gexec
SELECT format('ALTER ROLE %I LOGIN REPLICATION NOSUPERUSER NOCREATEDB NOCREATEROLE PASSWORD %L',
              'lab5_replication', :'replication_password') \gexec

-- A physical slot keeps WAL while the local replica is briefly unavailable.
SELECT 'SELECT pg_create_physical_replication_slot(''mypills_lab5'');'
WHERE NOT EXISTS (SELECT 1 FROM pg_replication_slots WHERE slot_name = 'mypills_lab5') \gexec

-- Exporters read PostgreSQL statistics with a dedicated account.
\getenv monitor_password LAB5_MONITOR_PASSWORD
SELECT 'CREATE ROLE lab5_monitor LOGIN'
WHERE NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'lab5_monitor') \gexec
SELECT format('ALTER ROLE lab5_monitor LOGIN NOSUPERUSER NOCREATEDB NOCREATEROLE NOREPLICATION PASSWORD %L',
              :'monitor_password') \gexec
GRANT pg_monitor TO lab5_monitor;
ALTER ROLE lab5_monitor SET default_transaction_read_only = on;
