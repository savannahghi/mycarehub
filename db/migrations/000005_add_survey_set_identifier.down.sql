BEGIN;

ALTER TABLE
    IF EXISTS "common_usersurveys"
    DROP COLUMN IF EXISTS "set_id";

COMMIT;