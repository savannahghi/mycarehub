BEGIN;

ALTER TABLE
    IF EXISTS "clients_client"
    DROP CONSTRAINT IF EXISTS "clients_client_program_id_fhir_patient_id_user_id_key",
    ADD CONSTRAINT "clients_client_user_id_program_id_unique" UNIQUE ("user_id", "program_id"),
    ADD CONSTRAINT "clients_client_fhir_patient_id_unique" UNIQUE ("fhir_patient_id");

COMMIT;