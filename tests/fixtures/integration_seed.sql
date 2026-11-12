-- Data for the three integration scenarios. IDs are mirrored in postgres_setup_test.go.
INSERT INTO Users (id, login, email, password, is_admin, sex, weight, age, is_pregnant, is_driver)
VALUES
  ('00000000-0000-4000-8000-000000000001', 'integration-cabinet-create', 'create@example.test', 'unused', false, false, 70, 30, false, false),
  ('00000000-0000-4000-8000-000000000002', 'integration-cabinet-owner', 'owner@example.test', 'unused', false, false, 70, 30, false, false),
  ('00000000-0000-4000-8000-000000000003', 'integration-cabinet-other', 'other@example.test', 'unused', false, false, 70, 30, false, false),
  ('00000000-0000-4000-8000-000000000004', 'integration-restock', 'restock@example.test', 'unused', false, false, 70, 30, false, false);

INSERT INTO Cabinets (id, user_id, name)
VALUES ('00000000-0000-4000-8000-000000000011', '00000000-0000-4000-8000-000000000002', 'Home');

INSERT INTO Form (id, name)
VALUES ('00000000-0000-4000-8000-000000000021', 'tablet');

INSERT INTO Unit (id, name)
VALUES ('00000000-0000-4000-8000-000000000022', 'piece');

INSERT INTO Medicine (id, form_id, unit_id, name, expire_time, effect_on_driver, effect_on_pregnant, method_of_application, is_prescription)
VALUES ('00000000-0000-4000-8000-000000000031', '00000000-0000-4000-8000-000000000021',
        '00000000-0000-4000-8000-000000000022', 'Test medicine', 12, false, false, 'oral', false);
