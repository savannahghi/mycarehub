BEGIN;

CREATE TABLE IF NOT EXISTS common_organisation (
    id uuid NOT NULL PRIMARY KEY,
    created timestamp WITH time zone NOT NULL,
    updated timestamp WITH time zone NOT NULL,
    created_by uuid,
    updated_by uuid,
    code integer NOT NULL UNIQUE,
    active boolean NOT NULL,
    org_code varchar(15) NOT NULL UNIQUE,
    organisation_name varchar(100) NOT NULL UNIQUE,
    email_address varchar(100) NOT NULL,
    phone_number varchar(128) NOT NULL,
    description text,
    postal_address varchar(100) NOT NULL,
    physical_address text NOT NULL,
    default_country varchar(255) NOT NULL,
    deleted_at timestamp WITH time zone
);

CREATE TABLE IF NOT EXISTS common_facility (
    id uuid NOT NULL PRIMARY KEY,
    active boolean NOT NULL,
    created timestamp WITH time zone NOT NULL,
    created_by uuid,
    updated timestamp WITH time zone NOT NULL,
    updated_by uuid,
    name text NOT NULL UNIQUE,
    mfl_code integer NOT NULL UNIQUE,
    county varchar(64) NOT NULL,
    organisation_id uuid NOT NULL CONSTRAINT common_facility_organisation_id_32a7aac6_fk_common_or REFERENCES common_organisation deferrable initially deferred,
    description text NOT NULL,
    deleted_at timestamp WITH time zone,
    phone varchar(15),
    fhir_organization_id varchar(64) UNIQUE
);

CREATE TABLE IF NOT EXISTS common_facilityattachment (
    id uuid NOT NULL PRIMARY KEY,
    active boolean NOT NULL,
    created timestamp WITH time zone NOT NULL,
    created_by uuid,
    updated timestamp WITH time zone NOT NULL,
    updated_by uuid,
    content_type varchar(100) NOT NULL,
    data varchar(65535) NOT NULL,
    title varchar(255) NOT NULL,
    creation_date timestamp WITH time zone NOT NULL,
    size integer,
    description text,
    aspect_ratio varchar(50),
    notes text NOT NULL,
    facility_id uuid NOT NULL CONSTRAINT common_facilityattac_facility_id_c043ede7_fk_common_fa REFERENCES common_facility deferrable initially deferred,
    organisation_id uuid NOT NULL CONSTRAINT common_facilityattac_organisation_id_e69b418f_fk_common_or REFERENCES common_organisation deferrable initially deferred,
    deleted_at timestamp WITH time zone
);

CREATE TABLE IF NOT EXISTS common_auditlog (
    id uuid NOT NULL PRIMARY KEY,
    active boolean NOT NULL,
    created timestamp WITH time zone NOT NULL,
    created_by uuid,
    updated timestamp WITH time zone NOT NULL,
    updated_by uuid,
    timestamp timestamp WITH time zone NOT NULL,
    record_type text NOT NULL,
    notes text NOT NULL,
    payload jsonb NOT NULL,
    organisation_id uuid NOT NULL CONSTRAINT common_auditlog_organisation_id_3c0012aa_fk_common_or REFERENCES common_organisation deferrable initially deferred,
    deleted_at timestamp WITH time zone
);

CREATE TABLE IF NOT EXISTS common_address (
    id uuid NOT NULL PRIMARY KEY,
    active boolean NOT NULL,
    created timestamp WITH time zone NOT NULL,
    created_by uuid,
    updated timestamp WITH time zone NOT NULL,
    updated_by uuid,
    address_type varchar(16) NOT NULL,
    text text NOT NULL,
    postal_code text NOT NULL,
    country varchar(255) NOT NULL,
    organisation_id uuid NOT NULL CONSTRAINT common_address_organisation_id_3a545a79_fk_common_or REFERENCES common_organisation deferrable initially deferred,
    deleted_at timestamp WITH time zone
);

CREATE TABLE IF NOT EXISTS clients_securityquestion (
    id uuid NOT NULL PRIMARY KEY,
    active boolean NOT NULL,
    created timestamp WITH time zone NOT NULL,
    created_by uuid,
    updated timestamp WITH time zone NOT NULL,
    updated_by uuid,
    stem text NOT NULL,
    description text NOT NULL,
    sequence integer NOT NULL,
    response_type varchar(32) NOT NULL,
    organisation_id uuid NOT NULL CONSTRAINT clients_securityques_organisation_id_e3a2aa56_fk_common_or REFERENCES common_organisation deferrable initially deferred,
    deleted_at timestamp WITH time zone,
    flavour varchar(32)
);

CREATE TABLE IF NOT EXISTS clients_relatedperson (
    id uuid NOT NULL PRIMARY KEY,
    active boolean NOT NULL,
    created timestamp WITH time zone NOT NULL,
    created_by uuid,
    updated timestamp WITH time zone NOT NULL,
    updated_by uuid,
    first_name text NOT NULL,
    last_name text NOT NULL,
    other_name text NOT NULL,
    date_of_birth date,
    gender varchar(16) NOT NULL,
    relationship_type varchar(64) NOT NULL,
    organisation_id uuid NOT NULL CONSTRAINT clients_relatedperso_organisation_id_70af9b64_fk_common_or REFERENCES common_organisation deferrable initially deferred,
    deleted_at timestamp WITH time zone
);

CREATE TABLE IF NOT EXISTS clients_relatedperson_addresses (
    id serial PRIMARY KEY,
    relatedperson_id uuid NOT NULL CONSTRAINT clients_relatedperso_relatedperson_id_49e2b56b_fk_clients_r REFERENCES clients_relatedperson deferrable initially deferred,
    address_id uuid NOT NULL CONSTRAINT clients_relatedperso_address_id_95fddcd9_fk_common_ad REFERENCES common_address deferrable initially deferred,
    CONSTRAINT clients_relatedperson_ad_relatedperson_id_address_16b5f53a_uniq UNIQUE (relatedperson_id, address_id)
);

CREATE TABLE IF NOT EXISTS clients_identifier (
    id uuid NOT NULL PRIMARY KEY,
    active boolean NOT NULL,
    created timestamp WITH time zone NOT NULL,
    created_by uuid,
    updated timestamp WITH time zone NOT NULL,
    updated_by uuid,
    identifier_type varchar(64) NOT NULL,
    identifier_value text NOT NULL,
    identifier_use varchar(64) NOT NULL,
    description text NOT NULL,
    valid_from timestamp WITH time zone NOT NULL,
    valid_to timestamp WITH time zone,
    is_primary_identifier boolean NOT NULL,
    organisation_id uuid NOT NULL CONSTRAINT clients_identifier_organisation_id_e1febd57_fk_common_or REFERENCES common_organisation deferrable initially deferred,
    deleted_at timestamp WITH time zone,
    CONSTRAINT clients_identifier_identifier_type_identifi_45671f8e_uniq UNIQUE (identifier_type, identifier_value)
);

CREATE TABLE IF NOT EXISTS clients_healthdiaryquote (
    id uuid NOT NULL PRIMARY KEY,
    active boolean NOT NULL,
    created timestamp WITH time zone NOT NULL,
    created_by uuid,
    updated timestamp WITH time zone NOT NULL,
    updated_by uuid,
    deleted_at timestamp WITH time zone,
    quote text NOT NULL UNIQUE,
    organisation_id uuid NOT NULL CONSTRAINT clients_healthdiaryq_organisation_id_fddd99d9_fk_common_or REFERENCES common_organisation deferrable initially deferred,
    by text NOT NULL
);

CREATE TABLE IF NOT EXISTS clients_caregiver (
    id uuid NOT NULL PRIMARY KEY,
    active boolean NOT NULL,
    created timestamp WITH time zone NOT NULL,
    created_by uuid,
    updated timestamp WITH time zone NOT NULL,
    updated_by uuid,
    deleted_at timestamp WITH time zone,
    first_name text NOT NULL,
    last_name text NOT NULL,
    caregiver_type varchar(64) NOT NULL,
    phone_number text,
    organisation_id uuid NOT NULL CONSTRAINT clients_caregiver_organisation_id_14fee135_fk_common_or REFERENCES common_organisation deferrable initially deferred
);

CREATE TABLE IF NOT EXISTS authority_authoritypermission (
    id uuid NOT NULL PRIMARY KEY,
    active boolean NOT NULL,
    created timestamp WITH time zone NOT NULL,
    created_by uuid,
    updated timestamp WITH time zone NOT NULL,
    updated_by uuid,
    deleted_at timestamp WITH time zone,
    name varchar(255) NOT NULL,
    organisation_id uuid NOT NULL CONSTRAINT authority_authorityp_organisation_id_d0abc669_fk_common_or REFERENCES common_organisation deferrable initially deferred
);

CREATE TABLE IF NOT EXISTS authority_authorityrole (
    id uuid NOT NULL PRIMARY KEY,
    active boolean NOT NULL,
    created timestamp WITH time zone NOT NULL,
    created_by uuid,
    updated timestamp WITH time zone NOT NULL,
    updated_by uuid,
    deleted_at timestamp WITH time zone,
    name varchar(255) NOT NULL,
    organisation_id uuid NOT NULL CONSTRAINT authority_authorityr_organisation_id_be9dbfc9_fk_common_or REFERENCES common_organisation deferrable initially deferred
);

CREATE TABLE IF NOT EXISTS authority_authorityrole_permissions (
    id serial PRIMARY KEY,
    authorityrole_id uuid NOT NULL CONSTRAINT authority_authorityr_authorityrole_id_88adf9a9_fk_authority REFERENCES authority_authorityrole deferrable initially deferred,
    authoritypermission_id uuid NOT NULL CONSTRAINT authority_authorityr_authoritypermission__8bc52019_fk_authority REFERENCES authority_authoritypermission deferrable initially deferred,
    CONSTRAINT authority_authorityrole__authorityrole_id_authori_6d452486_uniq UNIQUE (authorityrole_id, authoritypermission_id)
);

CREATE TABLE IF NOT EXISTS communities_community (
    id uuid NOT NULL PRIMARY KEY,
    active boolean NOT NULL,
    created timestamp WITH time zone NOT NULL,
    created_by uuid,
    updated timestamp WITH time zone NOT NULL,
    updated_by uuid,
    deleted_at timestamp WITH time zone,
    name varchar(64) NOT NULL,
    description varchar(150) NOT NULL,
    client_types varchar(64) [] NOT NULL,
    gender varchar(16) [],
    invite_only boolean NOT NULL,
    discoverable boolean NOT NULL,
    organisation_id uuid NOT NULL CONSTRAINT communities_communit_organisation_id_01650979_fk_common_or REFERENCES common_organisation deferrable initially deferred,
    max_age integer,
    min_age integer
);


CREATE TABLE IF NOT EXISTS screeningtools_screeningtoolsquestion (
    id uuid NOT NULL PRIMARY KEY,
    active boolean NOT NULL,
    created timestamp WITH time zone NOT NULL,
    created_by uuid,
    updated timestamp WITH time zone NOT NULL,
    updated_by uuid,
    deleted_at timestamp WITH time zone,
    question text NOT NULL,
    tool_type varchar(32) NOT NULL,
    response_choices jsonb,
    response_type varchar(32) NOT NULL,
    response_category varchar(32) NOT NULL,
    sequence integer NOT NULL,
    meta jsonb,
    organisation_id uuid NOT NULL CONSTRAINT screeningtools_scree_organisation_id_61a69743_fk_common_or REFERENCES common_organisation deferrable initially deferred
);

CREATE TABLE IF NOT EXISTS users_termsofservice (
    id serial PRIMARY KEY,
    text text NOT NULL,
    valid_from timestamp WITH time zone NOT NULL,
    valid_to timestamp WITH time zone,
    active boolean NOT NULL,
    created timestamp WITH time zone NOT NULL,
    created_by uuid,
    updated timestamp WITH time zone NOT NULL,
    updated_by uuid,
    deleted_at timestamp WITH time zone,
    flavour varchar(32)
);

CREATE TABLE IF NOT EXISTS users_user (
    PASSWORD varchar(128) NOT NULL,
    last_login timestamp WITH time zone,
    is_superuser boolean NOT NULL,
    username varchar(150) NOT NULL UNIQUE,
    email varchar(254) NOT NULL,
    is_staff boolean NOT NULL,
    is_active boolean NOT NULL,
    date_joined timestamp WITH time zone NOT NULL,
    id uuid NOT NULL PRIMARY KEY,
    name varchar(255) NOT NULL,
    is_approved boolean NOT NULL,
    approval_notified boolean NOT NULL,
    phone varchar(128),
    organisation_id uuid NOT NULL CONSTRAINT users_user_organisation_id_383db2a8_fk_common_organisation_id REFERENCES common_organisation deferrable initially deferred,
    date_of_birth date,
    failed_login_count integer NOT NULL,
    first_name varchar(150) NOT NULL,
    gender varchar(16),
    handle text NOT NULL,
    last_failed_login timestamp WITH time zone,
    last_name varchar(150) NOT NULL,
    last_successful_login timestamp WITH time zone,
    middle_name text NOT NULL,
    next_allowed_login timestamp WITH time zone NOT NULL,
    push_tokens varchar(256) [],
    user_type varchar(32),
    accepted_terms_of_service_id integer CONSTRAINT users_user_accepted_terms_of_se_17891a2c_fk_users_ter REFERENCES users_termsofservice deferrable initially deferred,
    flavour varchar(32),
    terms_accepted boolean NOT NULL,
    avatar text,
    is_suspended boolean NOT NULL,
    pin_change_required boolean,
    has_set_pin boolean,
    has_set_security_questions boolean,
    is_phone_verified boolean,
    failed_security_count integer NOT NULL,
    pin_update_required boolean,
    has_set_nickname boolean
);


CREATE TABLE IF NOT EXISTS common_contact (
    id uuid NOT NULL PRIMARY KEY,
    active boolean NOT NULL,
    created timestamp WITH time zone NOT NULL,
    created_by uuid,
    updated timestamp WITH time zone NOT NULL,
    updated_by uuid,
    contact_type varchar(16) NOT NULL,
    contact_value text NOT NULL,
    opted_in boolean NOT NULL,
    organisation_id uuid NOT NULL CONSTRAINT common_contact_organisation_id_85618961_fk_common_or REFERENCES common_organisation deferrable initially deferred,
    deleted_at timestamp WITH time zone,
    user_id uuid CONSTRAINT common_contact_user_id_1383bcb3_fk_users_user_id REFERENCES users_user deferrable initially deferred,
    flavour varchar(32),
    CONSTRAINT common_contact_contact_value_flavour_07f55910_uniq UNIQUE (contact_value, flavour)
);

CREATE TABLE IF NOT EXISTS staff_staff (
    id uuid NOT NULL PRIMARY KEY,
    active boolean NOT NULL,
    created timestamp WITH time zone NOT NULL,
    created_by uuid,
    updated timestamp WITH time zone NOT NULL,
    updated_by uuid,
    deleted_at timestamp WITH time zone,
    staff_number varchar(255),
    default_facility_id uuid NOT NULL CONSTRAINT staff_staff_default_facility_id_ac90e91c_fk_common_facility_id REFERENCES common_facility deferrable initially deferred,
    organisation_id uuid NOT NULL CONSTRAINT staff_staff_organisation_id_e73d0fb1_fk_common_organisation_id REFERENCES common_organisation deferrable initially deferred,
    user_id uuid UNIQUE CONSTRAINT staff_staff_user_id_abd8ba40_fk_users_user_id REFERENCES users_user deferrable initially deferred
);

CREATE TABLE IF NOT EXISTS staff_staff_facilities (
    id serial PRIMARY KEY,
    staff_id uuid NOT NULL CONSTRAINT staff_staff_facilities_staff_id_2fdf3b17_fk_staff_staff_id REFERENCES staff_staff deferrable initially deferred,
    facility_id uuid NOT NULL CONSTRAINT staff_staff_faciliti_facility_id_e6dbfdf8_fk_common_fa REFERENCES common_facility deferrable initially deferred,
    CONSTRAINT staff_staff_facilities_staff_id_facility_id_81a27bda_uniq UNIQUE (staff_id, facility_id)
);

CREATE TABLE IF NOT EXISTS clients_client (
    id uuid NOT NULL PRIMARY KEY,
    active boolean NOT NULL,
    created timestamp WITH time zone NOT NULL,
    created_by uuid,
    updated timestamp WITH time zone NOT NULL,
    updated_by uuid,
    enrollment_date timestamp WITH time zone NOT NULL,
    fhir_patient_id text UNIQUE,
    emr_health_record_id text UNIQUE,
    treatment_buddy text,
    counselled boolean NOT NULL,
    languages varchar(150) [],
    chv_id uuid CONSTRAINT clients_client_chv_id_7599da97_fk_users_user_id REFERENCES users_user deferrable initially deferred,
    current_facility_id uuid NOT NULL CONSTRAINT clients_client_current_facility_id_32d5e567_fk_common_fa REFERENCES common_facility deferrable initially deferred,
    organisation_id uuid NOT NULL CONSTRAINT clients_client_organisation_id_641b9f26_fk_common_or REFERENCES common_organisation deferrable initially deferred,
    user_id uuid UNIQUE CONSTRAINT clients_client_user_id_57c8fc4a_fk_users_user_id REFERENCES users_user deferrable initially deferred,
    deleted_at timestamp WITH time zone,
    caregiver_id uuid UNIQUE CONSTRAINT clients_client_caregiver_id_0033fe31_fk_clients_caregiver_id REFERENCES clients_caregiver deferrable initially deferred,
    client_types varchar(64) [] NOT NULL
);

CREATE TABLE IF NOT EXISTS clients_client_addresses (
    id serial PRIMARY KEY,
    client_id uuid NOT NULL CONSTRAINT clients_client_addre_client_id_82c2383b_fk_clients_c REFERENCES clients_client deferrable initially deferred,
    address_id uuid NOT NULL CONSTRAINT clients_client_addre_address_id_3cc849da_fk_common_ad REFERENCES common_address deferrable initially deferred,
    CONSTRAINT clients_client_addresses_client_id_address_id_c422f723_uniq UNIQUE (client_id, address_id)
);

CREATE TABLE IF NOT EXISTS clients_client_contacts (
    id serial PRIMARY KEY,
    client_id uuid NOT NULL CONSTRAINT clients_client_contacts_client_id_71aac335_fk_clients_client_id REFERENCES clients_client deferrable initially deferred,
    contact_id uuid NOT NULL CONSTRAINT clients_client_conta_contact_id_f2666e22_fk_common_co REFERENCES common_contact deferrable initially deferred,
    CONSTRAINT clients_client_contacts_client_id_contact_id_6bf99dcf_uniq UNIQUE (client_id, contact_id)
);

CREATE TABLE IF NOT EXISTS clients_relatedperson_contacts (
    id serial PRIMARY KEY,
    relatedperson_id uuid NOT NULL CONSTRAINT clients_relatedperso_relatedperson_id_6387a794_fk_clients_r REFERENCES clients_relatedperson deferrable initially deferred,
    contact_id uuid NOT NULL CONSTRAINT clients_relatedperso_contact_id_bf7953d0_fk_common_co REFERENCES common_contact deferrable initially deferred,
    CONSTRAINT clients_relatedperson_co_relatedperson_id_contact_0b11ce1d_uniq UNIQUE (relatedperson_id, contact_id)
);

CREATE TABLE IF NOT EXISTS clients_clientfacility (
    id uuid NOT NULL PRIMARY KEY,
    active boolean NOT NULL,
    created timestamp WITH time zone NOT NULL,
    created_by uuid,
    updated timestamp WITH time zone NOT NULL,
    updated_by uuid,
    assigned timestamp WITH time zone NOT NULL,
    transferred_out timestamp WITH time zone,
    client_id uuid NOT NULL CONSTRAINT clients_clientfacility_client_id_b75fe1c7_fk_clients_client_id REFERENCES clients_client deferrable initially deferred,
    facility_id uuid NOT NULL CONSTRAINT clients_clientfacili_facility_id_0eef5c43_fk_common_fa REFERENCES common_facility deferrable initially deferred,
    organisation_id uuid NOT NULL CONSTRAINT clients_clientfacili_organisation_id_4011ecb9_fk_common_or REFERENCES common_organisation deferrable initially deferred,
    notes text NOT NULL,
    deleted_at timestamp WITH time zone,
    CONSTRAINT clients_clientfacility_client_id_facility_id_b9a5c3bf_uniq UNIQUE (client_id, facility_id)
);

CREATE TABLE IF NOT EXISTS clients_client_identifiers (
    id serial PRIMARY KEY,
    client_id uuid NOT NULL CONSTRAINT clients_client_ident_client_id_24873e4d_fk_clients_c REFERENCES clients_client deferrable initially deferred,
    identifier_id uuid NOT NULL CONSTRAINT clients_client_ident_identifier_id_03c1aedc_fk_clients_i REFERENCES clients_identifier deferrable initially deferred,
    CONSTRAINT clients_client_identifie_client_id_identifier_id_94271e77_uniq UNIQUE (client_id, identifier_id)
);

CREATE TABLE IF NOT EXISTS clients_client_related_persons (
    id serial PRIMARY KEY,
    client_id uuid NOT NULL CONSTRAINT clients_client_relat_client_id_d013dcfc_fk_clients_c REFERENCES clients_client deferrable initially deferred,
    relatedperson_id uuid NOT NULL CONSTRAINT clients_client_relat_relatedperson_id_4217c22b_fk_clients_r REFERENCES clients_relatedperson deferrable initially deferred,
    CONSTRAINT clients_client_related_p_client_id_relatedperson__128d6bde_uniq UNIQUE (client_id, relatedperson_id)
);

CREATE TABLE IF NOT EXISTS clients_securityquestionresponse (
    id uuid NOT NULL PRIMARY KEY,
    active boolean NOT NULL,
    created timestamp WITH time zone NOT NULL,
    created_by uuid,
    updated timestamp WITH time zone NOT NULL,
    updated_by uuid,
    timestamp timestamp WITH time zone NOT NULL,
    response text NOT NULL,
    organisation_id uuid NOT NULL CONSTRAINT clients_securityques_organisation_id_8c3f6d29_fk_common_or REFERENCES common_organisation deferrable initially deferred,
    question_id uuid NOT NULL CONSTRAINT clients_securityques_question_id_1efd91bd_fk_clients_s REFERENCES clients_securityquestion deferrable initially deferred,
    user_id uuid NOT NULL CONSTRAINT clients_securityques_user_id_ff6cd092_fk_users_use REFERENCES users_user deferrable initially deferred,
    deleted_at timestamp WITH time zone,
    is_correct boolean NOT NULL,
    CONSTRAINT clients_securityquestion_user_id_question_id_e9ed4932_uniq UNIQUE (user_id, question_id)
);

CREATE TABLE IF NOT EXISTS clients_healthdiaryentry (
    id uuid NOT NULL PRIMARY KEY,
    active boolean NOT NULL,
    created timestamp WITH time zone NOT NULL,
    created_by uuid,
    updated timestamp WITH time zone NOT NULL,
    updated_by uuid,
    deleted_at timestamp WITH time zone,
    mood varchar(16) NOT NULL,
    note text,
    entry_type varchar(36) NOT NULL,
    share_with_health_worker boolean NOT NULL,
    shared_at timestamp WITH time zone,
    client_id uuid NOT NULL CONSTRAINT clients_healthdiarye_client_id_8c5de1dc_fk_clients_c REFERENCES clients_client deferrable initially deferred,
    organisation_id uuid NOT NULL CONSTRAINT clients_healthdiarye_organisation_id_78485008_fk_common_or REFERENCES common_organisation deferrable initially deferred
);

CREATE TABLE IF NOT EXISTS clients_healthdiaryattachment (
    id uuid NOT NULL PRIMARY KEY,
    active boolean NOT NULL,
    created timestamp WITH time zone NOT NULL,
    created_by uuid,
    updated timestamp WITH time zone NOT NULL,
    updated_by uuid,
    deleted_at timestamp WITH time zone,
    content_type varchar(100) NOT NULL,
    data varchar(65535) NOT NULL,
    title varchar(255) NOT NULL,
    creation_date timestamp WITH time zone NOT NULL,
    size integer,
    description text,
    aspect_ratio varchar(50),
    health_diary_entry_id uuid NOT NULL CONSTRAINT clients_healthdiarya_health_diary_entry_i_3b195bc0_fk_clients_h REFERENCES clients_healthdiaryentry deferrable initially deferred,
    organisation_id uuid NOT NULL CONSTRAINT clients_healthdiarya_organisation_id_128c56f5_fk_common_or REFERENCES common_organisation deferrable initially deferred
);

CREATE TABLE IF NOT EXISTS clients_servicerequest (
    id uuid NOT NULL PRIMARY KEY,
    active boolean NOT NULL,
    created timestamp WITH time zone NOT NULL,
    created_by uuid,
    updated timestamp WITH time zone NOT NULL,
    updated_by uuid,
    deleted_at timestamp WITH time zone,
    request_type varchar(36) NOT NULL,
    request text NOT NULL,
    STATUS varchar(36) NOT NULL,
    client_id uuid NOT NULL CONSTRAINT clients_servicerequest_client_id_640ced03_fk_clients_client_id REFERENCES clients_client deferrable initially deferred,
    organisation_id uuid NOT NULL CONSTRAINT clients_servicereque_organisation_id_8a95b043_fk_common_or REFERENCES common_organisation deferrable initially deferred,
    in_progress_at timestamp WITH time zone,
    in_progress_by_id uuid CONSTRAINT clients_servicereque_in_progress_by_id_c21e5f44_fk_staff_sta REFERENCES staff_staff deferrable initially deferred,
    resolved_at timestamp WITH time zone,
    resolved_by_id uuid CONSTRAINT clients_servicereque_resolved_by_id_492fdef8_fk_staff_sta REFERENCES staff_staff deferrable initially deferred,
    facility_id uuid CONSTRAINT clients_servicereque_facility_id_15a97356_fk_common_fa REFERENCES common_facility deferrable initially deferred,
    meta jsonb
);

CREATE TABLE IF NOT EXISTS staff_staff_addresses (
    id serial PRIMARY KEY,
    staff_id uuid NOT NULL CONSTRAINT staff_staff_addresses_staff_id_1620f1f1_fk_staff_staff_id REFERENCES staff_staff deferrable initially deferred,
    address_id uuid NOT NULL CONSTRAINT staff_staff_addresses_address_id_bdc155a3_fk_common_address_id REFERENCES common_address deferrable initially deferred,
    CONSTRAINT staff_staff_addresses_staff_id_address_id_44c681f2_uniq UNIQUE (staff_id, address_id)
);

CREATE TABLE IF NOT EXISTS staff_staff_contacts (
    id serial PRIMARY KEY,
    staff_id uuid NOT NULL CONSTRAINT staff_staff_contacts_staff_id_9e46919d_fk_staff_staff_id REFERENCES staff_staff deferrable initially deferred,
    contact_id uuid NOT NULL CONSTRAINT staff_staff_contacts_contact_id_1119fe01_fk_common_contact_id REFERENCES common_contact deferrable initially deferred,
    CONSTRAINT staff_staff_contacts_staff_id_contact_id_7d66d096_uniq UNIQUE (staff_id, contact_id)
);

CREATE TABLE IF NOT EXISTS staff_staff_identifiers (
    id serial PRIMARY KEY,
    staff_id uuid NOT NULL CONSTRAINT staff_staff_identifiers_staff_id_4cc212da_fk_staff_staff_id REFERENCES staff_staff deferrable initially deferred,
    identifier_id uuid NOT NULL CONSTRAINT staff_staff_identifi_identifier_id_eee797f4_fk_clients_i REFERENCES clients_identifier deferrable initially deferred,
    CONSTRAINT staff_staff_identifiers_staff_id_identifier_id_0606a0f6_uniq UNIQUE (staff_id, identifier_id)
);

CREATE TABLE IF NOT EXISTS appointments_appointment (
    id uuid NOT NULL PRIMARY KEY,
    active boolean NOT NULL,
    created timestamp WITH time zone NOT NULL,
    created_by uuid,
    updated timestamp WITH time zone NOT NULL,
    updated_by uuid,
    deleted_at timestamp WITH time zone,
    reason text,
    provider varchar(36),
    date date,
    client_id uuid NOT NULL CONSTRAINT appointments_appoint_client_id_c189cb9e_fk_clients_c REFERENCES clients_client deferrable initially deferred,
    organisation_id uuid NOT NULL CONSTRAINT appointments_appoint_organisation_id_44b89706_fk_common_or REFERENCES common_organisation deferrable initially deferred,
    staff_id uuid CONSTRAINT appointments_appointment_staff_id_64350136_fk_staff_staff_id REFERENCES staff_staff deferrable initially deferred,
    facility_id uuid CONSTRAINT appointments_appoint_facility_id_52ace4b7_fk_common_fa REFERENCES common_facility deferrable initially deferred,
    has_rescheduled_appointment boolean NOT NULL,
    external_id varchar(128) UNIQUE
);

CREATE TABLE IF NOT EXISTS authority_authorityrole_users (
    id serial PRIMARY KEY,
    authorityrole_id uuid NOT NULL CONSTRAINT authority_authorityr_authorityrole_id_5e240d69_fk_authority REFERENCES authority_authorityrole deferrable initially deferred,
    user_id uuid NOT NULL CONSTRAINT authority_authorityrole_users_user_id_d225811e_fk_users_user_id REFERENCES users_user deferrable initially deferred,
    CONSTRAINT authority_authorityrole__authorityrole_id_user_id_81c95b9c_uniq UNIQUE (authorityrole_id, user_id)
);

CREATE TABLE IF NOT EXISTS common_notification (
    id uuid NOT NULL PRIMARY KEY,
    active boolean NOT NULL,
    created timestamp WITH time zone NOT NULL,
    created_by uuid,
    updated timestamp WITH time zone NOT NULL,
    updated_by uuid,
    deleted_at timestamp WITH time zone,
    title varchar(64) NOT NULL,
    body text NOT NULL,
    notification_type varchar(32) NOT NULL,
    flavour varchar(32) NOT NULL,
    is_read boolean NOT NULL,
    facility_id uuid CONSTRAINT common_notification_facility_id_e51d5197_fk_common_facility_id REFERENCES common_facility deferrable initially deferred,
    organisation_id uuid NOT NULL CONSTRAINT common_notification_organisation_id_6d681c83_fk_common_or REFERENCES common_organisation deferrable initially deferred,
    user_id uuid CONSTRAINT common_notification_user_id_b851b5a2_fk_users_user_id REFERENCES users_user deferrable initially deferred
);

CREATE TABLE IF NOT EXISTS common_usersurveys (
    id uuid NOT NULL PRIMARY KEY,
    active boolean NOT NULL,
    created timestamp WITH time zone NOT NULL,
    created_by uuid,
    updated timestamp WITH time zone NOT NULL,
    updated_by uuid,
    deleted_at timestamp WITH time zone,
    link text NOT NULL,
    title text NOT NULL,
    description text,
    has_submitted boolean NOT NULL,
    organisation_id uuid NOT NULL CONSTRAINT common_usersurveys_organisation_id_83c7c493_fk_common_or REFERENCES common_organisation deferrable initially deferred,
    user_id uuid NOT NULL CONSTRAINT common_usersurveys_user_id_9c94c3eb_fk_users_user_id REFERENCES users_user deferrable initially deferred,
    form_id text,
    link_id integer,
    project_id integer,
    token text
);

CREATE TABLE IF NOT EXISTS common_feedback (
    id uuid NOT NULL PRIMARY KEY,
    active boolean NOT NULL,
    created timestamp WITH time zone NOT NULL,
    created_by uuid,
    updated timestamp WITH time zone NOT NULL,
    updated_by uuid,
    deleted_at timestamp WITH time zone,
    feedback_type varchar(32) NOT NULL,
    satisfaction_level integer,
    service_name varchar(32),
    feedback text,
    requires_followup boolean NOT NULL,
    organisation_id uuid NOT NULL CONSTRAINT common_feedback_organisation_id_cbd4e0e9_fk_common_or REFERENCES common_organisation deferrable initially deferred,
    user_id uuid NOT NULL CONSTRAINT common_feedback_user_id_3367f3cd_fk_users_user_id REFERENCES users_user deferrable initially deferred,
    phone_number varchar(32)
);

CREATE TABLE IF NOT EXISTS communities_community_clients (
    id serial PRIMARY KEY,
    community_id uuid NOT NULL CONSTRAINT communities_communit_community_id_a2523222_fk_communiti REFERENCES communities_community deferrable initially deferred,
    client_id uuid NOT NULL CONSTRAINT communities_communit_client_id_62435084_fk_clients_c REFERENCES clients_client deferrable initially deferred,
    CONSTRAINT communities_community_cl_community_id_client_id_9a2a0b79_uniq UNIQUE (community_id, client_id)
);

CREATE TABLE IF NOT EXISTS communities_community_staff (
    id serial PRIMARY KEY,
    community_id uuid NOT NULL CONSTRAINT communities_communit_community_id_c63e8e97_fk_communiti REFERENCES communities_community deferrable initially deferred,
    staff_id uuid NOT NULL CONSTRAINT communities_community_staff_staff_id_6ac91363_fk_staff_staff_id REFERENCES staff_staff deferrable initially deferred,
    CONSTRAINT communities_community_staff_community_id_staff_id_744e85cb_uniq UNIQUE (community_id, staff_id)
);

CREATE TABLE IF NOT EXISTS screeningtools_screeningtoolsresponse (
    id uuid NOT NULL PRIMARY KEY,
    active boolean NOT NULL,
    created timestamp WITH time zone NOT NULL,
    created_by uuid,
    updated timestamp WITH time zone NOT NULL,
    updated_by uuid,
    deleted_at timestamp WITH time zone,
    response text NOT NULL,
    client_id uuid NOT NULL CONSTRAINT screeningtools_scree_client_id_4daffd54_fk_clients_c REFERENCES clients_client deferrable initially deferred,
    organisation_id uuid NOT NULL CONSTRAINT screeningtools_scree_organisation_id_5a1fb2fd_fk_common_or REFERENCES common_organisation deferrable initially deferred,
    question_id uuid NOT NULL CONSTRAINT screeningtools_scree_question_id_a1dc424b_fk_screening REFERENCES screeningtools_screeningtoolsquestion deferrable initially deferred
);

CREATE TABLE IF NOT EXISTS staff_servicerequest (
    id uuid NOT NULL PRIMARY KEY,
    active boolean NOT NULL,
    created timestamp WITH time zone NOT NULL,
    created_by uuid,
    updated timestamp WITH time zone NOT NULL,
    updated_by uuid,
    deleted_at timestamp WITH time zone,
    request_type varchar(36) NOT NULL,
    request text NOT NULL,
    STATUS varchar(36) NOT NULL,
    resolved_at timestamp WITH time zone,
    meta jsonb,
    organisation_id uuid NOT NULL CONSTRAINT staff_servicerequest_organisation_id_4b38eea3_fk_common_or REFERENCES common_organisation deferrable initially deferred,
    resolved_by_id uuid CONSTRAINT staff_servicerequest_resolved_by_id_b680d0a6_fk_staff_staff_id REFERENCES staff_staff deferrable initially deferred,
    staff_id uuid NOT NULL CONSTRAINT staff_servicerequest_staff_id_d52e14cc_fk_staff_staff_id REFERENCES staff_staff deferrable initially deferred,
    facility_id uuid CONSTRAINT staff_servicerequest_facility_id_3c840b76_fk_common_facility_id REFERENCES common_facility deferrable initially deferred
);

CREATE TABLE IF NOT EXISTS users_metric (
    id serial PRIMARY KEY,
    timestamp timestamp WITH time zone NOT NULL,
    payload jsonb NOT NULL,
    metric_type varchar(32) NOT NULL,
    active boolean NOT NULL,
    created timestamp WITH time zone NOT NULL,
    created_by uuid,
    updated timestamp WITH time zone NOT NULL,
    updated_by uuid,
    user_id uuid NOT NULL CONSTRAINT users_metric_user_id_e7321b28_fk_users_user_id REFERENCES users_user deferrable initially deferred,
    deleted_at timestamp WITH time zone
);

CREATE TABLE IF NOT EXISTS users_userpin (
    id serial PRIMARY KEY,
    hashed_pin text NOT NULL,
    valid_from timestamp WITH time zone NOT NULL,
    valid_to timestamp WITH time zone NOT NULL,
    user_type varchar(32),
    active boolean NOT NULL,
    created timestamp WITH time zone NOT NULL,
    created_by uuid,
    updated timestamp WITH time zone NOT NULL,
    updated_by uuid,
    user_id uuid NOT NULL CONSTRAINT users_userpin_user_id_aba00bd1_fk_users_user_id REFERENCES users_user deferrable initially deferred,
    flavour varchar(32),
    salt text,
    deleted_at timestamp WITH time zone
);

CREATE TABLE IF NOT EXISTS users_userotp (
    id serial PRIMARY KEY,
    is_valid boolean NOT NULL,
    generated_at timestamp WITH time zone NOT NULL,
    valid_until timestamp WITH time zone,
    channel varchar(10) NOT NULL,
    flavour varchar(32),
    phonenumber text NOT NULL,
    user_id uuid NOT NULL CONSTRAINT users_userotp_user_id_7b41e28a_fk_users_user_id REFERENCES users_user deferrable initially deferred,
    created timestamp WITH time zone NOT NULL,
    created_by uuid,
    updated timestamp WITH time zone NOT NULL,
    updated_by uuid,
    otp varchar(8) NOT NULL,
    deleted_at timestamp WITH time zone
);

COMMIT;