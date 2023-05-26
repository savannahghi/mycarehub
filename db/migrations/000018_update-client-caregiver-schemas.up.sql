BEGIN;

ALTER TABLE
    IF EXISTS "clients_healthdiaryentry"
    ADD COLUMN IF NOT EXISTS "caregiver_id" uuid;

ALTER TABLE
    IF EXISTS "clients_servicerequest"
    ADD COLUMN IF NOT EXISTS "caregiver_id" uuid;

ALTER TABLE
    IF EXISTS "questionnaires_screeningtoolresponse"
    ADD COLUMN IF NOT EXISTS "caregiver_id" uuid;

ALTER TABLE
    IF EXISTS "common_usersurveys"
    ADD COLUMN IF NOT EXISTS "caregiver_id" uuid;

ALTER TABLE
    IF EXISTS "appointments_appointment"
    ADD COLUMN IF NOT EXISTS "caregiver_id" uuid;


COMMIT;

