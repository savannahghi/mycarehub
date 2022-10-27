BEGIN;

ALTER TABLE common_feedback
    DROP COLUMN IF EXISTS gender  ,
    DROP COLUMN IF EXISTS facility_id;

COMMIT;