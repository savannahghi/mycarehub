BEGIN;

ALTER TABLE
    IF EXISTS "common_organisation"
DROP
    CONSTRAINT "common_organisation_created_by_fkey";

ALTER TABLE
    IF EXISTS "common_organisation"
DROP
    CONSTRAINT "common_organisation_updated_by_fkey";

ALTER TABLE
    IF EXISTS "common_facility"
DROP
    CONSTRAINT "common_facility_created_by_fkey";

ALTER TABLE
    IF EXISTS "common_facility"
DROP
    CONSTRAINT "common_facility_updated_by_fkey";

ALTER TABLE
    IF EXISTS "common_facility_identifier"
DROP
    CONSTRAINT "common_facility_identifier_created_by_fkey";

ALTER TABLE
    IF EXISTS "common_facility_identifier"
DROP
    CONSTRAINT "common_facility_identifier_updated_by_fkey";

ALTER TABLE
    IF EXISTS "common_facility_identifier"
DROP
    CONSTRAINT "common_facility_identifier_facility_id_fkey";

ALTER TABLE
    IF EXISTS "users_user"
DROP
    CONSTRAINT "users_user_created_by_fkey";

ALTER TABLE
    IF EXISTS "users_user"
DROP
    CONSTRAINT "users_user_updated_by_fkey";

ALTER TABLE
    IF EXISTS "users_user"
DROP
    CONSTRAINT "users_user_accepted_terms_of_service_id_fkey";

ALTER TABLE
    IF EXISTS "users_user"
DROP
    CONSTRAINT "users_user_current_program_id_fkey";

ALTER TABLE
    IF EXISTS "users_user"
DROP
    CONSTRAINT "users_user_organisation_id_fkey";

ALTER TABLE
    IF EXISTS "common_contact"
DROP
    CONSTRAINT "common_contact_created_by_fkey";

ALTER TABLE
    IF EXISTS "common_contact"
DROP
    CONSTRAINT "common_contact_updated_by_fkey";

ALTER TABLE
    IF EXISTS "common_contact"
DROP
    CONSTRAINT "common_contact_user_id_fkey";

ALTER TABLE
    IF EXISTS "common_contact"
DROP
    CONSTRAINT "common_contact_organisation_id_fkey";

ALTER TABLE
    IF EXISTS "common_program"
DROP
    CONSTRAINT "common_program_created_by_fkey";

ALTER TABLE
    IF EXISTS "common_program"
DROP
    CONSTRAINT "common_program_updated_by_fkey";

ALTER TABLE
    IF EXISTS "common_program"
DROP
    CONSTRAINT "common_program_organisation_id_fkey";

ALTER TABLE
    IF EXISTS "common_program_user"
DROP
    CONSTRAINT "common_program_user_user_id_fkey";

ALTER TABLE
    IF EXISTS "common_program_user"
DROP
    CONSTRAINT "common_program_user_program_id_fkey";

ALTER TABLE
    IF EXISTS "common_program_facility"
DROP
    CONSTRAINT "common_program_facility_facility_id_fkey";

ALTER TABLE
    IF EXISTS "common_program_facility"
DROP
    CONSTRAINT "common_program_facility_program_id_fkey";

ALTER TABLE
    IF EXISTS "common_auditlog"
DROP
    CONSTRAINT "common_auditlog_created_by_fkey";

ALTER TABLE
    IF EXISTS "common_auditlog"
DROP
    CONSTRAINT "common_auditlog_updated_by_fkey";

ALTER TABLE
    IF EXISTS "common_auditlog"
DROP
    CONSTRAINT "common_auditlog_organisation_id_fkey";

ALTER TABLE
    IF EXISTS "common_address"
DROP
    CONSTRAINT "common_address_created_by_fkey";

ALTER TABLE
    IF EXISTS "common_address"
DROP
    CONSTRAINT "common_address_updated_by_fkey";

ALTER TABLE
    IF EXISTS "common_address"
DROP
    CONSTRAINT "common_address_organisation_id_fkey";

ALTER TABLE
    IF EXISTS "common_securityquestion"
DROP
    CONSTRAINT "common_securityquestion_created_by_fkey";

ALTER TABLE
    IF EXISTS "common_securityquestion"
DROP
    CONSTRAINT "common_securityquestion_updated_by_fkey";

ALTER TABLE
    IF EXISTS "clients_relatedperson"
DROP
    CONSTRAINT "clients_relatedperson_created_by_fkey";

ALTER TABLE
    IF EXISTS "clients_relatedperson"
DROP
    CONSTRAINT "clients_relatedperson_updated_by_fkey";

ALTER TABLE
    IF EXISTS "clients_relatedperson"
DROP
    CONSTRAINT "clients_relatedperson_organisation_id_fkey";

ALTER TABLE
    IF EXISTS "clients_relatedperson"
DROP
    CONSTRAINT "clients_relatedperson_program_id_fkey";

ALTER TABLE
    IF EXISTS "clients_relatedperson_addresses"
DROP
    CONSTRAINT "clients_relatedperson_addresses_relatedperson_id_fkey";

ALTER TABLE
    IF EXISTS "clients_relatedperson_addresses"
DROP
    CONSTRAINT "clients_relatedperson_addresses_address_id_fkey";

ALTER TABLE
    IF EXISTS "clients_identifier"
DROP
    CONSTRAINT "clients_identifier_created_by_fkey";

ALTER TABLE
    IF EXISTS "clients_identifier"
DROP
    CONSTRAINT "clients_identifier_updated_by_fkey";

ALTER TABLE
    IF EXISTS "clients_identifier"
DROP
    CONSTRAINT "clients_identifier_organisation_id_fkey";

ALTER TABLE
    IF EXISTS "clients_identifier"
DROP
    CONSTRAINT "clients_identifier_program_id_fkey";

ALTER TABLE
    IF EXISTS "clients_healthdiaryquote"
DROP
    CONSTRAINT "clients_healthdiaryquote_created_by_fkey";

ALTER TABLE
    IF EXISTS "clients_healthdiaryquote"
DROP
    CONSTRAINT "clients_healthdiaryquote_updated_by_fkey";

ALTER TABLE
    IF EXISTS "clients_healthdiaryquote"
DROP
    CONSTRAINT "clients_healthdiaryquote_program_id_fkey";

ALTER TABLE
    IF EXISTS "clients_healthdiaryquote"
DROP
    CONSTRAINT "clients_healthdiaryquote_organisation_id_fkey";

ALTER TABLE
    IF EXISTS "authority_authoritypermission"
DROP
    CONSTRAINT "authority_authoritypermission_created_by_fkey";

ALTER TABLE
    IF EXISTS "authority_authoritypermission"
DROP
    CONSTRAINT "authority_authoritypermission_updated_by_fkey";

ALTER TABLE
    IF EXISTS "authority_authoritypermission"
DROP
    CONSTRAINT "authority_authoritypermission_organisation_id_fkey";

ALTER TABLE
    IF EXISTS "authority_authorityrole"
DROP
    CONSTRAINT "authority_authorityrole_created_by_fkey";

ALTER TABLE
    IF EXISTS "authority_authorityrole"
DROP
    CONSTRAINT "authority_authorityrole_updated_by_fkey";

ALTER TABLE
    IF EXISTS "authority_authorityrole"
DROP
    CONSTRAINT "authority_authorityrole_organisation_id_fkey";

ALTER TABLE
    IF EXISTS "authority_authorityrole_permissions"
DROP
    CONSTRAINT "authority_authorityrole_permissions_authorityrole_id_fkey";

ALTER TABLE
    IF EXISTS "authority_authorityrole_permissions"
DROP
    CONSTRAINT "authority_authorityrole_permissions_authoritypermission_id_fkey";

ALTER TABLE
    IF EXISTS "communities_community"
DROP
    CONSTRAINT "communities_community_created_by_fkey";

ALTER TABLE
    IF EXISTS "communities_community"
DROP
    CONSTRAINT "communities_community_updated_by_fkey";

ALTER TABLE
    IF EXISTS "communities_community"
DROP
    CONSTRAINT "communities_community_program_id_fkey";

ALTER TABLE
    IF EXISTS "communities_community"
DROP
    CONSTRAINT "communities_community_organisation_id_fkey";

ALTER TABLE
    IF EXISTS "screeningtools_screeningtoolsquestion"
DROP
    CONSTRAINT "screeningtools_screeningtoolsquestion_created_by_fkey";

ALTER TABLE
    IF EXISTS "screeningtools_screeningtoolsquestion"
DROP
    CONSTRAINT "screeningtools_screeningtoolsquestion_updated_by_fkey";

ALTER TABLE
    IF EXISTS "screeningtools_screeningtoolsquestion"
DROP
    CONSTRAINT "screeningtools_screeningtoolsquestion_organisation_id_fkey";

ALTER TABLE
    IF EXISTS "screeningtools_screeningtoolsquestion"
DROP
    CONSTRAINT "screeningtools_screeningtoolsquestion_program_id_fkey";

ALTER TABLE
    IF EXISTS "users_termsofservice"
DROP
    CONSTRAINT "users_termsofservice_created_by_fkey";

ALTER TABLE
    IF EXISTS "users_termsofservice"
DROP
    CONSTRAINT "users_termsofservice_updated_by_fkey";

ALTER TABLE
    IF EXISTS "staff_staff"
DROP
    CONSTRAINT "staff_staff_created_by_fkey";

ALTER TABLE
    IF EXISTS "staff_staff"
DROP
    CONSTRAINT "staff_staff_updated_by_fkey";

ALTER TABLE
    IF EXISTS "staff_staff"
DROP
    CONSTRAINT "staff_staff_current_facility_id_fkey";

ALTER TABLE
    IF EXISTS "staff_staff"
DROP
    CONSTRAINT "staff_staff_organisation_id_fkey";

ALTER TABLE
    IF EXISTS "staff_staff"
DROP
    CONSTRAINT "staff_staff_user_id_fkey";

ALTER TABLE
    IF EXISTS "staff_staff"
DROP
    CONSTRAINT "staff_staff_program_id_fkey";

ALTER TABLE
    IF EXISTS "staff_staff_facilities"
DROP
    CONSTRAINT "staff_staff_facilities_staff_id_fkey";

ALTER TABLE
    IF EXISTS "staff_staff_facilities"
ADD
    FOREIGN KEY ("facility_id") REFERENCES "common_facility" ("id");

ALTER TABLE
    IF EXISTS "clients_client"
ADD
    FOREIGN KEY ("created_by") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "clients_client"
ADD
    FOREIGN KEY ("updated_by") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "clients_client"
ADD
    FOREIGN KEY ("current_facility_id") REFERENCES "common_facility" ("id");

ALTER TABLE
    IF EXISTS "clients_client"
ADD
    FOREIGN KEY ("organisation_id") REFERENCES "common_organisation" ("id");

ALTER TABLE
    IF EXISTS "clients_client"
ADD
    FOREIGN KEY ("user_id") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "clients_client"
ADD
    FOREIGN KEY ("program_id") REFERENCES "common_program" ("id");

ALTER TABLE
    IF EXISTS "clients_relatedperson_contacts"
ADD
    FOREIGN KEY ("relatedperson_id") REFERENCES "clients_relatedperson" ("id");

ALTER TABLE
    IF EXISTS "clients_relatedperson_contacts"
ADD
    FOREIGN KEY ("contact_id") REFERENCES "common_contact" ("id");

ALTER TABLE
    IF EXISTS "clients_clientfacility"
ADD
    FOREIGN KEY ("created_by") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "clients_clientfacility"
ADD
    FOREIGN KEY ("updated_by") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "clients_clientfacility"
ADD
    FOREIGN KEY ("client_id") REFERENCES "clients_client" ("id");

ALTER TABLE
    IF EXISTS "clients_clientfacility"
ADD
    FOREIGN KEY ("facility_id") REFERENCES "common_facility" ("id");

ALTER TABLE
    IF EXISTS "clients_clientfacility"
ADD
    FOREIGN KEY ("organisation_id") REFERENCES "common_organisation" ("id");

ALTER TABLE
    IF EXISTS "clients_clientfacility"
ADD
    FOREIGN KEY ("program_id") REFERENCES "common_program" ("id");

ALTER TABLE
    IF EXISTS "clients_client_identifiers"
ADD
    FOREIGN KEY ("client_id") REFERENCES "clients_client" ("id");

ALTER TABLE
    IF EXISTS "clients_client_identifiers"
ADD
    FOREIGN KEY ("identifier_id") REFERENCES "clients_identifier" ("id");

ALTER TABLE
    IF EXISTS "clients_client_related_persons"
ADD
    FOREIGN KEY ("client_id") REFERENCES "clients_client" ("id");

ALTER TABLE
    IF EXISTS "clients_client_related_persons"
ADD
    FOREIGN KEY ("relatedperson_id") REFERENCES "clients_relatedperson" ("id");

ALTER TABLE
    IF EXISTS "common_securityquestionresponse"
ADD
    FOREIGN KEY ("created_by") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "common_securityquestionresponse"
ADD
    FOREIGN KEY ("updated_by") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "common_securityquestionresponse"
ADD
    FOREIGN KEY ("organisation_id") REFERENCES "common_organisation" ("id");

ALTER TABLE
    IF EXISTS "common_securityquestionresponse"
ADD
    FOREIGN KEY ("question_id") REFERENCES "common_securityquestion" ("id");

ALTER TABLE
    IF EXISTS "common_securityquestionresponse"
ADD
    FOREIGN KEY ("user_id") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "clients_healthdiaryentry"
ADD
    FOREIGN KEY ("created_by") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "clients_healthdiaryentry"
ADD
    FOREIGN KEY ("updated_by") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "clients_healthdiaryentry"
ADD
    FOREIGN KEY ("client_id") REFERENCES "clients_client" ("id");

ALTER TABLE
    IF EXISTS "clients_healthdiaryentry"
ADD
    FOREIGN KEY ("organisation_id") REFERENCES "common_organisation" ("id");

ALTER TABLE
    IF EXISTS "clients_healthdiaryentry"
ADD
    FOREIGN KEY ("program_id") REFERENCES "common_program" ("id");

ALTER TABLE
    IF EXISTS "clients_servicerequest"
ADD
    FOREIGN KEY ("created_by") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "clients_servicerequest"
ADD
    FOREIGN KEY ("updated_by") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "clients_servicerequest"
ADD
    FOREIGN KEY ("in_progress_by_id") REFERENCES "staff_staff" ("id");

ALTER TABLE
    IF EXISTS "clients_servicerequest"
ADD
    FOREIGN KEY ("resolved_by_id") REFERENCES "staff_staff" ("id");

ALTER TABLE
    IF EXISTS "clients_servicerequest"
ADD
    FOREIGN KEY ("facility_id") REFERENCES "common_facility" ("id");

ALTER TABLE
    IF EXISTS "clients_servicerequest"
ADD
    FOREIGN KEY ("client_id") REFERENCES "clients_client" ("id");

ALTER TABLE
    IF EXISTS "clients_servicerequest"
ADD
    FOREIGN KEY ("organisation_id") REFERENCES "common_organisation" ("id");

ALTER TABLE
    IF EXISTS "clients_servicerequest"
ADD
    FOREIGN KEY ("program_id") REFERENCES "common_program" ("id");

ALTER TABLE
    IF EXISTS "staff_staff_identifiers"
ADD
    FOREIGN KEY ("staff_id") REFERENCES "staff_staff" ("id");

ALTER TABLE
    IF EXISTS "staff_staff_identifiers"
ADD
    FOREIGN KEY ("identifier_id") REFERENCES "clients_identifier" ("id");

ALTER TABLE
    IF EXISTS "appointments_appointment"
ADD
    FOREIGN KEY ("created_by") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "appointments_appointment"
ADD
    FOREIGN KEY ("updated_by") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "appointments_appointment"
ADD
    FOREIGN KEY ("client_id") REFERENCES "clients_client" ("id");

ALTER TABLE
    IF EXISTS "appointments_appointment"
ADD
    FOREIGN KEY ("organisation_id") REFERENCES "common_organisation" ("id");

ALTER TABLE
    IF EXISTS "appointments_appointment"
ADD
    FOREIGN KEY ("staff_id") REFERENCES "staff_staff" ("id");

ALTER TABLE
    IF EXISTS "appointments_appointment"
ADD
    FOREIGN KEY ("facility_id") REFERENCES "common_facility" ("id");

ALTER TABLE
    IF EXISTS "appointments_appointment"
ADD
    FOREIGN KEY ("program_id") REFERENCES "common_program" ("id");

ALTER TABLE
    IF EXISTS "authority_authorityrole_users"
ADD
    FOREIGN KEY ("authorityrole_id") REFERENCES "authority_authorityrole" ("id");

ALTER TABLE
    IF EXISTS "authority_authorityrole_users"
ADD
    FOREIGN KEY ("user_id") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "common_notification"
ADD
    FOREIGN KEY ("created_by") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "common_notification"
ADD
    FOREIGN KEY ("updated_by") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "common_notification"
ADD
    FOREIGN KEY ("facility_id") REFERENCES "common_facility" ("id");

ALTER TABLE
    IF EXISTS "common_notification"
ADD
    FOREIGN KEY ("organisation_id") REFERENCES "common_organisation" ("id");

ALTER TABLE
    IF EXISTS "common_notification"
ADD
    FOREIGN KEY ("user_id") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "common_notification"
ADD
    FOREIGN KEY ("program_id") REFERENCES "common_program" ("id");

ALTER TABLE
    IF EXISTS "common_usersurveys"
ADD
    FOREIGN KEY ("created_by") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "common_usersurveys"
ADD
    FOREIGN KEY ("updated_by") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "common_usersurveys"
ADD
    FOREIGN KEY ("organisation_id") REFERENCES "common_organisation" ("id");

ALTER TABLE
    IF EXISTS "common_usersurveys"
ADD
    FOREIGN KEY ("user_id") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "common_usersurveys"
ADD
    FOREIGN KEY ("program_id") REFERENCES "common_program" ("id");

ALTER TABLE
    IF EXISTS "common_feedback"
ADD
    FOREIGN KEY ("created_by") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "common_feedback"
ADD
    FOREIGN KEY ("updated_by") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "common_feedback"
ADD
    FOREIGN KEY ("organisation_id") REFERENCES "common_organisation" ("id");

ALTER TABLE
    IF EXISTS "common_feedback"
ADD
    FOREIGN KEY ("user_id") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "common_feedback"
ADD
    FOREIGN KEY ("program_id") REFERENCES "common_program" ("id");

ALTER TABLE
    IF EXISTS "communities_community_clients"
ADD
    FOREIGN KEY ("community_id") REFERENCES "communities_community" ("id");

ALTER TABLE
    IF EXISTS "communities_community_clients"
ADD
    FOREIGN KEY ("client_id") REFERENCES "clients_client" ("id");

ALTER TABLE
    IF EXISTS "communities_community_staff"
ADD
    FOREIGN KEY ("community_id") REFERENCES "communities_community" ("id");

ALTER TABLE
    IF EXISTS "communities_community_staff"
ADD
    FOREIGN KEY ("staff_id") REFERENCES "staff_staff" ("id");

ALTER TABLE
    IF EXISTS "screeningtools_screeningtoolsresponse"
ADD
    FOREIGN KEY ("created_by") REFERENCES "users_user" ("id");

ALTER TABLE 
    IF EXISTS "screeningtools_screeningtoolsresponse"
DROP
    CONSTRAINT "screeningtools_screeningtoolsresponse_updated_by_fkey";

ALTER TABLE 
    IF EXISTS "screeningtools_screeningtoolsresponse"
DROP
    CONSTRAINT "screeningtools_screeningtoolsresponse_client_id_fkey";

ALTER TABLE 
    IF EXISTS "screeningtools_screeningtoolsresponse"
DROP
    CONSTRAINT "screeningtools_screeningtoolsresponse_organisation_id_fkey";

ALTER TABLE 
    IF EXISTS "screeningtools_screeningtoolsresponse"
DROP
    CONSTRAINT "screeningtools_screeningtoolsresponse_question_id_fkey";

ALTER TABLE 
    IF EXISTS "screeningtools_screeningtoolsresponse"
DROP
    CONSTRAINT "screeningtools_screeningtoolsresponse_program_id_fkey";

ALTER TABLE 
    IF EXISTS "staff_servicerequest"
DROP
    CONSTRAINT "staff_servicerequest_created_by_fkey";

ALTER TABLE 
    IF EXISTS "staff_servicerequest"
DROP
    CONSTRAINT "staff_servicerequest_updated_by_fkey";

ALTER TABLE 
    IF EXISTS "staff_servicerequest"
DROP
    CONSTRAINT "staff_servicerequest_organisation_id_fkey";

ALTER TABLE 
    IF EXISTS "staff_servicerequest"
DROP
    CONSTRAINT "staff_servicerequest_resolved_by_id_fkey";

ALTER TABLE 
    IF EXISTS "staff_servicerequest"
DROP
    CONSTRAINT "staff_servicerequest_staff_id_fkey";

ALTER TABLE 
    IF EXISTS "staff_servicerequest"
DROP
    CONSTRAINT "staff_servicerequest_facility_id_fkey";

ALTER TABLE 
    IF EXISTS "staff_servicerequest"
DROP
    CONSTRAINT "staff_servicerequest_program_id_fkey";

ALTER TABLE 
    IF EXISTS "users_metric"
DROP
    CONSTRAINT "users_metric_created_by_fkey";

ALTER TABLE 
    IF EXISTS "users_metric"
DROP
    CONSTRAINT "users_metric_updated_by_fkey";

ALTER TABLE 
    IF EXISTS "users_metric"
DROP
    CONSTRAINT "users_metric_organisation_id_fkey";

ALTER TABLE 
    IF EXISTS "users_userpin"
DROP
    CONSTRAINT "users_userpin_created_by_fkey";

ALTER TABLE 
    IF EXISTS "users_userpin"
DROP
    CONSTRAINT "users_userpin_updated_by_fkey";

ALTER TABLE 
    IF EXISTS "users_userpin"
DROP
    CONSTRAINT "users_userpin_user_id_fkey";

ALTER TABLE 
    IF EXISTS "users_userpin"
DROP
    CONSTRAINT "users_userpin_organisation_id_fkey";

ALTER TABLE 
    IF EXISTS "users_userotp"
DROP
    CONSTRAINT "users_userotp_user_id_fkey";

ALTER TABLE 
    IF EXISTS "users_userotp"
DROP
    CONSTRAINT "users_userotp_created_by_fkey";

ALTER TABLE 
    IF EXISTS "users_userotp"
DROP
    CONSTRAINT "users_userotp_updated_by_fkey";

ALTER TABLE 
    IF EXISTS "questionnaires_questionnaire"
DROP
    CONSTRAINT "questionnaires_questionnaire_created_by_fkey";

ALTER TABLE 
    IF EXISTS "questionnaires_questionnaire"
DROP
    CONSTRAINT "questionnaires_questionnaire_updated_by_fkey";

ALTER TABLE 
    IF EXISTS "questionnaires_questionnaire"
DROP
    CONSTRAINT "questionnaires_questionnaire_organisation_id_fkey";

ALTER TABLE 
    IF EXISTS "questionnaires_questionnaire"
DROP
    CONSTRAINT "questionnaires_questionnaire_program_id_fkey";

ALTER TABLE 
    IF EXISTS "questionnaires_screeningtool"
DROP
    CONSTRAINT "questionnaires_screeningtool_created_by_fkey";

ALTER TABLE 
    IF EXISTS "questionnaires_screeningtool"
DROP
    CONSTRAINT "questionnaires_screeningtool_updated_by_fkey";

ALTER TABLE 
    IF EXISTS "questionnaires_screeningtool"
DROP
    CONSTRAINT "questionnaires_screeningtool_organisation_id_fkey";

ALTER TABLE 
    IF EXISTS "questionnaires_screeningtool"
DROP
    CONSTRAINT "questionnaires_screeningtool_questionnaire_id_fkey";

ALTER TABLE 
    IF EXISTS "questionnaires_screeningtool"
DROP
    CONSTRAINT "questionnaires_screeningtool_program_id_fkey";

ALTER TABLE 
    IF EXISTS "questionnaires_question"
DROP
    CONSTRAINT "questionnaires_question_created_by_fkey";

ALTER TABLE 
    IF EXISTS "questionnaires_question"
DROP
    CONSTRAINT "questionnaires_question_updated_by_fkey";

ALTER TABLE 
    IF EXISTS "questionnaires_question"
DROP
    CONSTRAINT "questionnaires_question_organisation_id_fkey";

ALTER TABLE 
    IF EXISTS "questionnaires_question"
DROP
    CONSTRAINT "questionnaires_question_questionnaire_id_fkey";

ALTER TABLE 
    IF EXISTS "questionnaires_question"
DROP
    CONSTRAINT "questionnaires_question_program_id_fkey";

ALTER TABLE 
    IF EXISTS "questionnaires_questioninputchoice"
DROP
    CONSTRAINT "questionnaires_questioninputchoice_created_by_fkey";

ALTER TABLE 
    IF EXISTS "questionnaires_questioninputchoice"
DROP
    CONSTRAINT "questionnaires_questioninputchoice_updated_by_fkey";

ALTER TABLE 
    IF EXISTS "questionnaires_questioninputchoice"
DROP
    CONSTRAINT "questionnaires_questioninputchoice_organisation_id_fkey";

ALTER TABLE 
    IF EXISTS "questionnaires_questioninputchoice"
DROP
    CONSTRAINT "questionnaires_questioninputchoice_question_id_fkey";

ALTER TABLE 
    IF EXISTS "questionnaires_questioninputchoice"
DROP
    CONSTRAINT "questionnaires_questioninputchoice_program_id_fkey";

ALTER TABLE 
    IF EXISTS "questionnaires_screeningtoolresponse"
DROP
    CONSTRAINT "questionnaires_screeningtoolresponse_created_by_fkey";

ALTER TABLE 
    IF EXISTS "questionnaires_screeningtoolresponse"
DROP
    CONSTRAINT "questionnaires_screeningtoolresponse_updated_by_fkey";

ALTER TABLE 
    IF EXISTS "questionnaires_screeningtoolresponse"
DROP
    CONSTRAINT "questionnaires_screeningtoolresponse_organisation_id_fkey";

ALTER TABLE 
    IF EXISTS "questionnaires_screeningtoolresponse"
DROP
    CONSTRAINT "questionnaires_screeningtoolresponse_screeningtool_id_fkey";

ALTER TABLE 
    IF EXISTS "questionnaires_screeningtoolresponse"
DROP
    CONSTRAINT "questionnaires_screeningtoolresponse_facility_id_fkey";

ALTER TABLE 
    IF EXISTS "questionnaires_screeningtoolresponse"
DROP
    CONSTRAINT "questionnaires_screeningtoolresponse_client_id_fkey";

ALTER TABLE 
    IF EXISTS "questionnaires_screeningtoolresponse"
DROP
    CONSTRAINT "questionnaires_screeningtoolresponse_program_id_fkey";

ALTER TABLE 
    IF EXISTS "questionnaires_screeningtoolquestionresponse"
DROP
    CONSTRAINT "questionnaires_screeningtoolquestionresponse_created_by_fkey";

ALTER TABLE 
    IF EXISTS "questionnaires_screeningtoolquestionresponse"
DROP
    CONSTRAINT "questionnaires_screeningtoolquestionresponse_updated_by_fkey";

ALTER TABLE 
    IF EXISTS "questionnaires_screeningtoolquestionresponse"
DROP
    CONSTRAINT "questionnaires_screeningtoolquestionresponse_organisation_id_fkey";

ALTER TABLE 
    IF EXISTS "questionnaires_screeningtoolquestionresponse"
DROP
    CONSTRAINT "questionnaires_screeningtoolquestionresponse_screeningtoolresponse_id_fkey";

ALTER TABLE 
    IF EXISTS "questionnaires_screeningtoolquestionresponse"
DROP
   CONSTRAINT "questionnaires_screeningtoolquestionresponse_question_id_fkey";

ALTER TABLE 
    IF EXISTS "questionnaires_screeningtoolquestionresponse"
DROP
   CONSTRAINT "questionnaires_screeningtoolquestionresponse_program_id_fkey";

ALTER TABLE 
    IF EXISTS "questionnaires_screeningtoolquestionresponse"
DROP
    CONSTRAINT "questionnaires_screeningtoolquestionresponse_facility_id_fkey";

ALTER TABLE 
    IF EXISTS "clients_client_facilities"
DROP
    CONSTRAINT "clients_client_facilities_client_id_fkey";

ALTER TABLE 
    IF EXISTS "clients_client_facilities"
DROP
    CONSTRAINT "clients_client_facilities_facility_id_fkey";

ALTER TABLE 
    IF EXISTS "caregivers_caregiver"
DROP
    CONSTRAINT "caregivers_caregiver_created_by_fkey";

ALTER TABLE 
    IF EXISTS "caregivers_caregiver"
DROP
    CONSTRAINT "caregivers_caregiver_updated_by_fkey";

ALTER TABLE 
    IF EXISTS "caregivers_caregiver"
DROP
    CONSTRAINT "caregivers_caregiver_organisation_id_fkey";

ALTER TABLE 
    IF EXISTS "caregivers_caregiver"
DROP
    CONSTRAINT "caregivers_caregiver_user_id_fkey";

ALTER TABLE 
    IF EXISTS "caregivers_caregiver"
DROP
    CONSTRAINT "caregivers_caregiver_program_id_fkey";

ALTER TABLE 
    IF EXISTS "caregivers_caregiver_client"
DROP
    CONSTRAINT "caregivers_caregiver_client_caregiver_id_fkey";

ALTER TABLE 
    IF EXISTS "caregivers_caregiver_client"
DROP
    CONSTRAINT "caregivers_caregiver_client_client_id_fkey";

ALTER TABLE 
    IF EXISTS "caregivers_caregiver_client"
DROP
    CONSTRAINT "caregivers_caregiver_client_created_by_fkey";

ALTER TABLE 
    IF EXISTS "caregivers_caregiver_client"
DROP
    CONSTRAINT "caregivers_caregiver_client_updated_by_fkey";

ALTER TABLE 
    IF EXISTS "caregivers_caregiver_client"
DROP
    CONSTRAINT "caregivers_caregiver_client_organisation_id_fkey";

ALTER TABLE 
    IF EXISTS "caregivers_caregiver_client"
DROP
    CONSTRAINT "caregivers_caregiver_client_assigned_by_fkey";

ALTER TABLE 
    IF EXISTS "caregivers_caregiver_client"
DROP
    CONSTRAINT "caregivers_caregiver_client_program_id_fkey";

DROP TABLE IF EXISTS "caregivers_caregiver_client";

DROP TABLE IF EXISTS "caregivers_caregiver";

DROP TABLE IF EXISTS "clients_client_facilities";

DROP TABLE IF EXISTS "questionnaires_screeningtoolquestionresponse";

DROP TABLE IF EXISTS "questionnaires_screeningtoolresponse";

DROP TABLE IF EXISTS "questionnaires_questioninputchoice";

DROP TABLE IF EXISTS "questionnaires_question";

DROP TABLE IF EXISTS "questionnaires_screeningtool";

DROP TABLE IF EXISTS "questionnaires_questionnaire";

DROP TABLE IF EXISTS "users_userotp";

DROP TABLE IF EXISTS "users_userpin";

DROP TABLE IF EXISTS "users_metric";

DROP TABLE IF EXISTS "staff_servicerequest";

DROP TABLE IF EXISTS "screeningtools_screeningtoolsresponse";

DROP TABLE IF EXISTS "communities_community_staff";

DROP TABLE IF EXISTS "communities_community_clients";

DROP TABLE IF EXISTS "common_feedback";

DROP TABLE IF EXISTS "common_usersurveys";

DROP TABLE IF EXISTS "common_notification";

DROP TABLE IF EXISTS "authority_authorityrole_users";

DROP TABLE IF EXISTS "appointments_appointment";

DROP TABLE IF EXISTS "staff_staff_identifiers";

DROP TABLE IF EXISTS "clients_servicerequest";

DROP TABLE IF EXISTS "clients_healthdiaryentry";

DROP TABLE IF EXISTS "common_securityquestionresponse";

DROP TABLE IF EXISTS "clients_client_related_persons";

DROP TABLE IF EXISTS "clients_client_identifiers";

DROP TABLE IF EXISTS "clients_clientfacility";

DROP TABLE IF EXISTS "clients_relatedperson_contacts";

DROP TABLE IF EXISTS "clients_client";

DROP TABLE IF EXISTS "staff_staff_facilities";

DROP TABLE IF EXISTS "staff_staff";

DROP TABLE IF EXISTS "users_termsofservice";

DROP TABLE IF EXISTS "screeningtools_screeningtoolsquestion";

DROP TABLE IF EXISTS "communities_community";

DROP TABLE IF EXISTS "authority_authorityrole_permissions";

DROP TABLE IF EXISTS "authority_authorityrole";

DROP TABLE IF EXISTS "authority_authoritypermission";

DROP TABLE IF EXISTS "clients_healthdiaryquote";

DROP TABLE IF EXISTS "clients_identifier";

DROP TABLE IF EXISTS "clients_relatedperson_DROPresses";

DROP TABLE IF EXISTS "clients_relatedperson";

DROP TABLE IF EXISTS "common_securityquestion";

DROP TABLE IF EXISTS "common_DROPress";

DROP TABLE IF EXISTS "common_auditlog";

DROP TABLE IF EXISTS "common_program_facility";

DROP TABLE IF EXISTS "common_program_user";

DROP TABLE IF EXISTS "common_program";

DROP TABLE IF EXISTS "common_contact";

DROP TABLE IF EXISTS "users_user";

DROP TABLE IF EXISTS "common_facility_identifier";

DROP TABLE IF EXISTS "common_facility";

DROP TABLE IF EXISTS "common_organisation";

COMMIT;