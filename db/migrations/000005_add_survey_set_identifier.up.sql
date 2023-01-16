BEGIN;

ALTER TABLE
    IF EXISTS "common_usersurveys"
    ADD COLUMN IF NOT EXISTS "set_id" uuid;

COMMIT;