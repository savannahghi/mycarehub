BEGIN;

ALTER TABLE
    IF EXISTS "screeningtools_screeningtoolsresponse"
    DROP CONSTRAINT IF EXISTS "screeningtools_screeningtoolsresponse_created_by_fkey",
    DROP CONSTRAINT IF EXISTS "screeningtools_screeningtoolsresponse_updated_by_fkey",
    DROP CONSTRAINT IF EXISTS "screeningtools_screeningtoolsresponse_client_id_fkey",
    DROP CONSTRAINT IF EXISTS "screeningtools_screeningtoolsresponse_organisation_id_fkey",
    DROP CONSTRAINT IF EXISTS "screeningtools_screeningtoolsresponse_question_id_fkey",
    DROP CONSTRAINT IF EXISTS "screeningtools_screeningtoolsresponse_program_id_fkey";

DROP TABLE IF EXISTS "screeningtools_screeningtoolsresponse";

ALTER TABLE
    IF EXISTS "screeningtools_screeningtoolsquestion"
    DROP CONSTRAINT IF EXISTS "screeningtools_screeningtoolsquestion_created_by_fkey",
    DROP CONSTRAINT IF EXISTS "screeningtools_screeningtoolsquestion_updated_by_fkey",
    DROP CONSTRAINT IF EXISTS "screeningtools_screeningtoolsquestion_organisation_id_fkey",
    DROP CONSTRAINT IF EXISTS "screeningtools_screeningtoolsquestion_program_id_fkey";

DROP TABLE IF EXISTS "screeningtools_screeningtoolsquestion";



COMMIT;