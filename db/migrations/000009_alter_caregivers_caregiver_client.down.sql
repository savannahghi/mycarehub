BEGIN;

ALTER TABLE IF EXISTS caregivers_caregiver_client ALTER COLUMN caregiver_consent SET DATA TYPE boolean USING (caregiver_consent::boolean);

ALTER TABLE IF EXISTS caregivers_caregiver_client ALTER COLUMN client_consent SET DATA TYPE boolean USING (client_consent::boolean);

COMMIT;