BEGIN;
ALTER TABLE common_usersurveys DROP COLUMN IF EXISTS submitted_at;
COMMIT;
