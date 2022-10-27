BEGIN;

ALTER TABLE common_feedback
    ADD COLUMN IF NOT EXISTS gender varchar(15) NOT NULL,
    ADD COLUMN IF NOT EXISTS facility_id varchar(255) NOT NULL;

COMMIT;