BEGIN;

ALTER TABLE "common_organisation"
  ADD COLUMN code integer UNIQUE NOT NULL DEFAULT 0;

COMMIT;