BEGIN;

CREATE TABLE IF NOT EXISTS caregivers_caregiver (
    id uuid NOT NULL PRIMARY KEY,
    active boolean NOT NULL,
    created timestamp WITH time zone NOT NULL,
    created_by uuid,
    updated timestamp WITH time zone NOT NULL,
    updated_by uuid,
    organisation_id uuid NOT NULL REFERENCES common_organisation DEFERRABLE INITIALLY DEFERRED,
    user_id uuid NOT NULL UNIQUE REFERENCES users_user DEFERRABLE INITIALLY DEFERRED,
    caregiver_number varchar(255)
);

CREATE TABLE IF NOT EXISTS caregivers_caregiver_client (
    caregiver_id uuid NOT NULL REFERENCES caregivers_caregiver DEFERRABLE INITIALLY DEFERRED,
    client_id uuid NOT NULL REFERENCES clients_client DEFERRABLE INITIALLY DEFERRED,
    active boolean NOT NULL,
    created timestamp WITH time zone NOT NULL,
    created_by uuid,
    updated timestamp WITH time zone NOT NULL,
    updated_by uuid,
    organisation_id uuid NOT NULL REFERENCES common_organisation DEFERRABLE INITIALLY DEFERRED,
    assigned_by uuid NOT NULL REFERENCES staff_staff DEFERRABLE INITIALLY DEFERRED,
    relationship_type varchar(32) NOT NULL,
    caregiver_consent boolean ,
    caregiver_consent_at timestamp WITH time zone,
    client_consent boolean ,
    client_consent_at timestamp WITH time zone,

    PRIMARY KEY (caregiver_id, client_id)
);

COMMIT;