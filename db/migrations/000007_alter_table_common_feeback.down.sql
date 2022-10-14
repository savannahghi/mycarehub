BEGIN;

ALTER TABLE common_feedback
    DROP COLUMN gender  ,
    DROP COLUMN facility_id;

COMMIT;