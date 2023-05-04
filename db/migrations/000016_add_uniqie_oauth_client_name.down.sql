BEGIN;

ALTER TABLE
    IF EXISTS "oauth_client"
    DROP CONSTRAINT IF EXISTS "oauth_client_name_unique";
COMMIT;