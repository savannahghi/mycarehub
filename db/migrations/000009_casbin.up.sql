BEGIN;

CREATE TABLE IF NOT EXISTS casbin_policy (
    id serial PRIMARY KEY,
    active boolean,
    created timestamp WITH time zone,
    created_by uuid,
    updated timestamp WITH time zone,
    updated_by uuid,
    deleted_at timestamp WITH time zone,
    -- organisation_id uuid CONSTRAINT casbin_policy_organisation_id REFERENCES common_organisation deferrable initially deferred,
    p_type varchar(255),
    v0 varchar(255),
    v1 varchar(255),
    v2 varchar(255),
    v3 varchar(255),
    v4 varchar(255),
    v5 varchar(255)
);

CREATE TABLE IF NOT EXISTS user_groups (
    id uuid NOT NULL PRIMARY KEY,
    active boolean NOT NULL,
    created timestamp WITH time zone NOT NULL,
    created_by uuid,
    updated timestamp WITH time zone NOT NULL,
    updated_by uuid,
    organisation_id uuid NOT NULL CONSTRAINT user_groups_organisation REFERENCES common_organisation deferrable initially deferred,
    deleted_at timestamp WITH time zone,
    name text NOT NULL
);

CREATE TABLE IF NOT EXISTS users_user_user_groups (
    id serial PRIMARY KEY,
    user_id uuid NOT NULL CONSTRAINT user_id_users_user REFERENCES users_user deferrable initially deferred,
    group_id uuid NOT NULL CONSTRAINT group_id_user_groups REFERENCES user_groups deferrable initially deferred,
    CONSTRAINT user_id_group_id_unique UNIQUE (user_id, group_id)
);

COMMIT;