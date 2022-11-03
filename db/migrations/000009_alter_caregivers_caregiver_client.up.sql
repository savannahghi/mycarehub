BEGIN;

ALTER TABLE IF EXISTS caregivers_caregiver_client ALTER COLUMN caregiver_consent SET DATA TYPE varchar(255) USING (caregiver_consent::varchar(255));

ALTER TABLE IF EXISTS caregivers_caregiver_client ALTER COLUMN client_consent SET DATA TYPE varchar(255) USING (client_consent::varchar(255));

COMMIT;

