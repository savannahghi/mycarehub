BEGIN;

ALTER TABLE
    IF EXISTS "common_usersurveys"
    ADD COLUMN IF NOT EXISTS "genders" varchar(64) [],
    ADD COLUMN IF NOT EXISTS "age_range" jsonb,
    ADD COLUMN IF NOT EXISTS "client_types" varchar(64) [],
    ADD COLUMN IF NOT EXISTS "set_id" uuid;

COMMIT;