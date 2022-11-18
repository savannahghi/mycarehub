BEGIN;

ALTER TABLE
    IF EXISTS "common_organisation"
    DROP CONSTRAINT IF EXISTS "common_organisation_created_by_fkey";

ALTER TABLE
    IF EXISTS "common_organisation"
    DROP CONSTRAINT IF EXISTS "common_organisation_updated_by_fkey";

ALTER TABLE
    IF EXISTS "common_facility"
    DROP CONSTRAINT IF EXISTS "common_facility_created_by_fkey";

ALTER TABLE
    IF EXISTS "common_facility"
    DROP CONSTRAINT IF EXISTS "common_facility_updated_by_fkey";

ALTER TABLE
    IF EXISTS "common_facility_identifier"
    DROP CONSTRAINT IF EXISTS "common_facility_identifier_created_by_fkey";

ALTER TABLE
    IF EXISTS "common_facility_identifier"
    DROP CONSTRAINT IF EXISTS "common_facility_identifier_updated_by_fkey";

ALTER TABLE
    IF EXISTS "common_facility_identifier"
    DROP CONSTRAINT IF EXISTS "common_facility_identifier_facility_id_fkey";

ALTER TABLE
    IF EXISTS "users_user"
    DROP CONSTRAINT IF EXISTS "users_user_created_by_fkey";

ALTER TABLE
    IF EXISTS "users_user"
    DROP CONSTRAINT IF EXISTS "users_user_updated_by_fkey";

ALTER TABLE
    IF EXISTS "users_user"
    DROP CONSTRAINT IF EXISTS "users_user_accepted_terms_of_service_id_fkey";

ALTER TABLE
    IF EXISTS "users_user"
    DROP CONSTRAINT IF EXISTS "users_user_current_program_id_fkey";

ALTER TABLE
    IF EXISTS "users_user"
    DROP CONSTRAINT IF EXISTS "users_user_organisation_id_fkey";

ALTER TABLE
    IF EXISTS "common_contact"
    DROP CONSTRAINT IF EXISTS "common_contact_created_by_fkey";

ALTER TABLE
    IF EXISTS "common_contact"
    DROP CONSTRAINT IF EXISTS "common_contact_updated_by_fkey";

ALTER TABLE
    IF EXISTS "common_contact"
    DROP CONSTRAINT IF EXISTS "common_contact_user_id_fkey";

ALTER TABLE
    IF EXISTS "common_contact"
    DROP CONSTRAINT IF EXISTS "common_contact_organisation_id_fkey";

ALTER TABLE
    IF EXISTS "common_program"
    DROP CONSTRAINT IF EXISTS "common_program_created_by_fkey";

ALTER TABLE
    IF EXISTS "common_program"
    DROP CONSTRAINT IF EXISTS "common_program_updated_by_fkey";

ALTER TABLE
    IF EXISTS "common_program"
    DROP CONSTRAINT IF EXISTS "common_program_organisation_id_fkey";

ALTER TABLE
    IF EXISTS "common_program_user"
    DROP CONSTRAINT IF EXISTS "common_program_user_user_id_fkey";

ALTER TABLE
    IF EXISTS "common_program_user"
    DROP CONSTRAINT IF EXISTS "common_program_user_program_id_fkey";

ALTER TABLE
    IF EXISTS "common_program_facility"
    DROP CONSTRAINT IF EXISTS "common_program_facility_facility_id_fkey";

ALTER TABLE
    IF EXISTS "common_program_facility"
    DROP CONSTRAINT IF EXISTS "common_program_facility_program_id_fkey";

ALTER TABLE
    IF EXISTS "common_auditlog"
    DROP CONSTRAINT IF EXISTS "common_auditlog_created_by_fkey";

ALTER TABLE
    IF EXISTS "common_auditlog"
    DROP CONSTRAINT IF EXISTS "common_auditlog_updated_by_fkey";

ALTER TABLE
    IF EXISTS "common_auditlog"
    DROP CONSTRAINT IF EXISTS "common_auditlog_organisation_id_fkey";

ALTER TABLE
    IF EXISTS "common_address"
    DROP CONSTRAINT IF EXISTS "common_address_created_by_fkey";

ALTER TABLE
    IF EXISTS "common_address"
    DROP CONSTRAINT IF EXISTS "common_address_updated_by_fkey";

ALTER TABLE
    IF EXISTS "common_address"
    DROP CONSTRAINT IF EXISTS "common_address_organisation_id_fkey";

ALTER TABLE
    IF EXISTS "common_securityquestion"
    DROP CONSTRAINT IF EXISTS "common_securityquestion_created_by_fkey";

ALTER TABLE
    IF EXISTS "common_securityquestion"
    DROP CONSTRAINT IF EXISTS "common_securityquestion_updated_by_fkey";

ALTER TABLE
    IF EXISTS "clients_relatedperson"
    DROP CONSTRAINT IF EXISTS "clients_relatedperson_created_by_fkey";

ALTER TABLE
    IF EXISTS "clients_relatedperson"
    DROP CONSTRAINT IF EXISTS "clients_relatedperson_updated_by_fkey";

ALTER TABLE
    IF EXISTS "clients_relatedperson"
    DROP CONSTRAINT IF EXISTS "clients_relatedperson_organisation_id_fkey";

ALTER TABLE
    IF EXISTS "clients_relatedperson"
    DROP CONSTRAINT IF EXISTS "clients_relatedperson_program_id_fkey";

ALTER TABLE
    IF EXISTS "clients_relatedperson_addresses"
    DROP CONSTRAINT IF EXISTS "clients_relatedperson_addresses_relatedperson_id_fkey";

ALTER TABLE
    IF EXISTS "clients_relatedperson_addresses"
    DROP CONSTRAINT IF EXISTS "clients_relatedperson_addresses_address_id_fkey";

ALTER TABLE
    IF EXISTS "clients_identifier"
    DROP CONSTRAINT IF EXISTS "clients_identifier_created_by_fkey";

ALTER TABLE
    IF EXISTS "clients_identifier"
    DROP CONSTRAINT IF EXISTS "clients_identifier_updated_by_fkey";

ALTER TABLE
    IF EXISTS "clients_identifier"
    DROP CONSTRAINT IF EXISTS "clients_identifier_organisation_id_fkey";

ALTER TABLE
    IF EXISTS "clients_identifier"
    DROP CONSTRAINT IF EXISTS "clients_identifier_program_id_fkey";

ALTER TABLE
    IF EXISTS "clients_healthdiaryquote"
    DROP CONSTRAINT IF EXISTS "clients_healthdiaryquote_created_by_fkey";

ALTER TABLE
    IF EXISTS "clients_healthdiaryquote"
    DROP CONSTRAINT IF EXISTS "clients_healthdiaryquote_updated_by_fkey";

ALTER TABLE
    IF EXISTS "clients_healthdiaryquote"
    DROP CONSTRAINT IF EXISTS "clients_healthdiaryquote_program_id_fkey";

ALTER TABLE
    IF EXISTS "clients_healthdiaryquote"
    DROP CONSTRAINT IF EXISTS "clients_healthdiaryquote_organisation_id_fkey";

ALTER TABLE
    IF EXISTS "authority_authoritypermission"
    DROP CONSTRAINT IF EXISTS "authority_authoritypermission_created_by_fkey";

ALTER TABLE
    IF EXISTS "authority_authoritypermission"
    DROP CONSTRAINT IF EXISTS "authority_authoritypermission_updated_by_fkey";

ALTER TABLE
    IF EXISTS "authority_authoritypermission"
    DROP CONSTRAINT IF EXISTS "authority_authoritypermission_organisation_id_fkey";

ALTER TABLE
    IF EXISTS "authority_authorityrole"
    DROP CONSTRAINT IF EXISTS "authority_authorityrole_created_by_fkey";

ALTER TABLE
    IF EXISTS "authority_authorityrole"
    DROP CONSTRAINT IF EXISTS "authority_authorityrole_updated_by_fkey";

ALTER TABLE
    IF EXISTS "authority_authorityrole"
    DROP CONSTRAINT IF EXISTS "authority_authorityrole_organisation_id_fkey";

ALTER TABLE
    IF EXISTS "authority_authorityrole_permissions"
    DROP CONSTRAINT IF EXISTS "authority_authorityrole_permissions_authorityrole_id_fkey";

ALTER TABLE
    IF EXISTS "authority_authorityrole_permissions"
    DROP CONSTRAINT IF EXISTS "authority_authorityrole_permissions_authoritypermission_id_fkey";

ALTER TABLE
    IF EXISTS "communities_community"
    DROP CONSTRAINT IF EXISTS "communities_community_created_by_fkey";

ALTER TABLE
    IF EXISTS "communities_community"
    DROP CONSTRAINT IF EXISTS "communities_community_updated_by_fkey";

ALTER TABLE
    IF EXISTS "communities_community"
    DROP CONSTRAINT IF EXISTS "communities_community_program_id_fkey";

ALTER TABLE
    IF EXISTS "communities_community"
    DROP CONSTRAINT IF EXISTS "communities_community_organisation_id_fkey";

ALTER TABLE
    IF EXISTS "screeningtools_screeningtoolsquestion"
    DROP CONSTRAINT IF EXISTS "screeningtools_screeningtoolsquestion_created_by_fkey";

ALTER TABLE
    IF EXISTS "screeningtools_screeningtoolsquestion"
    DROP CONSTRAINT IF EXISTS "screeningtools_screeningtoolsquestion_updated_by_fkey";

ALTER TABLE
    IF EXISTS "screeningtools_screeningtoolsquestion"
    DROP CONSTRAINT IF EXISTS "screeningtools_screeningtoolsquestion_organisation_id_fkey";

ALTER TABLE
    IF EXISTS "screeningtools_screeningtoolsquestion"
    DROP CONSTRAINT IF EXISTS "screeningtools_screeningtoolsquestion_program_id_fkey";

ALTER TABLE
    IF EXISTS "users_termsofservice"
    DROP CONSTRAINT IF EXISTS "users_termsofservice_created_by_fkey";

ALTER TABLE
    IF EXISTS "users_termsofservice"
    DROP CONSTRAINT IF EXISTS "users_termsofservice_updated_by_fkey";

ALTER TABLE
    IF EXISTS "staff_staff"
    DROP CONSTRAINT IF EXISTS "staff_staff_created_by_fkey";

ALTER TABLE
    IF EXISTS "staff_staff"
    DROP CONSTRAINT IF EXISTS "staff_staff_updated_by_fkey";

ALTER TABLE
    IF EXISTS "staff_staff"
    DROP CONSTRAINT IF EXISTS "staff_staff_current_facility_id_fkey";

ALTER TABLE
    IF EXISTS "staff_staff"
    DROP CONSTRAINT IF EXISTS "staff_staff_organisation_id_fkey";

ALTER TABLE
    IF EXISTS "staff_staff"
    DROP CONSTRAINT IF EXISTS "staff_staff_user_id_fkey";

ALTER TABLE
    IF EXISTS "staff_staff"
    DROP CONSTRAINT IF EXISTS "staff_staff_program_id_fkey";

ALTER TABLE
    IF EXISTS "staff_staff_facilities"
    DROP CONSTRAINT IF EXISTS "staff_staff_facilities_staff_id_fkey";

ALTER TABLE
    IF EXISTS "staff_staff_facilities"
    DROP CONSTRAINT IF EXISTS "staff_staff_facilities_facility_id_fkey";

ALTER TABLE
    IF EXISTS "clients_client"
    DROP CONSTRAINT IF EXISTS "clients_client_created_by_fkey";

ALTER TABLE
    IF EXISTS "clients_client"
    DROP CONSTRAINT IF EXISTS "clients_client_updated_by_fkey";

ALTER TABLE
    IF EXISTS "clients_client"
    DROP CONSTRAINT IF EXISTS "clients_client_current_facility_id_fkey";

ALTER TABLE
    IF EXISTS "clients_client"
    DROP CONSTRAINT IF EXISTS "clients_client_organisation_id_fkey";

ALTER TABLE
    IF EXISTS "clients_client"
    DROP CONSTRAINT IF EXISTS "clients_client_user_id_fkey";

ALTER TABLE
    IF EXISTS "clients_client"
    DROP CONSTRAINT IF EXISTS "clients_client_program_id_fkey";

ALTER TABLE
    IF EXISTS "clients_relatedperson_contacts"
    DROP CONSTRAINT IF EXISTS "clients_relatedperson_contacts_relatedperson_id_fkey";

ALTER TABLE
    IF EXISTS "clients_relatedperson_contacts"
    DROP CONSTRAINT IF EXISTS "clients_relatedperson_contacts_contact_id_fkey";

ALTER TABLE
    IF EXISTS "clients_clientfacility"
    DROP CONSTRAINT IF EXISTS "clients_clientfacility_created_by_fkey";

ALTER TABLE
    IF EXISTS "clients_clientfacility"
    DROP CONSTRAINT IF EXISTS "clients_clientfacility_updated_by_fkey";

ALTER TABLE
    IF EXISTS "clients_clientfacility"
    DROP CONSTRAINT IF EXISTS "clients_clientfacility_client_id_fkey";

ALTER TABLE
    IF EXISTS "clients_clientfacility"
    DROP CONSTRAINT IF EXISTS "clients_clientfacility_facility_id_fkey";

ALTER TABLE
    IF EXISTS "clients_clientfacility"
    DROP CONSTRAINT IF EXISTS "clients_clientfacility_organisation_id_fkey";

ALTER TABLE
    IF EXISTS "clients_clientfacility"
    DROP CONSTRAINT IF EXISTS "clients_clientfacility_program_id_fkey";

ALTER TABLE
    IF EXISTS "clients_client_identifiers"
    DROP CONSTRAINT IF EXISTS "clients_client_identifiers_client_id_fkey";

ALTER TABLE
    IF EXISTS "clients_client_identifiers"
    DROP CONSTRAINT IF EXISTS "clients_client_identifiers_identifier_id_fkey";

ALTER TABLE
    IF EXISTS "clients_client_related_persons"
    DROP CONSTRAINT IF EXISTS "clients_client_related_persons_client_id_fkey";

ALTER TABLE
    IF EXISTS "clients_client_related_persons"
    DROP CONSTRAINT IF EXISTS "clients_client_related_persons_relatedperson_id_fkey";

ALTER TABLE
    IF EXISTS "common_securityquestionresponse"
    DROP CONSTRAINT IF EXISTS "common_securityquestionresponse_created_by_fkey";

ALTER TABLE
    IF EXISTS "common_securityquestionresponse"
    DROP CONSTRAINT IF EXISTS "common_securityquestionresponse_updated_by_fkey";

ALTER TABLE
    IF EXISTS "common_securityquestionresponse"
    DROP CONSTRAINT IF EXISTS "common_securityquestionresponse_organisation_id_fkey";

ALTER TABLE
    IF EXISTS "common_securityquestionresponse"
    DROP CONSTRAINT IF EXISTS "common_securityquestionresponse_question_id_fkey";

ALTER TABLE
    IF EXISTS "common_securityquestionresponse"
    DROP CONSTRAINT IF EXISTS "common_securityquestionresponse_user_id_fkey";

ALTER TABLE
    IF EXISTS "clients_healthdiaryentry"
    DROP CONSTRAINT IF EXISTS "clients_healthdiaryentry_created_by_fkey";

ALTER TABLE
    IF EXISTS "clients_healthdiaryentry"
    DROP CONSTRAINT IF EXISTS "clients_healthdiaryentry_updated_by_fkey";

ALTER TABLE
    IF EXISTS "clients_healthdiaryentry"
    DROP CONSTRAINT IF EXISTS "clients_healthdiaryentry_client_id_fkey";

ALTER TABLE
    IF EXISTS "clients_healthdiaryentry"
    DROP CONSTRAINT IF EXISTS "clients_healthdiaryentry_organisation_id_fkey";

ALTER TABLE
    IF EXISTS "clients_healthdiaryentry"
    DROP CONSTRAINT IF EXISTS "clients_healthdiaryentry_program_id_fkey";

ALTER TABLE
    IF EXISTS "clients_servicerequest"
    DROP CONSTRAINT IF EXISTS "clients_servicerequest_created_by_fkey";

ALTER TABLE
    IF EXISTS "clients_servicerequest"
    DROP CONSTRAINT IF EXISTS "clients_servicerequest_updated_by_fkey";

ALTER TABLE
    IF EXISTS "clients_servicerequest"
    DROP CONSTRAINT IF EXISTS "clients_servicerequest_in_progress_by_id_fkey";

ALTER TABLE
    IF EXISTS "clients_servicerequest"
    DROP CONSTRAINT IF EXISTS "clients_servicerequest_resolved_by_id_fkey";

ALTER TABLE
    IF EXISTS "clients_servicerequest"
    DROP CONSTRAINT IF EXISTS "clients_servicerequest_facility_id_fkey";

ALTER TABLE
    IF EXISTS "clients_servicerequest"
    DROP CONSTRAINT IF EXISTS "clients_servicerequest_client_id_fkey";

ALTER TABLE
    IF EXISTS "clients_servicerequest"
    DROP CONSTRAINT IF EXISTS "clients_servicerequest_organisation_id_fkey";

ALTER TABLE
    IF EXISTS "clients_servicerequest"
    DROP CONSTRAINT IF EXISTS "clients_servicerequest_program_id_fkey";

ALTER TABLE
    IF EXISTS "staff_staff_identifiers"
    DROP CONSTRAINT IF EXISTS "staff_staff_identifiers_staff_id_fkey";

ALTER TABLE
    IF EXISTS "staff_staff_identifiers"
    DROP CONSTRAINT IF EXISTS "staff_staff_identifiers_identifier_id_fkey";

ALTER TABLE
    IF EXISTS "appointments_appointment"
    DROP CONSTRAINT IF EXISTS "appointments_appointment_created_by_fkey";

ALTER TABLE
    IF EXISTS "appointments_appointment"
    DROP CONSTRAINT IF EXISTS "appointments_appointment_updated_by_fkey";

ALTER TABLE
    IF EXISTS "appointments_appointment"
    DROP CONSTRAINT IF EXISTS "appointments_appointment_client_id_fkey";

ALTER TABLE
    IF EXISTS "appointments_appointment"
    DROP CONSTRAINT IF EXISTS "appointments_appointment_organisation_id_fkey";

ALTER TABLE
    IF EXISTS "appointments_appointment"
    DROP CONSTRAINT IF EXISTS "appointments_appointment_staff_id_fkey";

ALTER TABLE
    IF EXISTS "appointments_appointment"
    DROP CONSTRAINT IF EXISTS "appointments_appointment_facility_id_fkey";

ALTER TABLE
    IF EXISTS "appointments_appointment"
    DROP CONSTRAINT IF EXISTS "appointments_appointment_program_id_fkey";

ALTER TABLE
    IF EXISTS "authority_authorityrole_users"
    DROP CONSTRAINT IF EXISTS "authority_authorityrole_users_authorityrole_id_fkey";

ALTER TABLE
    IF EXISTS "authority_authorityrole_users"
    DROP CONSTRAINT IF EXISTS "authority_authorityrole_users_user_id_fkey";

ALTER TABLE
    IF EXISTS "common_notification"
    DROP CONSTRAINT IF EXISTS "common_notification_created_by_fkey";

ALTER TABLE
    IF EXISTS "common_notification"
    DROP CONSTRAINT IF EXISTS "common_notification_updated_by_fkey";

ALTER TABLE
    IF EXISTS "common_notification"
    DROP CONSTRAINT IF EXISTS "common_notification_facility_id_fkey";

ALTER TABLE
    IF EXISTS "common_notification"
    DROP CONSTRAINT IF EXISTS "common_notification_organisation_id_fkey";

ALTER TABLE
    IF EXISTS "common_notification"
    DROP CONSTRAINT IF EXISTS "common_notification_user_id_fkey";

ALTER TABLE
    IF EXISTS "common_notification"
    DROP CONSTRAINT IF EXISTS "common_notification_program_id_fkey";

ALTER TABLE
    IF EXISTS "common_usersurveys"
    DROP CONSTRAINT IF EXISTS "common_usersurveys_created_by_fkey";

ALTER TABLE
    IF EXISTS "common_usersurveys"
    DROP CONSTRAINT IF EXISTS "common_usersurveys_updated_by_fkey";

ALTER TABLE
    IF EXISTS "common_usersurveys"
    DROP CONSTRAINT IF EXISTS "common_usersurveys_organisation_id_fkey";

ALTER TABLE
    IF EXISTS "common_usersurveys"
    DROP CONSTRAINT IF EXISTS "common_usersurveys_user_id_fkey";

ALTER TABLE
    IF EXISTS "common_usersurveys"
    DROP CONSTRAINT IF EXISTS "common_usersurveys_program_id_fkey";

ALTER TABLE
    IF EXISTS "common_feedback"
    DROP CONSTRAINT IF EXISTS "common_feedback_created_by_fkey";

ALTER TABLE
    IF EXISTS "common_feedback"
    DROP CONSTRAINT IF EXISTS "common_feedback_updated_by_fkey";

ALTER TABLE
    IF EXISTS "common_feedback"
    DROP CONSTRAINT IF EXISTS "common_feedback_organisation_id_fkey";

ALTER TABLE
    IF EXISTS "common_feedback"
    DROP CONSTRAINT IF EXISTS "common_feedback_user_id_fkey";

ALTER TABLE
    IF EXISTS "common_feedback"
    DROP CONSTRAINT IF EXISTS "common_feedback_program_id_fkey";

ALTER TABLE
    IF EXISTS "communities_community_clients"
    DROP CONSTRAINT IF EXISTS "communities_community_clients_community_id_fkey";

ALTER TABLE
    IF EXISTS "communities_community_clients"
    DROP CONSTRAINT IF EXISTS "communities_community_clients_client_id_fkey";

ALTER TABLE
    IF EXISTS "communities_community_staff"
    DROP CONSTRAINT IF EXISTS "communities_community_staff_community_id_fkey";

ALTER TABLE
    IF EXISTS "communities_community_staff"
    DROP CONSTRAINT IF EXISTS "communities_community_staff_staff_id_fkey";

ALTER TABLE
    IF EXISTS "screeningtools_screeningtoolsresponse"
    DROP CONSTRAINT IF EXISTS "screeningtools_screeningtoolsresponse_created_by_fkey";

ALTER TABLE
    IF EXISTS "screeningtools_screeningtoolsresponse"
    DROP CONSTRAINT IF EXISTS "screeningtools_screeningtoolsresponse_updated_by_fkey";

ALTER TABLE
    IF EXISTS "screeningtools_screeningtoolsresponse"
    DROP CONSTRAINT IF EXISTS "screeningtools_screeningtoolsresponse_client_id_fkey";

ALTER TABLE
    IF EXISTS "screeningtools_screeningtoolsresponse"
    DROP CONSTRAINT IF EXISTS "screeningtools_screeningtoolsresponse_organisation_id_fkey";

ALTER TABLE
    IF EXISTS "screeningtools_screeningtoolsresponse"
    DROP CONSTRAINT IF EXISTS "screeningtools_screeningtoolsresponse_question_id_fkey";

ALTER TABLE
    IF EXISTS "screeningtools_screeningtoolsresponse"
    DROP CONSTRAINT IF EXISTS "screeningtools_screeningtoolsresponse_program_id_fkey";

ALTER TABLE
    IF EXISTS "staff_servicerequest"
    DROP CONSTRAINT IF EXISTS "staff_servicerequest_created_by_fkey";

ALTER TABLE
    IF EXISTS "staff_servicerequest"
    DROP CONSTRAINT IF EXISTS "staff_servicerequest_updated_by_fkey";

ALTER TABLE
    IF EXISTS "staff_servicerequest"
    DROP CONSTRAINT IF EXISTS "staff_servicerequest_organisation_id_fkey";

ALTER TABLE
    IF EXISTS "staff_servicerequest"
    DROP CONSTRAINT IF EXISTS "staff_servicerequest_resolved_by_id_fkey";

ALTER TABLE
    IF EXISTS "staff_servicerequest"
    DROP CONSTRAINT IF EXISTS "staff_servicerequest_staff_id_fkey";

ALTER TABLE
    IF EXISTS "staff_servicerequest"
    DROP CONSTRAINT IF EXISTS "staff_servicerequest_facility_id_fkey";

ALTER TABLE
    IF EXISTS "staff_servicerequest"
    DROP CONSTRAINT IF EXISTS "staff_servicerequest_program_id_fkey";

ALTER TABLE
    IF EXISTS "users_metric"
    DROP CONSTRAINT IF EXISTS "users_metric_created_by_fkey";

ALTER TABLE
    IF EXISTS "users_metric"
    DROP CONSTRAINT IF EXISTS "users_metric_updated_by_fkey";

ALTER TABLE
    IF EXISTS "users_metric"
    DROP CONSTRAINT IF EXISTS "users_metric_organisation_id_fkey";

ALTER TABLE
    IF EXISTS "users_userpin"
    DROP CONSTRAINT IF EXISTS "users_userpin_created_by_fkey";

ALTER TABLE
    IF EXISTS "users_userpin"
    DROP CONSTRAINT IF EXISTS "users_userpin_updated_by_fkey";

ALTER TABLE
    IF EXISTS "users_userpin"
    DROP CONSTRAINT IF EXISTS "users_userpin_user_id_fkey";

ALTER TABLE
    IF EXISTS "users_userpin"
    DROP CONSTRAINT IF EXISTS "users_userpin_organisation_id_fkey";

ALTER TABLE
    IF EXISTS "users_userotp"
    DROP CONSTRAINT IF EXISTS "users_userotp_user_id_fkey";

ALTER TABLE
    IF EXISTS "users_userotp"
    DROP CONSTRAINT IF EXISTS "users_userotp_created_by_fkey";

ALTER TABLE
    IF EXISTS "users_userotp"
    DROP CONSTRAINT IF EXISTS "users_userotp_updated_by_fkey";

ALTER TABLE
    IF EXISTS "questionnaires_questionnaire"
    DROP CONSTRAINT IF EXISTS "questionnaires_questionnaire_created_by_fkey";

ALTER TABLE
    IF EXISTS "questionnaires_questionnaire"
    DROP CONSTRAINT IF EXISTS "questionnaires_questionnaire_updated_by_fkey";

ALTER TABLE
    IF EXISTS "questionnaires_questionnaire"
    DROP CONSTRAINT IF EXISTS "questionnaires_questionnaire_organisation_id_fkey";

ALTER TABLE
    IF EXISTS "questionnaires_questionnaire"
    DROP CONSTRAINT IF EXISTS "questionnaires_questionnaire_program_id_fkey";

ALTER TABLE
    IF EXISTS "questionnaires_screeningtool"
    DROP CONSTRAINT IF EXISTS "questionnaires_screeningtool_created_by_fkey";

ALTER TABLE
    IF EXISTS "questionnaires_screeningtool"
    DROP CONSTRAINT IF EXISTS "questionnaires_screeningtool_updated_by_fkey";

ALTER TABLE
    IF EXISTS "questionnaires_screeningtool"
    DROP CONSTRAINT IF EXISTS "questionnaires_screeningtool_organisation_id_fkey";

ALTER TABLE
    IF EXISTS "questionnaires_screeningtool"
    DROP CONSTRAINT IF EXISTS "questionnaires_screeningtool_questionnaire_id_fkey";

ALTER TABLE
    IF EXISTS "questionnaires_screeningtool"
    DROP CONSTRAINT IF EXISTS "questionnaires_screeningtool_program_id_fkey";

ALTER TABLE
    IF EXISTS "questionnaires_question"
    DROP CONSTRAINT IF EXISTS "questionnaires_question_created_by_fkey";

ALTER TABLE
    IF EXISTS "questionnaires_question"
    DROP CONSTRAINT IF EXISTS "questionnaires_question_updated_by_fkey";

ALTER TABLE
    IF EXISTS "questionnaires_question"
    DROP CONSTRAINT IF EXISTS "questionnaires_question_organisation_id_fkey";

ALTER TABLE
    IF EXISTS "questionnaires_question"
    DROP CONSTRAINT IF EXISTS "questionnaires_question_questionnaire_id_fkey";

ALTER TABLE
    IF EXISTS "questionnaires_question"
    DROP CONSTRAINT IF EXISTS "questionnaires_question_program_id_fkey";

ALTER TABLE
    IF EXISTS "questionnaires_questioninputchoice"
    DROP CONSTRAINT IF EXISTS "questionnaires_questioninputchoice_created_by_fkey";

ALTER TABLE
    IF EXISTS "questionnaires_questioninputchoice"
    DROP CONSTRAINT IF EXISTS "questionnaires_questioninputchoice_updated_by_fkey";

ALTER TABLE
    IF EXISTS "questionnaires_questioninputchoice"
    DROP CONSTRAINT IF EXISTS "questionnaires_questioninputchoice_organisation_id_fkey";

ALTER TABLE
    IF EXISTS "questionnaires_questioninputchoice"
    DROP CONSTRAINT IF EXISTS "questionnaires_questioninputchoice_question_id_fkey";

ALTER TABLE
    IF EXISTS "questionnaires_questioninputchoice"
    DROP CONSTRAINT IF EXISTS "questionnaires_questioninputchoice_program_id_fkey";

ALTER TABLE
    IF EXISTS "questionnaires_screeningtoolresponse"
    DROP CONSTRAINT IF EXISTS "questionnaires_screeningtoolresponse_created_by_fkey";

ALTER TABLE
    IF EXISTS "questionnaires_screeningtoolresponse"
    DROP CONSTRAINT IF EXISTS "questionnaires_screeningtoolresponse_updated_by_fkey";

ALTER TABLE
    IF EXISTS "questionnaires_screeningtoolresponse"
    DROP CONSTRAINT IF EXISTS "questionnaires_screeningtoolresponse_organisation_id_fkey";

ALTER TABLE
    IF EXISTS "questionnaires_screeningtoolresponse"
    DROP CONSTRAINT IF EXISTS "questionnaires_screeningtoolresponse_screeningtool_id_fkey";

ALTER TABLE
    IF EXISTS "questionnaires_screeningtoolresponse"
    DROP CONSTRAINT IF EXISTS "questionnaires_screeningtoolresponse_facility_id_fkey";

ALTER TABLE
    IF EXISTS "questionnaires_screeningtoolresponse"
    DROP CONSTRAINT IF EXISTS "questionnaires_screeningtoolresponse_client_id_fkey";

ALTER TABLE
    IF EXISTS "questionnaires_screeningtoolresponse"
    DROP CONSTRAINT IF EXISTS "questionnaires_screeningtoolresponse_program_id_fkey";

ALTER TABLE
    IF EXISTS "questionnaires_screeningtoolquestionresponse"
    DROP CONSTRAINT IF EXISTS "questionnaires_screeningtoolquestionresponse_created_by_fkey";

ALTER TABLE
    IF EXISTS "questionnaires_screeningtoolquestionresponse"
    DROP CONSTRAINT IF EXISTS "questionnaires_screeningtoolquestionresponse_updated_by_fkey";

ALTER TABLE
    IF EXISTS "questionnaires_screeningtoolquestionresponse"
    DROP CONSTRAINT IF EXISTS "questionnes_screeningtoolquestionresponse_organisation_id_fkey";

ALTER TABLE
    IF EXISTS "questionnaires_screeningtoolquestionresponse"
    DROP CONSTRAINT IF EXISTS "questionnaires_screeniuestionresponse_screeniesponse_id_fkey";

ALTER TABLE
    IF EXISTS "questionnaires_screeningtoolquestionresponse"
    DROP CONSTRAINT IF EXISTS "questionnaires_screeningtoolquestionresponse_question_id_fkey";

ALTER TABLE
    IF EXISTS "questionnaires_screeningtoolquestionresponse"
    DROP CONSTRAINT IF EXISTS "questionnaires_screeningtoolquestionresponse_program_id_fkey";

ALTER TABLE
    IF EXISTS "questionnaires_screeningtoolquestionresponse"
    DROP CONSTRAINT IF EXISTS "questionnaires_screeningtoolquestionresponse_facility_id_fkey";

ALTER TABLE
    IF EXISTS "clients_client_facilities"
    DROP CONSTRAINT IF EXISTS "clients_client_facilities_client_id_fkey";

ALTER TABLE
    IF EXISTS "clients_client_facilities"
    DROP CONSTRAINT IF EXISTS "clients_client_facilities_facility_id_fkey";

ALTER TABLE
    IF EXISTS "caregivers_caregiver"
    DROP CONSTRAINT IF EXISTS "caregivers_caregiver_created_by_fkey";

ALTER TABLE
    IF EXISTS "caregivers_caregiver"
    DROP CONSTRAINT IF EXISTS "caregivers_caregiver_updated_by_fkey";

ALTER TABLE
    IF EXISTS "caregivers_caregiver"
    DROP CONSTRAINT IF EXISTS "caregivers_caregiver_organisation_id_fkey";

ALTER TABLE
    IF EXISTS "caregivers_caregiver"
    DROP CONSTRAINT IF EXISTS "caregivers_caregiver_user_id_fkey";

ALTER TABLE
    IF EXISTS "caregivers_caregiver"
    DROP CONSTRAINT IF EXISTS "caregivers_caregiver_program_id_fkey";

ALTER TABLE
    IF EXISTS "caregivers_caregiver_client"
    DROP CONSTRAINT IF EXISTS "caregivers_caregiver_client_caregiver_id_fkey";

ALTER TABLE
    IF EXISTS "caregivers_caregiver_client"
    DROP CONSTRAINT IF EXISTS "caregivers_caregiver_client_client_id_fkey";

ALTER TABLE
    IF EXISTS "caregivers_caregiver_client"
    DROP CONSTRAINT IF EXISTS "caregivers_caregiver_client_created_by_fkey";

ALTER TABLE
    IF EXISTS "caregivers_caregiver_client"
    DROP CONSTRAINT IF EXISTS "caregivers_caregiver_client_updated_by_fkey";

ALTER TABLE
    IF EXISTS "caregivers_caregiver_client"
    DROP CONSTRAINT IF EXISTS "caregivers_caregiver_client_organisation_id_fkey";

ALTER TABLE
    IF EXISTS "caregivers_caregiver_client"
    DROP CONSTRAINT IF EXISTS "caregivers_caregiver_client_assigned_by_fkey";

ALTER TABLE
    IF EXISTS "caregivers_caregiver_client"
    DROP CONSTRAINT IF EXISTS "caregivers_caregiver_client_program_id_fkey";

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

DROP TABLE IF EXISTS "clients_relatedperson_addresses";

DROP TABLE IF EXISTS "clients_relatedperson";

DROP TABLE IF EXISTS "common_securityquestion";

DROP TABLE IF EXISTS "common_address";

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