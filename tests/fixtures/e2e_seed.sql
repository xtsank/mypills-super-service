-- Only catalog data is seeded; the end-to-end test creates its user and stock via HTTP.
INSERT INTO Form (id, name)
VALUES ('00000000-0000-4000-8000-000000000021', 'tablet');

INSERT INTO Unit (id, name)
VALUES ('00000000-0000-4000-8000-000000000022', 'piece');

INSERT INTO Illness (id, name)
VALUES ('00000000-0000-4000-8000-000000000041', 'Test headache');

INSERT INTO Medicine (id, form_id, unit_id, name, expire_time, effect_on_driver, effect_on_pregnant, method_of_application, is_prescription)
VALUES
  ('00000000-0000-4000-8000-000000000042', '00000000-0000-4000-8000-000000000021', '00000000-0000-4000-8000-000000000022', 'Driver-safe medicine', 24, false, false, 'oral', false),
  ('00000000-0000-4000-8000-000000000043', '00000000-0000-4000-8000-000000000021', '00000000-0000-4000-8000-000000000022', 'Driver-unsafe medicine', 24, true, false, 'oral', false);

INSERT INTO Recommendations (medicine_id, illness_id)
VALUES
  ('00000000-0000-4000-8000-000000000042', '00000000-0000-4000-8000-000000000041'),
  ('00000000-0000-4000-8000-000000000043', '00000000-0000-4000-8000-000000000041');

INSERT INTO Dosage (id, medicine_id, value_from, value_to, dosage_type, dosage_value, number_of_doses_per_day)
VALUES
  ('00000000-0000-4000-8000-000000000044', '00000000-0000-4000-8000-000000000042', 18, 65, 'age', 1, 2),
  ('00000000-0000-4000-8000-000000000045', '00000000-0000-4000-8000-000000000043', 18, 65, 'age', 1, 2);
