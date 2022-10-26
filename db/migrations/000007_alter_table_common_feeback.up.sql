BEGIN;

ALTER TABLE common_feedback
    ADD COLUMN gender varchar(15),
    ADD COLUMN facility_id varchar(255);

COMMIT;