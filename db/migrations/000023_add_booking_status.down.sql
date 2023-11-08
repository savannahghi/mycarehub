BEGIN;

ALTER TABLE
    IF EXISTS "service_booking"
    DROP COLUMN IF EXISTS "booking_status";

COMMIT;