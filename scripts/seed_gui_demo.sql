-- Локальные демонстрационные данные для проверки сценариев GUI.
-- Повторный запуск не создаёт дубликаты.
BEGIN;

CREATE EXTENSION IF NOT EXISTS pgcrypto;

INSERT INTO Users (id, login, email, password, is_admin, sex, weight, age, is_pregnant, is_driver)
VALUES ('fd000001-0000-4000-8000-000000000001', 'demo_admin', 'demo@example.test',
        crypt('demo12345', gen_salt('bf')), true, false, 70, 30, false, false)
ON CONFLICT DO NOTHING;

INSERT INTO Form (id, name)
VALUES ('fd000001-0000-4000-8000-000000000002', 'Таблетки')
ON CONFLICT DO NOTHING;

INSERT INTO Unit (id, name)
VALUES ('fd000001-0000-4000-8000-000000000003', 'табл.')
ON CONFLICT DO NOTHING;

INSERT INTO Illness (id, name)
VALUES ('fd000001-0000-4000-8000-000000000004', 'Головная боль')
ON CONFLICT DO NOTHING;

INSERT INTO Substance (id, name)
VALUES ('fd000001-0000-4000-8000-000000000005', 'Парацетамол'),
       ('fd000001-0000-4000-8000-000000000006', 'Амоксициллин')
ON CONFLICT DO NOTHING;

INSERT INTO User_Illness (user_id, illness_id)
SELECT 'fd000001-0000-4000-8000-000000000001', 'fd000001-0000-4000-8000-000000000004'
WHERE NOT EXISTS (
    SELECT 1 FROM User_Illness
    WHERE user_id = 'fd000001-0000-4000-8000-000000000001'
      AND illness_id = 'fd000001-0000-4000-8000-000000000004'
);

INSERT INTO User_Substance (user_id, substance_id)
SELECT 'fd000001-0000-4000-8000-000000000001', 'fd000001-0000-4000-8000-000000000006'
WHERE NOT EXISTS (
    SELECT 1 FROM User_Substance
    WHERE user_id = 'fd000001-0000-4000-8000-000000000001'
      AND substance_id = 'fd000001-0000-4000-8000-000000000006'
);

INSERT INTO Medicine (id, form_id, unit_id, name, expire_time,
                      effect_on_driver, effect_on_pregnant, method_of_application, is_prescription)
VALUES ('fd000001-0000-4000-8000-000000000007',
        'fd000001-0000-4000-8000-000000000002',
        'fd000001-0000-4000-8000-000000000003',
        'Парацетамол (демо)', 24, false, false, 'Внутрь', false)
ON CONFLICT DO NOTHING;

INSERT INTO Medicine_Substance (id, medicine_id, substance_id, concentration)
VALUES ('fd000001-0000-4000-8000-000000000008',
        'fd000001-0000-4000-8000-000000000007',
        'fd000001-0000-4000-8000-000000000005', 500)
ON CONFLICT DO NOTHING;

INSERT INTO Dosage (id, medicine_id, value_from, value_to, dosage_type,
                    dosage_value, number_of_doses_per_day)
VALUES ('fd000001-0000-4000-8000-000000000009',
        'fd000001-0000-4000-8000-000000000007', 18, 120, 'age', 1, 2)
ON CONFLICT DO NOTHING;

INSERT INTO Recommendations (id, medicine_id, illness_id)
VALUES ('fd000001-0000-4000-8000-00000000000a',
        'fd000001-0000-4000-8000-000000000007',
        'fd000001-0000-4000-8000-000000000004')
ON CONFLICT DO NOTHING;

INSERT INTO Cabinets (id, user_id, name)
VALUES ('fd000001-0000-4000-8000-00000000000b',
        'fd000001-0000-4000-8000-000000000001', 'Домашняя (демо)')
ON CONFLICT DO NOTHING;

INSERT INTO User_Medicine (id, user_id, cabinet_id, medicine_id,
                           date_of_manufacture, quantity, minimum_quantity)
VALUES ('fd000001-0000-4000-8000-00000000000c',
        'fd000001-0000-4000-8000-000000000001',
        'fd000001-0000-4000-8000-00000000000b',
        'fd000001-0000-4000-8000-000000000007', CURRENT_DATE, 2, 5),
       ('fd000001-0000-4000-8000-00000000000d',
        'fd000001-0000-4000-8000-000000000001',
        'fd000001-0000-4000-8000-00000000000b',
        'fd000001-0000-4000-8000-000000000007', CURRENT_DATE - INTERVAL '3 years', 3, 1)
ON CONFLICT DO NOTHING;

INSERT INTO Restock_Items (id, user_id, medicine_id, desired_quantity)
VALUES ('fd000001-0000-4000-8000-00000000000e',
        'fd000001-0000-4000-8000-000000000001',
        'fd000001-0000-4000-8000-000000000007', 10)
ON CONFLICT DO NOTHING;

-- Дополнительные значения всех справочников.
INSERT INTO Form (id, name)
VALUES ('fd000001-0000-4000-8000-00000000000f', 'Капсулы'),
       ('fd000001-0000-4000-8000-000000000010', 'Сироп'),
       ('fd000001-0000-4000-8000-000000000011', 'Спрей')
ON CONFLICT DO NOTHING;

INSERT INTO Unit (id, name)
VALUES ('fd000001-0000-4000-8000-000000000012', 'капс.'),
       ('fd000001-0000-4000-8000-000000000013', 'мл'),
       ('fd000001-0000-4000-8000-000000000014', 'впрыск.')
ON CONFLICT DO NOTHING;

INSERT INTO Illness (id, name)
VALUES ('fd000001-0000-4000-8000-000000000015', 'Повышенная температура'),
       ('fd000001-0000-4000-8000-000000000016', 'Аллергический ринит'),
       ('fd000001-0000-4000-8000-000000000017', 'Бактериальная инфекция'),
       ('fd000001-0000-4000-8000-000000000018', 'Гастрит')
ON CONFLICT DO NOTHING;

INSERT INTO Substance (id, name)
VALUES ('fd000001-0000-4000-8000-000000000019', 'Ибупрофен'),
       ('fd000001-0000-4000-8000-00000000001a', 'Лоратадин'),
       ('fd000001-0000-4000-8000-00000000001b', 'Натрия хлорид')
ON CONFLICT DO NOTHING;

INSERT INTO User_Illness (user_id, illness_id)
SELECT 'fd000001-0000-4000-8000-000000000001', 'fd000001-0000-4000-8000-000000000018'
WHERE NOT EXISTS (
    SELECT 1 FROM User_Illness
    WHERE user_id = 'fd000001-0000-4000-8000-000000000001'
      AND illness_id = 'fd000001-0000-4000-8000-000000000018'
);

INSERT INTO Medicine (id, form_id, unit_id, name, expire_time,
                      effect_on_driver, effect_on_pregnant, method_of_application, is_prescription)
VALUES ('fd000001-0000-4000-8000-00000000001c',
        'fd000001-0000-4000-8000-000000000002',
        'fd000001-0000-4000-8000-000000000003',
        'Ибупрофен (демо)', 24, false, false, 'Внутрь', false),
       ('fd000001-0000-4000-8000-00000000001d',
        'fd000001-0000-4000-8000-00000000000f',
        'fd000001-0000-4000-8000-000000000012',
        'Лоратадин (демо)', 36, false, false, 'Внутрь', false),
       ('fd000001-0000-4000-8000-00000000001e',
        'fd000001-0000-4000-8000-00000000000f',
        'fd000001-0000-4000-8000-000000000012',
        'Амоксициллин (демо)', 24, false, false, 'Внутрь', true),
       ('fd000001-0000-4000-8000-00000000001f',
        'fd000001-0000-4000-8000-000000000011',
        'fd000001-0000-4000-8000-000000000014',
        'Солевой спрей (демо)', 24, false, false, 'Назально', false)
ON CONFLICT DO NOTHING;

INSERT INTO Medicine_Substance (id, medicine_id, substance_id, concentration)
VALUES ('fd000001-0000-4000-8000-000000000020',
        'fd000001-0000-4000-8000-00000000001c',
        'fd000001-0000-4000-8000-000000000019', 200),
       ('fd000001-0000-4000-8000-000000000021',
        'fd000001-0000-4000-8000-00000000001d',
        'fd000001-0000-4000-8000-00000000001a', 10),
       ('fd000001-0000-4000-8000-000000000022',
        'fd000001-0000-4000-8000-00000000001e',
        'fd000001-0000-4000-8000-000000000006', 500),
       ('fd000001-0000-4000-8000-000000000023',
        'fd000001-0000-4000-8000-00000000001f',
        'fd000001-0000-4000-8000-00000000001b', 9)
ON CONFLICT DO NOTHING;

INSERT INTO Dosage (id, medicine_id, value_from, value_to, dosage_type,
                    dosage_value, number_of_doses_per_day)
VALUES ('fd000001-0000-4000-8000-000000000024',
        'fd000001-0000-4000-8000-00000000001c', 18, 120, 'age', 1, 2),
       ('fd000001-0000-4000-8000-000000000025',
        'fd000001-0000-4000-8000-00000000001d', 18, 120, 'age', 1, 1),
       ('fd000001-0000-4000-8000-000000000026',
        'fd000001-0000-4000-8000-00000000001e', 18, 120, 'age', 1, 3),
       ('fd000001-0000-4000-8000-000000000027',
        'fd000001-0000-4000-8000-00000000001f', 18, 120, 'age', 1, 2)
ON CONFLICT DO NOTHING;

INSERT INTO Recommendations (id, medicine_id, illness_id)
VALUES ('fd000001-0000-4000-8000-000000000028',
        'fd000001-0000-4000-8000-000000000007',
        'fd000001-0000-4000-8000-000000000015'),
       ('fd000001-0000-4000-8000-000000000029',
        'fd000001-0000-4000-8000-00000000001c',
        'fd000001-0000-4000-8000-000000000004'),
       ('fd000001-0000-4000-8000-00000000002a',
        'fd000001-0000-4000-8000-00000000001c',
        'fd000001-0000-4000-8000-000000000015'),
       ('fd000001-0000-4000-8000-00000000002b',
        'fd000001-0000-4000-8000-00000000001d',
        'fd000001-0000-4000-8000-000000000016'),
       ('fd000001-0000-4000-8000-00000000002c',
        'fd000001-0000-4000-8000-00000000001e',
        'fd000001-0000-4000-8000-000000000017'),
       ('fd000001-0000-4000-8000-00000000002d',
        'fd000001-0000-4000-8000-00000000001f',
        'fd000001-0000-4000-8000-000000000016')
ON CONFLICT DO NOTHING;

INSERT INTO Contraindications (id, medicine_id, illness_id)
VALUES ('fd000001-0000-4000-8000-00000000002e',
        'fd000001-0000-4000-8000-00000000001c',
        'fd000001-0000-4000-8000-000000000018')
ON CONFLICT DO NOTHING;

INSERT INTO Cabinets (id, user_id, name)
VALUES ('fd000001-0000-4000-8000-00000000002f',
        'fd000001-0000-4000-8000-000000000001', 'В поездку (демо)')
ON CONFLICT DO NOTHING;

INSERT INTO User_Medicine (id, user_id, cabinet_id, medicine_id,
                           date_of_manufacture, quantity, minimum_quantity)
VALUES ('fd000001-0000-4000-8000-000000000030',
        'fd000001-0000-4000-8000-000000000001',
        'fd000001-0000-4000-8000-00000000000b',
        'fd000001-0000-4000-8000-00000000001c', CURRENT_DATE, 8, 2),
       ('fd000001-0000-4000-8000-000000000031',
        'fd000001-0000-4000-8000-000000000001',
        'fd000001-0000-4000-8000-00000000002f',
        'fd000001-0000-4000-8000-00000000001d', CURRENT_DATE, 4, 2),
       ('fd000001-0000-4000-8000-000000000032',
        'fd000001-0000-4000-8000-000000000001',
        'fd000001-0000-4000-8000-00000000002f',
        'fd000001-0000-4000-8000-00000000001e', CURRENT_DATE, 10, 2)
ON CONFLICT DO NOTHING;

INSERT INTO Restock_Items (id, user_id, medicine_id, desired_quantity)
VALUES ('fd000001-0000-4000-8000-000000000033',
        'fd000001-0000-4000-8000-000000000001',
        'fd000001-0000-4000-8000-00000000001d', 8)
ON CONFLICT DO NOTHING;

COMMIT;
