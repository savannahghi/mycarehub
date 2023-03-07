BEGIN;

ALTER TABLE
    IF EXISTS "users_user"
    DROP COLUMN IF EXISTS "current_usertype";

COMMIT;