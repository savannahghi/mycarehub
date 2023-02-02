BEGIN;

ALTER TABLE
    IF EXISTS "common_program"
    ADD COLUMN IF NOT EXISTS "description" text;

COMMIT;