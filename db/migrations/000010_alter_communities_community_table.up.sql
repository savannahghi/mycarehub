BEGIN;

ALTER TABLE
    IF EXISTS "communities_community"
    DROP COLUMN IF EXISTS "invite_only",
    DROP COLUMN IF EXISTS "discoverable",
    ADD COLUMN IF NOT EXISTS "room_id" text;

COMMIT;