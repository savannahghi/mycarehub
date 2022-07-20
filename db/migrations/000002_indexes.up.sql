BEGIN;

-- common_facility
CREATE INDEX IF NOT EXISTS common_facility_name_21b959d6_like ON common_facility (name text_pattern_ops);

CREATE INDEX IF NOT EXISTS common_facility_organisation_id_32a7aac6 ON common_facility (organisation_id);

CREATE INDEX IF NOT EXISTS common_facility_fhir_organization_id_8a5ed448_like ON common_facility (fhir_organization_id varchar_pattern_ops);

-- common_facilityattachment
CREATE INDEX IF NOT EXISTS common_facilityattachment_facility_id_c043ede7 ON common_facilityattachment (facility_id);

CREATE INDEX IF NOT EXISTS common_facilityattachment_organisation_id_e69b418f ON common_facilityattachment (organisation_id);

-- common_organisation
CREATE INDEX IF NOT EXISTS common_organisation_org_code_ab2d714c_like ON common_organisation (org_code varchar_pattern_ops);

CREATE INDEX IF NOT EXISTS common_organisation_organisation_name_c466bc44_like ON common_organisation (organisation_name varchar_pattern_ops);

-- communities
CREATE INDEX IF NOT EXISTS communities_community_staff_community_id_c63e8e97 ON communities_community_staff (community_id);

CREATE INDEX IF NOT EXISTS communities_community_staff_staff_id_6ac91363 ON communities_community_staff (staff_id);

-- userotp
CREATE INDEX IF NOT EXISTS users_userotp_user_id_7b41e28a ON users_userotp (user_id);

-- userpin
CREATE INDEX IF NOT EXISTS users_userpin_user_id_aba00bd1 ON users_userpin (user_id);

CREATE INDEX IF NOT EXISTS users_userpin_user_id_user_type_c3e27cf8_idx ON users_userpin (user_id, user_type);

-- security questions
CREATE INDEX IF NOT EXISTS clients_securityquestionresponse_organisation_id_8c3f6d29 ON clients_securityquestionresponse (organisation_id);

CREATE INDEX IF NOT EXISTS clients_securityquestionresponse_question_id_1efd91bd ON clients_securityquestionresponse (question_id);

CREATE INDEX IF NOT EXISTS clients_securityquestionresponse_user_id_ff6cd092 ON clients_securityquestionresponse (user_id);

-- health diary
CREATE INDEX IF NOT EXISTS clients_healthdiaryentry_client_id_8c5de1dc ON clients_healthdiaryentry (client_id);

CREATE INDEX IF NOT EXISTS clients_healthdiaryentry_organisation_id_78485008 ON clients_healthdiaryentry (organisation_id);

-- health diary attachment
CREATE INDEX IF NOT EXISTS clients_healthdiaryattachment_health_diary_entry_id_3b195bc0 ON clients_healthdiaryattachment (health_diary_entry_id);

CREATE INDEX IF NOT EXISTS clients_healthdiaryattachment_organisation_id_128c56f5 ON clients_healthdiaryattachment (organisation_id);

-- client service request
CREATE INDEX IF NOT EXISTS clients_servicerequest_client_id_640ced03 ON clients_servicerequest (client_id);

CREATE INDEX IF NOT EXISTS clients_servicerequest_organisation_id_8a95b043 ON clients_servicerequest (organisation_id);

CREATE INDEX IF NOT EXISTS clients_servicerequest_in_progress_by_id_c21e5f44 ON clients_servicerequest (in_progress_by_id);

CREATE INDEX IF NOT EXISTS clients_servicerequest_resolved_by_id_492fdef8 ON clients_servicerequest (resolved_by_id);

CREATE INDEX IF NOT EXISTS clients_servicerequest_facility_id_15a97356 ON clients_servicerequest (facility_id);

-- staff addresses
CREATE INDEX IF NOT EXISTS staff_staff_addresses_staff_id_1620f1f1 ON staff_staff_addresses (staff_id);

CREATE INDEX IF NOT EXISTS staff_staff_addresses_address_id_bdc155a3 ON staff_staff_addresses (address_id);

-- staff contacts
CREATE INDEX IF NOT EXISTS staff_staff_contacts_staff_id_9e46919d ON staff_staff_contacts (staff_id);

CREATE INDEX IF NOT EXISTS staff_staff_contacts_contact_id_1119fe01 ON staff_staff_contacts (contact_id);

-- staff identifiers
CREATE INDEX IF NOT EXISTS staff_staff_identifiers_staff_id_4cc212da ON staff_staff_identifiers (staff_id);

CREATE INDEX IF NOT EXISTS staff_staff_identifiers_identifier_id_eee797f4 ON staff_staff_identifiers (identifier_id);

-- appointments
CREATE INDEX IF NOT EXISTS appointments_appointment_client_id_c189cb9e ON appointments_appointment (client_id);

CREATE INDEX IF NOT EXISTS appointments_appointment_organisation_id_44b89706 ON appointments_appointment (organisation_id);

CREATE INDEX IF NOT EXISTS appointments_appointment_staff_id_64350136 ON appointments_appointment (staff_id);

CREATE INDEX IF NOT EXISTS appointments_appointment_facility_id_52ace4b7 ON appointments_appointment (facility_id);

CREATE INDEX IF NOT EXISTS appointments_appointment_external_id_53cd13e6_like ON appointments_appointment (external_id varchar_pattern_ops);

-- authority role
CREATE INDEX IF NOT EXISTS authority_authorityrole_users_authorityrole_id_5e240d69 ON authority_authorityrole_users (authorityrole_id);

CREATE INDEX IF NOT EXISTS authority_authorityrole_users_user_id_d225811e ON authority_authorityrole_users (user_id);

-- notifications
CREATE INDEX IF NOT EXISTS common_notification_facility_id_e51d5197 ON common_notification (facility_id);

CREATE INDEX IF NOT EXISTS common_notification_organisation_id_6d681c83 ON common_notification (organisation_id);

CREATE INDEX IF NOT EXISTS common_notification_user_id_b851b5a2 ON common_notification (user_id);

-- user surveys
CREATE INDEX IF NOT EXISTS common_usersurveys_organisation_id_83c7c493 ON common_usersurveys (organisation_id);

CREATE INDEX IF NOT EXISTS common_usersurveys_user_id_9c94c3eb ON common_usersurveys (user_id);

-- feedback
CREATE INDEX IF NOT EXISTS common_feedback_organisation_id_cbd4e0e9 ON common_feedback (organisation_id);

CREATE INDEX IF NOT EXISTS common_feedback_user_id_3367f3cd ON common_feedback (user_id);

-- communities
CREATE INDEX IF NOT EXISTS communities_community_clients_community_id_a2523222 ON communities_community_clients (community_id);

CREATE INDEX IF NOT EXISTS communities_community_clients_client_id_62435084 ON communities_community_clients (client_id);

-- screening tools
CREATE INDEX IF NOT EXISTS screeningtools_screeningtoolsresponse_client_id_4daffd54 ON screeningtools_screeningtoolsresponse (client_id);

CREATE INDEX IF NOT EXISTS screeningtools_screeningtoolsresponse_organisation_id_5a1fb2fd ON screeningtools_screeningtoolsresponse (organisation_id);

CREATE INDEX IF NOT EXISTS screeningtools_screeningtoolsresponse_question_id_a1dc424b ON screeningtools_screeningtoolsresponse (question_id);

-- survey question
CREATE INDEX IF NOT EXISTS clients_securityquestion_organisation_id_e3a2aa56 ON clients_securityquestion (organisation_id);

-- client related person
CREATE INDEX IF NOT EXISTS clients_relatedperson_organisation_id_70af9b64 ON clients_relatedperson (organisation_id);

-- client related person id
CREATE INDEX IF NOT EXISTS clients_relatedperson_addresses_relatedperson_id_49e2b56b ON clients_relatedperson_addresses (relatedperson_id);

CREATE INDEX IF NOT EXISTS clients_relatedperson_addresses_address_id_95fddcd9 ON clients_relatedperson_addresses (address_id);

-- client identifiers
CREATE INDEX IF NOT EXISTS clients_identifier_organisation_id_e1febd57 ON clients_identifier (organisation_id);

-- clients_healthdiaryquote
CREATE INDEX IF NOT EXISTS clients_healthdiaryquote_quote_78ebad75_like ON clients_healthdiaryquote (quote text_pattern_ops);

CREATE INDEX IF NOT EXISTS clients_healthdiaryquote_organisation_id_fddd99d9 ON clients_healthdiaryquote (organisation_id);

-- clients_caregiver
CREATE INDEX IF NOT EXISTS clients_caregiver_organisation_id_14fee135 ON clients_caregiver (organisation_id);

-- authority_authoritypermission
CREATE INDEX IF NOT EXISTS authority_authoritypermission_organisation_id_d0abc669 ON authority_authoritypermission (organisation_id);

-- common_auditlog
CREATE INDEX IF NOT EXISTS common_auditlog_organisation_id_3c0012aa ON common_auditlog (organisation_id);

-- common_address
CREATE INDEX IF NOT EXISTS common_address_organisation_id_3a545a79 ON common_address (organisation_id);

-- authority_authorityrole
CREATE INDEX IF NOT EXISTS authority_authorityrole_organisation_id_be9dbfc9 ON authority_authorityrole (organisation_id);

CREATE INDEX IF NOT EXISTS authority_authorityrole_permissions_authorityrole_id_88adf9a9 ON authority_authorityrole_permissions (authorityrole_id);

CREATE INDEX IF NOT EXISTS authority_authorityrole_pe_authoritypermission_id_8bc52019 ON authority_authorityrole_permissions (authoritypermission_id);

-- communities
CREATE INDEX IF NOT EXISTS communities_community_organisation_id_01650979 ON communities_community (organisation_id);

-- screeningtools_screeningtoolsquestion
CREATE INDEX IF NOT EXISTS screeningtools_screeningtoolsquestion_organisation_id_61a69743 ON screeningtools_screeningtoolsquestion (organisation_id);


-- users_user
CREATE INDEX IF NOT EXISTS users_user_username_06e46fe6_like ON users_user (username varchar_pattern_ops);

CREATE INDEX IF NOT EXISTS users_user_organisation_id_383db2a8 ON users_user (organisation_id);

CREATE INDEX IF NOT EXISTS users_user_accepted_terms_of_service_id_17891a2c ON users_user (accepted_terms_of_service_id);

-- common_contact
CREATE INDEX IF NOT EXISTS common_contact_organisation_id_85618961 ON common_contact (organisation_id);

CREATE INDEX IF NOT EXISTS common_contact_user_id_id_0aea1c8e ON common_contact (user_id);

-- staff_staff
CREATE INDEX IF NOT EXISTS staff_staff_organisation_id_e73d0fb1 ON staff_staff (organisation_id);

CREATE INDEX IF NOT EXISTS staff_staff_default_facility_id_ac90e91c ON staff_staff (default_facility_id);

CREATE INDEX IF NOT EXISTS staff_staff_facilities_staff_id_2fdf3b17 ON staff_staff_facilities (staff_id);

CREATE INDEX IF NOT EXISTS staff_staff_facilities_facility_id_e6dbfdf8 ON staff_staff_facilities (facility_id);

-- clients_client
CREATE INDEX IF NOT EXISTS clients_client_fhir_patient_id_fd3230c2_like ON clients_client (fhir_patient_id text_pattern_ops);

CREATE INDEX IF NOT EXISTS clients_client_emr_health_record_id_c37d5bfa_like ON clients_client (emr_health_record_id text_pattern_ops);

CREATE INDEX IF NOT EXISTS clients_client_chv_id_7599da97 ON clients_client (chv_id);

CREATE INDEX IF NOT EXISTS clients_client_current_facility_id_32d5e567 ON clients_client (current_facility_id);

CREATE INDEX IF NOT EXISTS clients_client_organisation_id_641b9f26 ON clients_client (organisation_id);

CREATE INDEX IF NOT EXISTS clients_client_addresses_client_id_82c2383b ON clients_client_addresses (client_id);

CREATE INDEX IF NOT EXISTS clients_client_addresses_address_id_3cc849da ON clients_client_addresses (address_id);

CREATE INDEX IF NOT EXISTS clients_client_contacts_client_id_71aac335 ON clients_client_contacts (client_id);

CREATE INDEX IF NOT EXISTS clients_client_contacts_contact_id_f2666e22 ON clients_client_contacts (contact_id);

-- clients_relatedperson
CREATE INDEX IF NOT EXISTS clients_relatedperson_contacts_relatedperson_id_6387a794 ON clients_relatedperson_contacts (relatedperson_id);

CREATE INDEX IF NOT EXISTS clients_relatedperson_contacts_contact_id_bf7953d0 ON clients_relatedperson_contacts (contact_id);

-- clients_clientfacility
CREATE INDEX IF NOT EXISTS clients_clientfacility_client_id_b75fe1c7 ON clients_clientfacility (client_id);

CREATE INDEX IF NOT EXISTS clients_clientfacility_facility_id_0eef5c43 ON clients_clientfacility (facility_id);

CREATE INDEX IF NOT EXISTS clients_clientfacility_organisation_id_4011ecb9 ON clients_clientfacility (organisation_id);

-- clients_client
CREATE INDEX IF NOT EXISTS clients_client_identifiers_client_id_24873e4d ON clients_client_identifiers (client_id);

CREATE INDEX IF NOT EXISTS clients_client_identifiers_identifier_id_03c1aedc ON clients_client_identifiers (identifier_id);

CREATE INDEX IF NOT EXISTS clients_client_related_persons_client_id_d013dcfc ON clients_client_related_persons (client_id);

CREATE INDEX IF NOT EXISTS clients_client_related_persons_relatedperson_id_4217c22b ON clients_client_related_persons (relatedperson_id);

-- staff_servicerequest
CREATE INDEX IF NOT EXISTS staff_servicerequest_organisation_id_4b38eea3 ON staff_servicerequest (organisation_id);

CREATE INDEX IF NOT EXISTS staff_servicerequest_resolved_by_id_b680d0a6 ON staff_servicerequest (resolved_by_id);

CREATE INDEX IF NOT EXISTS staff_servicerequest_staff_id_d52e14cc ON staff_servicerequest (staff_id);

CREATE INDEX IF NOT EXISTS staff_servicerequest_facility_id_3c840b76 ON staff_servicerequest (facility_id);

CREATE INDEX IF NOT EXISTS users_metric_user_id_e7321b28 ON users_metric (user_id);

COMMIT;