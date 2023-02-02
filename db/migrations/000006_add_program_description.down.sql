BEGIN;

ALTER TABLE
    IF EXISTS "common_program"
    DROP COLUMN IF EXISTS "description";

COMMIT;