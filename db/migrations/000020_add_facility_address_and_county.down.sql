BEGIN;

ALTER TABLE
    IF EXISTS "common_facility"
    DROP COLUMN IF EXISTS "address";

ALTER TABLE
    IF EXISTS "common_facility"
    DROP COLUMN IF EXISTS "county";

COMMIT;