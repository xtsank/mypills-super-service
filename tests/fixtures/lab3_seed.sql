BEGIN;

INSERT INTO Form (id, name)
VALUES ('10000000-0000-4000-8000-000000000001', 'Таблетка');

INSERT INTO Unit (id, name)
VALUES ('10000000-0000-4000-8000-000000000002', 'Штука');

INSERT INTO Medicine (
    id, form_id, unit_id, name, expire_time, effect_on_driver,
    effect_on_pregnant, method_of_application, is_prescription
)
SELECT
    ('20000000-0000-4000-8000-' || lpad(n::text, 12, '0'))::uuid,
    '10000000-0000-4000-8000-000000000001',
    '10000000-0000-4000-8000-000000000002',
    CASE
        WHEN n = 1 THEN 'Демо % концентрат'
        WHEN n = 2 THEN 'Демо _ раствор'
        ELSE
            (ARRAY['Демо Альфа', 'Демо Бета', 'Демо Гамма', 'Демо Дельта'])[1 + n % 4]
            || ' ' || lpad(n::text, 5, '0')
    END,
    12 + n % 25,
    n % 5 = 0,
    n % 3 = 0,
    'Тестовое описание',
    n % 7 = 0
FROM generate_series(1, 10000) AS series(n);

COMMIT;

ANALYZE Medicine;
