BEGIN;

ALTER TABLE
    IF EXISTS "clients_client"
    DROP CONSTRAINT IF EXISTS "clients_client_user_id_program_id_unique",
    DROP CONSTRAINT IF EXISTS "clients_client_fhir_patient_id_unique",
    ADD CONSTRAINT "clients_client_program_id_fhir_patient_id_user_id_key" UNIQUE ("program_id", "fhir_patient_id", "user_id");

COMMIT;