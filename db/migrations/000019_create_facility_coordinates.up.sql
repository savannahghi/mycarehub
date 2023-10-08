BEGIN;

ALTER TABLE
    IF EXISTS "common_facility"
    ADD COLUMN IF NOT EXISTS "latitude" float;

ALTER TABLE
    IF EXISTS "common_facility"
    ADD COLUMN IF NOT EXISTS "longitude" float;

COMMIT;