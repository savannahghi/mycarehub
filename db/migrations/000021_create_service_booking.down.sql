BEGIN;

ALTER TABLE
    IF EXISTS "service_booking"
    DROP CONSTRAINT IF EXISTS "service_booking_created_by_fkey";

ALTER TABLE
    IF EXISTS "service_booking"
    DROP CONSTRAINT IF EXISTS "service_booking_updated_by_fkey";

ALTER TABLE
    IF EXISTS "service_booking"
    DROP CONSTRAINT IF EXISTS "service_booking_client_id_fkey";

ALTER TABLE
    IF EXISTS "service_booking"
    DROP CONSTRAINT IF EXISTS "service_booking_organisation_id_fkey";

ALTER TABLE
    IF EXISTS "service_booking"
    DROP CONSTRAINT IF EXISTS "service_booking_program_id_fkey";

DROP TABLE IF EXISTS "service_booking";

COMMIT;