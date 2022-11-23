BEGIN;

ALTER TABLE IF EXISTS "common_organisation"
  ADD COLUMN IF NOT EXISTS code integer;

ALTER TABLE IF EXISTS "common_organisation"
  ADD CONSTRAINT common_organisation_code_unique UNIQUE (code);
  
ALTER TABLE IF EXISTS "common_organisation"
  ALTER COLUMN code  SET DEFAULT 0;

ALTER TABLE IF EXISTS "common_organisation"
  ALTER COLUMN code  SET NOT NULL;

COMMIT;