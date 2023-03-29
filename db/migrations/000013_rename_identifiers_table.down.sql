BEGIN;

ALTER TABLE common_identifiers
RENAME TO clients_identifier;

COMMIT;