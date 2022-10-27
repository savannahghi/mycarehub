BEGIN;

ALTER TABLE common_feedback
    ADD COLUMN IF NOT EXISTS gender varchar(15),
    ADD COLUMN IF NOT EXISTS facility_id varchar(255);

COMMIT;