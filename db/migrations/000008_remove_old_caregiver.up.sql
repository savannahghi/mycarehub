BEGIN;

DROP INDEX IF EXISTS clients_caregiver_organisation_id_14fee135;
ALTER TABLE clients_client DROP COLUMN IF EXISTS caregiver_id;
DROP TABLE IF EXISTS clients_caregiver; 

COMMIT;