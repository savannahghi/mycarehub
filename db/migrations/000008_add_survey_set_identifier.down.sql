BEGIN;

ALTER TABLE
    IF EXISTS "common_usersurveys"
    DROP COLUMN IF EXISTS "genders",
    DROP COLUMN IF EXISTS "age_range",
    DROP COLUMN IF EXISTS "client_types",
    DROP COLUMN IF EXISTS "set_id";

COMMIT;