BEGIN;

ALTER TABLE
    IF EXISTS "authority_authoritypermission"
    DROP CONSTRAINT IF EXISTS "authority_authoritypermission_organisation_id_fkey",
    DROP COLUMN IF EXISTS "organisation_id",
    ADD COLUMN IF NOT EXISTS "description" TEXT NOT NULL,
    ADD COLUMN IF NOT EXISTS "category" varchar(255) NOT NULL,
    ADD COLUMN IF NOT EXISTS "scope" TEXT UNIQUE NOT NULL; 

ALTER TABLE
    IF EXISTS "authority_authorityrole"
    ADD COLUMN IF NOT EXISTS "program_id" uuid NOT NULL,
    ADD CONSTRAINT "authority_authorityrole_program_id_fkey" FOREIGN KEY ("program_id") REFERENCES "common_program" ("id");

COMMIT;