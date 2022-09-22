BEGIN;

CREATE TABLE IF NOT EXISTS clients_caregiver (
    id uuid NOT NULL PRIMARY KEY,
    active boolean NOT NULL,
    created timestamp WITH time zone NOT NULL,
    created_by uuid,
    updated timestamp WITH time zone NOT NULL,
    updated_by uuid,
    deleted_at timestamp WITH time zone,
    first_name text NOT NULL,
    last_name text NOT NULL,
    caregiver_type varchar(64) NOT NULL,
    phone_number text,
    organisation_id uuid NOT NULL CONSTRAINT clients_caregiver_organisation_id_14fee135_fk_common_or REFERENCES common_organisation deferrable initially deferred
);

ALTER TABLE clients_client ADD COLUMN IF NOT EXISTS caregiver_id uuid UNIQUE CONSTRAINT clients_client_caregiver_id_0033fe31_fk_clients_caregiver_id REFERENCES clients_caregiver deferrable initially deferred;

CREATE INDEX IF NOT EXISTS clients_caregiver_organisation_id_14fee135 ON clients_caregiver (organisation_id);

COMMIT;