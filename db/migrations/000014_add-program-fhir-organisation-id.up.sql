BEGIN;

ALTER TABLE
    IF EXISTS "common_program"
    ADD COLUMN IF NOT EXISTS "fhir_organisation_id" text;

COMMIT;