BEGIN;

ALTER TABLE
    IF EXISTS "users_user"
    ADD COLUMN IF NOT EXISTS "current_usertype" varchar(255);

COMMIT;