BEGIN;

ALTER TABLE common_feedback
    ALTER gender DROP NOT NULL,
    ALTER facility_id DROP NOT NULL;

COMMIT;