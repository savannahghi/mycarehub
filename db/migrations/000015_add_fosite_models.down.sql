BEGIN;

ALTER TABLE
    IF EXISTS "oauth_session"
    DROP CONSTRAINT IF EXISTS "oauth_session_user_id_fkey";

ALTER TABLE
    IF EXISTS "oauth_access_token"
    DROP CONSTRAINT IF EXISTS "oauth_access_token_client_id_fkey";

ALTER TABLE
    IF EXISTS "oauth_access_token"
    DROP CONSTRAINT IF EXISTS "oauth_access_token_session_id_fkey";

ALTER TABLE
    IF EXISTS "oauth_authorization_code"
    DROP CONSTRAINT IF EXISTS "oauth_authorization_code_session_id_fkey";

ALTER TABLE
    IF EXISTS "oauth_authorization_code"
    DROP CONSTRAINT IF EXISTS "oauth_authorization_code_client_id_fkey";

ALTER TABLE
    IF EXISTS "oauth_pkce"
    DROP CONSTRAINT IF EXISTS "oauth_pkce_session_id_fkey";

ALTER TABLE
    IF EXISTS "oauth_pkce"
    DROP CONSTRAINT IF EXISTS "oauth_pkce_client_id_fkey";

ALTER TABLE
    IF EXISTS "oauth_refresh_token"
    DROP CONSTRAINT IF EXISTS "oauth_refresh_token_client_id_fkey";

ALTER TABLE
    IF EXISTS "oauth_refresh_token"
    DROP CONSTRAINT IF EXISTS "oauth_refresh_token_session_id_fkey";

DROP TABLE IF EXISTS "oauth_client_jwt";

DROP TABLE IF EXISTS "oauth_client";

DROP TABLE IF EXISTS "oauth_session";

DROP TABLE IF EXISTS "oauth_access_token";

DROP TABLE IF EXISTS "oauth_authorization_code";

DROP TABLE IF EXISTS "oauth_pkce";

DROP TABLE IF EXISTS "oauth_refresh_token";

COMMIT;