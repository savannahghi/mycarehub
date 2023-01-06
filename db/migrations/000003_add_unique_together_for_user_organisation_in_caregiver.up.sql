BEGIN;

ALTER TABLE
    IF EXISTS "caregivers_caregiver"
    ADD CONSTRAINT "caregivers_caregiver_organisation_id_user_id_unique" UNIQUE ("user_id", "organisation_id");

COMMIT;