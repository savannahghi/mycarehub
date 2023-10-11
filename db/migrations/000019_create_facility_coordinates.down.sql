BEGIN;

ALTER TABLE
    IF EXISTS "common_facility"
    DROP COLUMN IF EXISTS "latitude";

ALTER TABLE
    IF EXISTS "common_facility"
    DROP COLUMN IF EXISTS "longitude";

COMMIT;