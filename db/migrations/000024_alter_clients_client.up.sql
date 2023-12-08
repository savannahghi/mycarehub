BEGIN;

ALTER TABLE
    IF EXISTS "clients_client"
    ALTER COLUMN "client_types" DROP NOT NULL;

ALTER TABLE
    IF EXISTS "clients_client"
    ALTER COLUMN "current_facility_id" DROP NOT NULL;

ALTER TABLE
    IF EXISTS "clients_client"
    ALTER COLUMN "organisation_id" DROP NOT NULL;

ALTER TABLE
    IF EXISTS "clients_client"
    ALTER COLUMN "program_id" DROP NOT NULL;

COMMIT;