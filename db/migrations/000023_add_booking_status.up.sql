BEGIN;

ALTER TABLE
    IF EXISTS "service_booking"
    ADD COLUMN IF NOT EXISTS "booking_status" varchar(15);

COMMIT;