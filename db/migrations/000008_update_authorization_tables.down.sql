BEGIN;

ALTER TABLE
    IF EXISTS "authority_authoritypermission"
    ADD COLUMN IF NOT EXISTS "organisation_id" uuid NOT NULL,
    DROP COLUMN IF EXISTS "description",
    DROP COLUMN IF EXISTS "category",
    DROP COLUMN IF EXISTS "scope"; 

ALTER TABLE
    IF EXISTS "authority_authoritypermission"
    ADD CONSTRAINT "authority_authoritypermission_organisation_id_fkey" FOREIGN KEY ("organisation_id") REFERENCES "common_organisation" ("id");

ALTER TABLE
    IF EXISTS "authority_authorityrole"
    DROP CONSTRAINT IF EXISTS "authority_authorityrole_program_id_fkey",
    DROP COLUMN IF EXISTS "program_id";

COMMIT;