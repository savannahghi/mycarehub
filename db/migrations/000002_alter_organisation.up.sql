BEGIN;

ALTER TABLE "common_organisation"
  DROP COLUMN code;

COMMIT;