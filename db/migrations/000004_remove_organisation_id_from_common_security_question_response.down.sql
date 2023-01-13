BEGIN;

ALTER TABLE
    IF EXISTS "common_securityquestionresponse"
    ADD COLUMN IF NOT EXISTS "organisation_id" uuid,
    ADD CONSTRAINT "common_securityquestionresponse_organisation_id_fkey" FOREIGN KEY ("organisation_id") REFERENCES "common_organisation" ("id");

COMMIT;