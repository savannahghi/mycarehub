BEGIN;

ALTER TABLE
    IF EXISTS "common_program"
    DROP COLUMN IF EXISTS "fhir_organisation_id";

COMMIT;