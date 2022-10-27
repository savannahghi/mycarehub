BEGIN;

ALTER TABLE common_feedback
    ALTER gender SET NOT NULL,
    ALTER facility_id SET NOT NULL;

COMMIT;