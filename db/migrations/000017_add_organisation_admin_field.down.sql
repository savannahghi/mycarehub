BEGIN;

ALTER TABLE
    IF EXISTS "staff_staff"
    DROP COLUMN IF EXISTS "is_organisation_admin";

COMMIT;