BEGIN;

ALTER TABLE common_feedback
    ADD COLUMN gender varchar(15) NOT NULL,
    ADD COLUMN facility_id varchar(255) NOT NULL;

COMMIT;