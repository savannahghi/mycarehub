BEGIN;

ALTER TABLE
    IF EXISTS "common_facility"
    ADD COLUMN IF NOT EXISTS "address" varchar(100);

ALTER TABLE
    IF EXISTS "common_facility"
    ADD COLUMN IF NOT EXISTS "county" varchar(100);

COMMIT;