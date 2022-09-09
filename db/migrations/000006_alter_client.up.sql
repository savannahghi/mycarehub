BEGIN;

CREATE TABLE IF NOT EXISTS clients_client_facilities (
    id serial PRIMARY KEY,
    client_id uuid NOT NULL CONSTRAINT clients_client_facilities_client_id_2fdf3b17_fk_clients_client_id REFERENCES clients_client deferrable initially deferred,
    facility_id uuid NOT NULL CONSTRAINT clients_client_faciliti_facility_id_e6dbfdf8_fk_common_fa REFERENCES common_facility deferrable initially deferred,
    CONSTRAINT clients_client_facilities_client_id_facility_id_51a27bda_uniq UNIQUE (client_id, facility_id)
);

CREATE INDEX IF NOT EXISTS clients_client_facilities_client_id_4fdf3b17 ON clients_client_facilities (client_id);

CREATE INDEX IF NOT EXISTS clients_client_facilities_facility_id_f6dbfdf8 ON clients_client_facilities (facility_id);

COMMIT;