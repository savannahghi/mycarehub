BEGIN;

ALTER TABLE
    IF EXISTS "staff_staff"
    ADD COLUMN IF NOT EXISTS "is_organisation_admin" boolean;

COMMIT;