#!/usr/bin/env bash
set -euo pipefail

LOGIN=${1:-demo_user}
EMAIL=${2:-demo_user@example.com}
PASSWORD=${3:-demo_password}

DB_USER=user
DB_PASSWORD=password
DB_NAME=mypills
DB_PORT=5432
DB_HOST=db

export PGPASSWORD="$DB_PASSWORD"

psql "host=$DB_HOST port=$DB_PORT user=$DB_USER dbname=$DB_NAME" <<SQL
\set ON_ERROR_STOP on
\set login '${LOGIN}'
\set email '${EMAIL}'
\set password '${PASSWORD}'

WITH upserted AS (
  INSERT INTO Users (id, login, email, password, is_admin, sex, weight, age, is_pregnant, is_driver, notify_enabled, notify_interval_minutes, last_notified_at)
  VALUES (gen_random_uuid(), :'login', :'email', :'password', false, false, 70, 30, false, false, true, 1, NULL)
  ON CONFLICT (login)
  DO UPDATE SET
    email = EXCLUDED.email,
    notify_enabled = true,
    notify_interval_minutes = 1,
    last_notified_at = NULL
  RETURNING id
),
form AS (
  INSERT INTO Form (id, name) VALUES (gen_random_uuid(), 'Таблетки') RETURNING id
),
unit AS (
  INSERT INTO Unit (id, name) VALUES (gen_random_uuid(), 'шт') RETURNING id
),
illness AS (
  INSERT INTO Illness (id, name) VALUES (gen_random_uuid(), 'Простуда') RETURNING id
),
medicine AS (
  INSERT INTO Medicine (id, form_id, unit_id, name, expire_time, effect_on_driver, effect_on_pregnant, method_of_application, is_prescription)
  VALUES (gen_random_uuid(), (SELECT id FROM form), (SELECT id FROM unit), 'Тестамин', 1, false, false, 'перорально', false)
  RETURNING id
),
rec AS (
  INSERT INTO Recommendations (medicine_id, illness_id) VALUES ((SELECT id FROM medicine), (SELECT id FROM illness))
),
item AS (
  INSERT INTO User_Medicine (id, user_id, medicine_id, date_of_manufacture, quantity)
  VALUES (gen_random_uuid(), (SELECT id FROM upserted), (SELECT id FROM medicine), CURRENT_DATE - INTERVAL '2 months', 10)
)
SELECT (SELECT id FROM upserted) AS user_id, (SELECT id FROM medicine) AS medicine_id;
SQL

echo "Seed completed for login=$LOGIN (email=$EMAIL)"

