BEGIN;

ALTER TABLE
    IF EXISTS "service_booking"
    DROP COLUMN IF EXISTS "verification_code_status";

COMMIT;