BEGIN;

ALTER TABLE
    IF EXISTS "caregivers_caregiver"
    DROP CONSTRAINT IF EXISTS "caregivers_caregiver_program_id_fkey",
    DROP COLUMN IF EXISTS "program_id",
    ADD COLUMN IF NOT EXISTS "current_client" uuid,
    ADD COLUMN IF NOT EXISTS "current_facility" uuid,
    ADD CONSTRAINT "caregivers_caregiver_current_client_fkey" FOREIGN KEY ("current_client") REFERENCES "clients_client" ("id"),
    ADD CONSTRAINT "caregivers_caregiver_current_facility_fkey" FOREIGN KEY ("current_facility") REFERENCES "common_facility" ("id");

COMMIT;