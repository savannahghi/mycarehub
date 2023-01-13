BEGIN;

ALTER TABLE
    IF EXISTS "common_securityquestionresponse"
    DROP CONSTRAINT IF EXISTS "common_securityquestionresponse_organisation_id_fkey",
    DROP COLUMN IF EXISTS "organisation_id";

COMMIT;