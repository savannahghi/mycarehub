BEGIN;

CREATE TABLE IF NOT EXISTS "oauth_client_jwt" (
    "id" uuid PRIMARY KEY NOT NULL,
    "created_by" uuid,
    "created" timestamp NOT NULL,
    "updated_by" uuid,
    "updated" timestamp NOT NULL,
    "deleted_at" timestamp,
    "active" boolean NOT NULL,
    "jti" varchar(256),
    "expires_at" timestamp NOT NULL
);

CREATE TABLE IF NOT EXISTS "oauth_client" (
    "id" uuid PRIMARY KEY NOT NULL,
    "created_by" uuid,
    "created" timestamp NOT NULL,
    "updated_by" uuid,
    "updated" timestamp NOT NULL,
    "deleted_at" timestamp,
    "active" boolean NOT NULL,
    "name" varchar(256),
    "secret" varchar(256),
    "rotated_secrets" varchar(256) [],
    "public" boolean NOT NULL,
    "redirect_uris" varchar(256) [],
    "scopes" varchar(256) [],
    "audience" varchar(256) [],
    "grants" varchar(256) [],
    "response_types" varchar(256) [],
    "token_endpoint_auth_method" varchar(256)
);

CREATE TABLE IF NOT EXISTS "oauth_session" (
    "id" uuid PRIMARY KEY NOT NULL,
    "created_by" uuid,
    "created" timestamp NOT NULL,
    "updated_by" uuid,
    "updated" timestamp NOT NULL,
    "deleted_at" timestamp,
    "client_id" uuid,
    "username" varchar(150),
    "subject" varchar(150),
    "expires_at" jsonb,
    "extra" jsonb,
    "user_id" uuid
);

CREATE TABLE IF NOT EXISTS "oauth_access_token" (
    "id" uuid PRIMARY KEY NOT NULL,
    "created_by" uuid,
    "created" timestamp NOT NULL,
    "updated_by" uuid,
    "updated" timestamp NOT NULL,
    "deleted_at" timestamp,
    "active" boolean NOT NULL,
    "signature" varchar(256) unique,
    "requested_at" timestamp NOT NULL,
    "requested_scopes" varchar(256) [],
    "granted_scopes" varchar(256) [],
    "form" jsonb,
    "requested_audience" varchar(256) [],
    "granted_audience" varchar(256) [],
    "client_id" uuid,
    "session_id" uuid
);

CREATE TABLE IF NOT EXISTS "oauth_authorization_code" (
    "id" uuid PRIMARY KEY NOT NULL,
    "created_by" uuid,
    "created" timestamp NOT NULL,
    "updated_by" uuid,
    "updated" timestamp NOT NULL,
    "deleted_at" timestamp,
    "active" boolean NOT NULL,
    "code" varchar(256),
    "requested_at" timestamp NOT NULL,
    "requested_scopes" varchar(256) [],
    "granted_scopes" varchar(256) [],
    "form" jsonb,
    "requested_audience" varchar(256) [],
    "granted_audience" varchar(256) [],
    "session_id" uuid,
    "client_id" uuid
);

CREATE TABLE IF NOT EXISTS "oauth_pkce" (
    "id" uuid PRIMARY KEY NOT NULL,
    "created_by" uuid,
    "created" timestamp NOT NULL,
    "updated_by" uuid,
    "updated" timestamp NOT NULL,
    "deleted_at" timestamp,
    "active" boolean NOT NULL,
    "signature" varchar(256) unique,
    "requested_at" timestamp NOT NULL,
    "requested_scopes" varchar(256) [],
    "granted_scopes" varchar(256) [],
    "form" jsonb,
    "requested_audience" varchar(256) [],
    "granted_audience" varchar(256) [],
    "session_id" uuid,
    "client_id" uuid
);

CREATE TABLE IF NOT EXISTS "oauth_refresh_token" (
    "id" uuid PRIMARY KEY NOT NULL,
    "created_by" uuid,
    "created" timestamp NOT NULL,
    "updated_by" uuid,
    "updated" timestamp NOT NULL,
    "deleted_at" timestamp,
    "active" boolean NOT NULL,
    "signature" varchar(256) unique,
    "requested_at" timestamp NOT NULL,
    "requested_scopes" varchar(256) [],
    "granted_scopes" varchar(256) [],
    "form" jsonb,
    "requested_audience" varchar(256) [],
    "granted_audience" varchar(256) [],
    "client_id" uuid,
    "session_id" uuid
);

ALTER TABLE
    IF EXISTS "oauth_session"
ADD
    CONSTRAINT "oauth_session_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "oauth_access_token"
ADD
    CONSTRAINT "oauth_access_token_client_id_fkey" FOREIGN KEY ("client_id") REFERENCES "oauth_client" ("id");

ALTER TABLE
    IF EXISTS "oauth_access_token"
ADD
    CONSTRAINT "oauth_access_token_session_id_fkey" FOREIGN KEY ("session_id") REFERENCES "oauth_session" ("id");

ALTER TABLE
    IF EXISTS "oauth_authorization_code"
ADD
    CONSTRAINT "oauth_authorization_code_session_id_fkey" FOREIGN KEY ("session_id") REFERENCES "oauth_session" ("id");

ALTER TABLE
    IF EXISTS "oauth_authorization_code"
ADD
    CONSTRAINT "oauth_authorization_code_client_id_fkey" FOREIGN KEY ("client_id") REFERENCES "oauth_client" ("id");

ALTER TABLE
    IF EXISTS "oauth_pkce"
ADD
    CONSTRAINT "oauth_pkce_session_id_fkey" FOREIGN KEY ("session_id") REFERENCES "oauth_session" ("id");

ALTER TABLE
    IF EXISTS "oauth_pkce"
ADD
    CONSTRAINT "oauth_pkce_client_id_fkey" FOREIGN KEY ("client_id") REFERENCES "oauth_client" ("id");

ALTER TABLE
    IF EXISTS "oauth_refresh_token"
ADD
    CONSTRAINT "oauth_refresh_token_client_id_fkey" FOREIGN KEY ("client_id") REFERENCES "oauth_client" ("id");

ALTER TABLE
    IF EXISTS "oauth_refresh_token"
ADD
    CONSTRAINT "oauth_refresh_token_session_id_fkey" FOREIGN KEY ("session_id") REFERENCES "oauth_session" ("id");

COMMIT;