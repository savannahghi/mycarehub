BEGIN;

ALTER TABLE
    IF EXISTS "service_booking"
    ADD COLUMN IF NOT EXISTS "verification_code_status" varchar(15) NOT NULL;

COMMIT;