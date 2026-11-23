-- Move old positions into their owners' home cabinets in a single transaction.
-- No rows change on a repeat run unless new unassigned legacy data was imported.
BEGIN;

INSERT INTO cabinets (id, user_id, name)
SELECT gen_random_uuid(), owners.user_id, 'Домашняя аптечка'
FROM (SELECT DISTINCT user_id FROM user_medicine WHERE cabinet_id IS NULL) AS owners
ON CONFLICT (user_id, name) DO NOTHING;

UPDATE user_medicine AS item
SET cabinet_id = cabinet.id
FROM cabinets AS cabinet
WHERE item.cabinet_id IS NULL
  AND cabinet.user_id = item.user_id
  AND cabinet.name = 'Домашняя аптечка';

COMMIT;
