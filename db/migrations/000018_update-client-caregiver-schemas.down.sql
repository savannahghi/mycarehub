BEGIN;

ALTER TABLE
    IF EXISTS "clients_healthdiaryentry"
    DROP COLUMN IF EXISTS "caregiver_id";

ALTER TABLE
    IF EXISTS "clients_servicerequest"
    DROP COLUMN IF EXISTS "caregiver_id";

ALTER TABLE
    IF EXISTS "questionnaires_screeningtoolresponse"
    DROP COLUMN IF EXISTS "caregiver_id";

ALTER TABLE
    IF EXISTS "common_usersurveys"
    DROP COLUMN IF EXISTS "caregiver_id";

ALTER TABLE
    IF EXISTS "appointments_appointment"
    DROP COLUMN IF EXISTS "caregiver_id";

COMMIT;
