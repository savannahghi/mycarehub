BEGIN;

ALTER TABLE clients_identifier
RENAME TO common_identifiers;

COMMIT;