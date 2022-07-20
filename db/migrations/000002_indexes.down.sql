BEGIN;

-- common_facility
DROP INDEX IF EXISTS common_facility_name_21b959d6_like;

DROP INDEX IF EXISTS common_facility_organisation_id_32a7aac6;

DROP INDEX IF EXISTS common_facility_fhir_organization_id_8a5ed448_like;

-- common_facilityattachment
DROP INDEX IF EXISTS common_facilityattachment_facility_id_c043ede7;

DROP INDEX IF EXISTS common_facilityattachment_organisation_id_e69b418f;

-- common_organisation
DROP INDEX IF EXISTS common_organisation_org_code_ab2d714c_like;

DROP INDEX IF EXISTS common_organisation_organisation_name_c466bc44_like;

-- communities
DROP INDEX IF EXISTS communities_community_staff_community_id_c63e8e97;

DROP INDEX IF EXISTS communities_community_staff_staff_id_6ac91363;

-- userotp
DROP INDEX IF EXISTS users_userotp_user_id_7b41e28a;

-- userpin
DROP INDEX IF EXISTS users_userpin_user_id_aba00bd1;

DROP INDEX IF EXISTS users_userpin_user_id_user_type_c3e27cf8_idx;

-- security questions
DROP INDEX IF EXISTS clients_securityquestionresponse_organisation_id_8c3f6d29;

DROP INDEX IF EXISTS clients_securityquestionresponse_question_id_1efd91bd;

DROP INDEX IF EXISTS clients_securityquestionresponse_user_id_ff6cd092;

-- health diary
DROP INDEX IF EXISTS clients_healthdiaryentry_client_id_8c5de1dc;

DROP INDEX IF EXISTS clients_healthdiaryentry_organisation_id_78485008;

-- health diary attachment
DROP INDEX IF EXISTS clients_healthdiaryattachment_health_diary_entry_id_3b195bc0;

DROP INDEX IF EXISTS clients_healthdiaryattachment_organisation_id_128c56f5;

-- client service request
DROP INDEX IF EXISTS clients_servicerequest_client_id_640ced03;

DROP INDEX IF EXISTS clients_servicerequest_organisation_id_8a95b043;

DROP INDEX IF EXISTS clients_servicerequest_in_progress_by_id_c21e5f44;

DROP INDEX IF EXISTS clients_servicerequest_resolved_by_id_492fdef8;

DROP INDEX IF EXISTS clients_servicerequest_facility_id_15a97356;

-- staff addresses
DROP INDEX IF EXISTS staff_staff_addresses_staff_id_1620f1f1;

DROP INDEX IF EXISTS staff_staff_addresses_address_id_bdc155a3;

-- staff contacts
DROP INDEX IF EXISTS staff_staff_contacts_staff_id_9e46919d;

DROP INDEX IF EXISTS staff_staff_contacts_contact_id_1119fe01;

-- staff identifiers
DROP INDEX IF EXISTS staff_staff_identifiers_staff_id_4cc212da;

DROP INDEX IF EXISTS staff_staff_identifiers_identifier_id_eee797f4;

-- appointments
DROP INDEX IF EXISTS appointments_appointment_client_id_c189cb9e;

DROP INDEX IF EXISTS appointments_appointment_organisation_id_44b89706;

DROP INDEX IF EXISTS appointments_appointment_staff_id_64350136;

DROP INDEX IF EXISTS appointments_appointment_facility_id_52ace4b7;

DROP INDEX IF EXISTS appointments_appointment_external_id_53cd13e6_like;

-- authority role
DROP INDEX IF EXISTS authority_authorityrole_users_authorityrole_id_5e240d69;

DROP INDEX IF EXISTS authority_authorityrole_users_user_id_d225811e;

-- notifications
DROP INDEX IF EXISTS common_notification_facility_id_e51d5197;

DROP INDEX IF EXISTS common_notification_organisation_id_6d681c83;

DROP INDEX IF EXISTS common_notification_user_id_b851b5a2;

-- user surveys
DROP INDEX IF EXISTS common_usersurveys_organisation_id_83c7c493;

DROP INDEX IF EXISTS common_usersurveys_user_id_9c94c3eb;

-- feedback
DROP INDEX IF EXISTS common_feedback_organisation_id_cbd4e0e9;

DROP INDEX IF EXISTS common_feedback_user_id_3367f3cd;

-- communities
DROP INDEX IF EXISTS communities_community_clients_community_id_a2523222;

DROP INDEX IF EXISTS communities_community_clients_client_id_62435084;

-- screening tools
DROP INDEX IF EXISTS screeningtools_screeningtoolsresponse_client_id_4daffd54;

DROP INDEX IF EXISTS screeningtools_screeningtoolsresponse_organisation_id_5a1fb2fd;

DROP INDEX IF EXISTS screeningtools_screeningtoolsresponse_question_id_a1dc424b;

-- survey question
DROP INDEX IF EXISTS clients_securityquestion_organisation_id_e3a2aa56;

-- client related person
DROP INDEX IF EXISTS clients_relatedperson_organisation_id_70af9b64;

-- client related person id
DROP INDEX IF EXISTS clients_relatedperson_addresses_relatedperson_id_49e2b56b;

DROP INDEX IF EXISTS clients_relatedperson_addresses_address_id_95fddcd9;

-- client identifiers
DROP INDEX IF EXISTS clients_identifier_organisation_id_e1febd57;

-- clients_healthdiaryquote
DROP INDEX IF EXISTS clients_healthdiaryquote_quote_78ebad75_like;

DROP INDEX IF EXISTS clients_healthdiaryquote_organisation_id_fddd99d9;

-- clients_caregiver
DROP INDEX IF EXISTS clients_caregiver_organisation_id_14fee135;

-- authority_authoritypermission
DROP INDEX IF EXISTS authority_authoritypermission_organisation_id_d0abc669;

-- common_auditlog
DROP INDEX IF EXISTS common_auditlog_organisation_id_3c0012aa;

-- common_address
DROP INDEX IF EXISTS common_address_organisation_id_3a545a79;

-- authority_authorityrole
DROP INDEX IF EXISTS authority_authorityrole_organisation_id_be9dbfc9;

DROP INDEX IF EXISTS authority_authorityrole_permissions_authorityrole_id_88adf9a9;

DROP INDEX IF EXISTS authority_authorityrole_pe_authoritypermission_id_8bc52019;

-- communities
DROP INDEX IF EXISTS communities_community_organisation_id_01650979;

-- screeningtools_screeningtoolsquestion
DROP INDEX IF EXISTS screeningtools_screeningtoolsquestion_organisation_id_61a69743;


-- users_user
DROP INDEX IF EXISTS users_user_username_06e46fe6_like;

DROP INDEX IF EXISTS users_user_organisation_id_383db2a8;

DROP INDEX IF EXISTS users_user_accepted_terms_of_service_id_17891a2c;

-- common_contact
DROP INDEX IF EXISTS common_contact_organisation_id_85618961;

DROP INDEX IF EXISTS common_contact_user_id_id_0aea1c8e;

-- staff_staff
DROP INDEX IF EXISTS staff_staff_organisation_id_e73d0fb1;

DROP INDEX IF EXISTS staff_staff_default_facility_id_ac90e91c;

DROP INDEX IF EXISTS staff_staff_facilities_staff_id_2fdf3b17;

DROP INDEX IF EXISTS staff_staff_facilities_facility_id_e6dbfdf8;

-- clients_client
DROP INDEX IF EXISTS clients_client_fhir_patient_id_fd3230c2_like;

DROP INDEX IF EXISTS clients_client_emr_health_record_id_c37d5bfa_like;

DROP INDEX IF EXISTS clients_client_chv_id_7599da97;

DROP INDEX IF EXISTS clients_client_current_facility_id_32d5e567;

DROP INDEX IF EXISTS clients_client_organisation_id_641b9f26;

DROP INDEX IF EXISTS clients_client_addresses_client_id_82c2383b;

DROP INDEX IF EXISTS clients_client_addresses_address_id_3cc849da;

DROP INDEX IF EXISTS clients_client_contacts_client_id_71aac335;

DROP INDEX IF EXISTS clients_client_contacts_contact_id_f2666e22;

-- clients_relatedperson
DROP INDEX IF EXISTS clients_relatedperson_contacts_relatedperson_id_6387a794;

DROP INDEX IF EXISTS clients_relatedperson_contacts_contact_id_bf7953d0;

-- clients_clientfacility
DROP INDEX IF EXISTS clients_clientfacility_client_id_b75fe1c7;

DROP INDEX IF EXISTS clients_clientfacility_facility_id_0eef5c43;

DROP INDEX IF EXISTS clients_clientfacility_organisation_id_4011ecb9;

-- clients_client
DROP INDEX IF EXISTS clients_client_identifiers_client_id_24873e4d;

DROP INDEX IF EXISTS clients_client_identifiers_identifier_id_03c1aedc;

DROP INDEX IF EXISTS clients_client_related_persons_client_id_d013dcfc;

DROP INDEX IF EXISTS clients_client_related_persons_relatedperson_id_4217c22b;

-- staff_servicerequest
DROP INDEX IF EXISTS staff_servicerequest_organisation_id_4b38eea3;

DROP INDEX IF EXISTS staff_servicerequest_resolved_by_id_b680d0a6;

DROP INDEX IF EXISTS staff_servicerequest_staff_id_d52e14cc;

DROP INDEX IF EXISTS staff_servicerequest_facility_id_3c840b76;

DROP INDEX IF EXISTS users_metric_user_id_e7321b28;

COMMIT;