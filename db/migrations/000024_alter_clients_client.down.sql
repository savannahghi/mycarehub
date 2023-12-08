BEGIN;

ALTER TABLE
    IF EXISTS "clients_client"
    ALTER COLUMN "client_types" SET NOT NULL;

ALTER TABLE
    IF EXISTS "clients_client"
    ALTER COLUMN "current_facility_id" SET NOT NULL;

ALTER TABLE
    IF EXISTS "clients_client"
    ALTER COLUMN "organisation_id" SET NOT NULL;

ALTER TABLE
    IF EXISTS "clients_client"
    ALTER COLUMN "program_id" SET NOT NULL;

COMMIT;