BEGIN;

ALTER TABLE
    IF EXISTS "communities_community"
    ADD COLUMN IF NOT EXISTS "invite_only" boolean DEFAULT false,
    ADD COLUMN IF NOT EXISTS "discoverable" boolean DEFAULT false,
    DROP COLUMN IF EXISTS "room_id"; 

COMMIT;