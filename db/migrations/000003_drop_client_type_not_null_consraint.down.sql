BEGIN;

ALTER TABLE IF EXISTS "common_program"
    ADD COLUMN IF NOT EXISTS client_types varchar(64)[];

ALTER TABLE IF EXISTS "common_program"
    ALTER COLUMN client_types  SET NOT NULL;

COMMIT;