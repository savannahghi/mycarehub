BEGIN;

ALTER TABLE IF EXISTS "common_organisation"
  DROP COLUMN IF EXISTS code;

COMMIT;