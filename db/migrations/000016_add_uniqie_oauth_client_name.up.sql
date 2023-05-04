BEGIN;

ALTER TABLE
    IF EXISTS "oauth_client"
    ADD CONSTRAINT "oauth_client_name_unique" UNIQUE ("name");

COMMIT;