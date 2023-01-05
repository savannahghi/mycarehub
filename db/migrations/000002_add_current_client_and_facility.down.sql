BEGIN;

ALTER TABLE
    IF EXISTS "caregivers_caregiver"
    DROP CONSTRAINT IF EXISTS "caregivers_caregiver_current_client_fkey",
    DROP CONSTRAINT IF EXISTS "caregivers_caregiver_current_facility_fkey",
    DROP COLUMN IF EXISTS "current_client",
    DROP COLUMN IF EXISTS "current_facility",
    ADD COLUMN IF NOT EXISTS "program_id" uuid;

COMMIT;