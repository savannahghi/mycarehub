BEGIN;

ALTER TABLE
    IF EXISTS "caregivers_caregiver"
    DROP CONSTRAINT IF EXISTS "caregivers_caregiver_organisation_id_user_id_unique";

COMMIT;