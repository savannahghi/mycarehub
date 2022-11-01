BEGIN;

ALTER TABLE common_facility ADD COLUMN IF NOT EXISTS program_id uuid DEFAULT NULL;

ALTER TABLE common_auditlog ADD COLUMN IF NOT EXISTS program_id uuid DEFAULT NULL;

ALTER TABLE users_user ADD COLUMN IF NOT EXISTS program_id uuid DEFAULT NULL;

ALTER TABLE common_contact ADD COLUMN IF NOT EXISTS program_id uuid DEFAULT NULL;

ALTER TABLE users_termsofservice ADD COLUMN IF NOT EXISTS program_id uuid DEFAULT NULL;

ALTER TABLE clients_securityquestion ADD COLUMN IF NOT EXISTS program_id uuid DEFAULT NULL;

ALTER TABLE clients_securityquestionresponse ADD COLUMN IF NOT EXISTS program_id uuid DEFAULT NULL;

ALTER TABLE clients_client ADD COLUMN IF NOT EXISTS program_id uuid DEFAULT NULL;

ALTER TABLE clients_clientfacility ADD COLUMN IF NOT EXISTS program_id uuid DEFAULT NULL;

ALTER TABLE staff_staff ADD COLUMN IF NOT EXISTS program_id uuid DEFAULT NULL;

ALTER TABLE clients_healthdiaryentry ADD COLUMN IF NOT EXISTS program_id uuid DEFAULT NULL;

ALTER TABLE clients_servicerequest ADD COLUMN IF NOT EXISTS program_id uuid DEFAULT NULL;

ALTER TABLE staff_servicerequest ADD COLUMN IF NOT EXISTS program_id uuid DEFAULT NULL;

ALTER TABLE clients_healthdiaryquote ADD COLUMN IF NOT EXISTS program_id uuid DEFAULT NULL;

ALTER TABLE authority_authorityrole ADD COLUMN IF NOT EXISTS program_id uuid DEFAULT NULL;

ALTER TABLE authority_authoritypermission ADD COLUMN IF NOT EXISTS program_id uuid DEFAULT NULL;

ALTER TABLE communities_community ADD COLUMN IF NOT EXISTS program_id uuid DEFAULT NULL;

ALTER TABLE clients_identifier ADD COLUMN IF NOT EXISTS program_id uuid DEFAULT NULL;

ALTER TABLE clients_relatedperson ADD COLUMN IF NOT EXISTS program_id uuid DEFAULT NULL;

ALTER TABLE screeningtools_screeningtoolsquestion ADD COLUMN IF NOT EXISTS program_id uuid DEFAULT NULL;

ALTER TABLE screeningtools_screeningtoolsresponse ADD COLUMN IF NOT EXISTS program_id uuid DEFAULT NULL;

ALTER TABLE appointments_appointment ADD COLUMN IF NOT EXISTS program_id uuid DEFAULT NULL;

ALTER TABLE common_notification ADD COLUMN IF NOT EXISTS program_id uuid DEFAULT NULL;

ALTER TABLE common_usersurveys ADD COLUMN IF NOT EXISTS program_id uuid DEFAULT NULL;

ALTER TABLE common_feedback ADD COLUMN IF NOT EXISTS program_id uuid DEFAULT NULL;

ALTER TABLE questionnaires_questionnaire ADD COLUMN IF NOT EXISTS program_id uuid DEFAULT NULL;

ALTER TABLE questionnaires_screeningtool ADD COLUMN IF NOT EXISTS program_id uuid DEFAULT NULL;

ALTER TABLE questionnaires_question ADD COLUMN IF NOT EXISTS program_id uuid DEFAULT NULL;

ALTER TABLE questionnaires_questioninputchoice ADD COLUMN IF NOT EXISTS program_id uuid DEFAULT NULL;

ALTER TABLE questionnaires_screeningtoolresponse ADD COLUMN IF NOT EXISTS program_id uuid DEFAULT NULL;

ALTER TABLE questionnaires_screeningtoolquestionresponse ADD COLUMN IF NOT EXISTS program_id uuid DEFAULT NULL;

ALTER TABLE caregivers_caregiver ADD COLUMN IF NOT EXISTS program_id uuid DEFAULT NULL;

ALTER TABLE caregivers_caregiver_client ADD COLUMN IF NOT EXISTS program_id uuid DEFAULT NULL;

ALTER TABLE common_address ADD COLUMN IF NOT EXISTS program_id uuid DEFAULT NULL;

COMMIT;