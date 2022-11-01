BEGIN;

ALTER TABLE common_facility DROP COLUMN IF EXISTS program_id;

ALTER TABLE common_auditlog DROP COLUMN IF EXISTS program_id;

ALTER TABLE users_user DROP COLUMN IF EXISTS program_id;

ALTER TABLE common_contact DROP COLUMN IF EXISTS program_id;

ALTER TABLE users_termsofservice DROP COLUMN IF EXISTS program_id;

ALTER TABLE clients_securityquestion DROP COLUMN IF EXISTS program_id;

ALTER TABLE clients_securityquestionresponse DROP COLUMN IF EXISTS program_id;

ALTER TABLE clients_client DROP COLUMN IF EXISTS program_id;

ALTER TABLE clients_clientfacility DROP COLUMN IF EXISTS program_id;

ALTER TABLE staff_staff DROP COLUMN IF EXISTS program_id;

ALTER TABLE clients_healthdiaryentry DROP COLUMN IF EXISTS program_id;

ALTER TABLE clients_servicerequest DROP COLUMN IF EXISTS program_id;

ALTER TABLE staff_servicerequest DROP COLUMN IF EXISTS program_id;

ALTER TABLE clients_healthdiaryquote DROP COLUMN IF EXISTS program_id;

ALTER TABLE authority_authorityrole DROP COLUMN IF EXISTS program_id;

ALTER TABLE authority_authoritypermission DROP COLUMN IF EXISTS program_id;

ALTER TABLE communities_community DROP COLUMN IF EXISTS program_id;

ALTER TABLE clients_identifier DROP COLUMN IF EXISTS program_id;

ALTER TABLE clients_relatedperson DROP COLUMN IF EXISTS program_id;

ALTER TABLE screeningtools_screeningtoolsquestion DROP COLUMN IF EXISTS program_id;

ALTER TABLE screeningtools_screeningtoolsresponse DROP COLUMN IF EXISTS program_id;

ALTER TABLE appointments_appointment DROP COLUMN IF EXISTS program_id;

ALTER TABLE common_notification DROP COLUMN IF EXISTS program_id;

ALTER TABLE common_usersurveys DROP COLUMN IF EXISTS program_id;

ALTER TABLE common_feedback DROP COLUMN IF EXISTS program_id;

ALTER TABLE questionnaires_questionnaire DROP COLUMN IF EXISTS program_id;

ALTER TABLE questionnaires_screeningtool DROP COLUMN IF EXISTS program_id;

ALTER TABLE questionnaires_question DROP COLUMN IF EXISTS program_id;

ALTER TABLE questionnaires_questioninputchoice DROP COLUMN IF EXISTS program_id;

ALTER TABLE questionnaires_screeningtoolresponse DROP COLUMN IF EXISTS program_id;

ALTER TABLE questionnaires_screeningtoolquestionresponse DROP COLUMN IF EXISTS program_id;

ALTER TABLE caregivers_caregiver DROP COLUMN IF EXISTS program_id;

ALTER TABLE caregivers_caregiver_client DROP COLUMN IF EXISTS program_id;

ALTER TABLE common_address DROP COLUMN IF EXISTS program_id;

COMMIT;