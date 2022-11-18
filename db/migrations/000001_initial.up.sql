BEGIN;

CREATE TABLE IF NOT EXISTS "common_organisation"
(
    "id"               uuid PRIMARY KEY    NOT NULL,
    "active"           boolean             NOT NULL,
    "created"          timestamp           NOT NULL,
    "updated"          timestamp           NOT NULL,
    "created_by"       uuid,
    "updated_by"       uuid,
    "deleted_at"       timestamp,
    "code"             integer UNIQUE      NOT NULL,
    "org_code"         varchar(15) UNIQUE  NOT NULL,
    "name"             varchar(100) UNIQUE NOT NULL,
    "email_address"    varchar(100) UNIQUE NOT NULL,
    "phone_number"     varchar(128) UNIQUE NOT NULL,
    "description"      text,
    "postal_address"   varchar(100),
    "physical_address" text,
    "default_country"  varchar(255)        NOT NULL
);

CREATE TABLE IF NOT EXISTS "common_facility"
(
    "id"                   uuid PRIMARY KEY NOT NULL,
    "active"               boolean          NOT NULL,
    "created"              timestamp        NOT NULL,
    "created_by"           uuid,
    "updated"              timestamp        NOT NULL,
    "updated_by"           uuid,
    "deleted_at"           timestamp,
    "name"                 text             NOT NULL,
    "country"              varchar(64)      NOT NULL,
    "description"          text             NOT NULL,
    "phone"                varchar(15),
    "fhir_organization_id" varchar(64) UNIQUE
);

CREATE TABLE IF NOT EXISTS "common_facility_identifier"
(
    "id"               uuid PRIMARY KEY NOT NULL,
    "active"           boolean          NOT NULL,
    "created"          timestamp        NOT NULL,
    "created_by"       uuid,
    "updated"          timestamp        NOT NULL,
    "updated_by"       uuid,
    "deleted_at"       timestamp,
    "identifier_type"  text             NOT NULL,
    "identifier_value" text             NOT NULL,
    "facility_id"      uuid             NOT NULL,
    UNIQUE ("identifier_type", "identifier_value")
);

CREATE TABLE IF NOT EXISTS "users_user"
(
    "id"                           uuid PRIMARY KEY NOT NULL,
    "active"                       boolean          NOT NULL,
    "created"                      timestamp        NOT NULL,
    "created_by"                   uuid,
    "updated"                      timestamp        NOT NULL,
    "updated_by"                   uuid,
    "deleted_at"                   timestamp,
    "name"                         varchar(255)     NOT NULL,
    "username"                     varchar(150)     NOT NULL,
    "date_of_birth"                date,
    "gender"                       varchar(16),
    "push_tokens"                  varchar(256)[],
    "user_type"                    varchar(32),
    "flavour"                      varchar(32),
    "terms_accepted"               boolean          NOT NULL,
    "accepted_terms_of_service_id" integer,
    "languages"                    varchar(150)[],
    "avatar"                       text,
    "is_suspended"                 boolean          NOT NULL,
    "is_superuser"                 boolean          NOT NULL,
    "last_failed_login"            timestamp,
    "last_successful_login"        timestamp,
    "next_allowed_login"           timestamp        NOT NULL,
    "pin_change_required"          boolean,
    "has_set_pin"                  boolean,
    "has_set_security_questions"   boolean,
    "is_phone_verified"            boolean,
    "failed_login_count"           integer          NOT NULL,
    "failed_security_count"        integer          NOT NULL,
    "pin_update_required"          boolean,
    "has_set_username"             boolean,
    "current_program_id"           uuid,
    "organisation_id"              uuid             NOT NULL,
    UNIQUE ("username", "organisation_id")
);

CREATE TABLE IF NOT EXISTS "common_contact"
(
    "id"              uuid PRIMARY KEY NOT NULL,
    "active"          boolean          NOT NULL,
    "created"         timestamp        NOT NULL,
    "created_by"      uuid,
    "updated"         timestamp        NOT NULL,
    "updated_by"      uuid,
    "contact_type"    varchar(16)      NOT NULL,
    "contact_value"   text             NOT NULL,
    "opted_in"        boolean          NOT NULL,
    "deleted_at"      timestamp,
    "user_id"         uuid,
    "flavour"         varchar(32),
    "organisation_id" uuid             NOT NULL
);

CREATE TABLE IF NOT EXISTS "common_program"
(
    "id"              uuid PRIMARY KEY NOT NULL,
    "active"          boolean          NOT NULL,
    "created"         timestamp        NOT NULL,
    "created_by"      uuid,
    "updated"         timestamp        NOT NULL,
    "updated_by"      uuid,
    "deleted_at"      timestamp,
    "name"            text,
    "client_types"    varchar(64)[]    NOT NULL,
    "organisation_id" uuid             NOT NULL,
    UNIQUE ("organisation_id", "name")
);

CREATE TABLE IF NOT EXISTS "common_program_user"
(
    "id"         SERIAL PRIMARY KEY,
    "user_id"    uuid NOT NULL,
    "program_id" uuid NOT NULL,
    UNIQUE ("user_id", "program_id")
);

CREATE TABLE IF NOT EXISTS "common_program_facility"
(
    "id"          SERIAL PRIMARY KEY,
    "facility_id" uuid NOT NULL,
    "program_id"  uuid NOT NULL,
    UNIQUE ("facility_id", "program_id")
);

CREATE TABLE IF NOT EXISTS "common_auditlog"
(
    "id"              uuid PRIMARY KEY NOT NULL,
    "active"          boolean          NOT NULL,
    "created"         timestamp        NOT NULL,
    "created_by"      uuid,
    "updated"         timestamp        NOT NULL,
    "updated_by"      uuid,
    "deleted_at"      timestamp,
    "timestamp"       timestamp        NOT NULL,
    "record_type"     text             NOT NULL,
    "notes"           text             NOT NULL,
    "payload"         jsonb            NOT NULL,
    "organisation_id" uuid             NOT NULL
);

CREATE TABLE IF NOT EXISTS "common_address"
(
    "id"              uuid PRIMARY KEY NOT NULL,
    "active"          boolean          NOT NULL,
    "created"         timestamp        NOT NULL,
    "created_by"      uuid,
    "updated"         timestamp        NOT NULL,
    "updated_by"      uuid,
    "deleted_at"      timestamp,
    "address_type"    varchar(16)      NOT NULL,
    "text"            text             NOT NULL,
    "postal_code"     text             NOT NULL,
    "country"         varchar(255)     NOT NULL,
    "organisation_id" uuid             NOT NULL
);

CREATE TABLE IF NOT EXISTS "common_securityquestion"
(
    "id"            uuid PRIMARY KEY NOT NULL,
    "active"        boolean          NOT NULL,
    "created"       timestamp        NOT NULL,
    "created_by"    uuid,
    "updated"       timestamp        NOT NULL,
    "updated_by"    uuid,
    "deleted_at"    timestamp,
    "stem"          text             NOT NULL,
    "description"   text             NOT NULL,
    "sequence"      integer          NOT NULL,
    "response_type" varchar(32)      NOT NULL,
    "flavour"       varchar(32)
);

CREATE TABLE IF NOT EXISTS "clients_relatedperson"
(
    "id"                uuid PRIMARY KEY NOT NULL,
    "active"            boolean          NOT NULL,
    "created"           timestamp        NOT NULL,
    "created_by"        uuid,
    "updated"           timestamp        NOT NULL,
    "updated_by"        uuid,
    "deleted_at"        timestamp,
    "first_name"        text             NOT NULL,
    "last_name"         text             NOT NULL,
    "other_name"        text             NOT NULL,
    "date_of_birth"     date,
    "gender"            varchar(16)      NOT NULL,
    "relationship_type" varchar(64)      NOT NULL,
    "organisation_id"   uuid             NOT NULL,
    "program_id"        uuid             NOT NULL
);

CREATE TABLE IF NOT EXISTS "clients_relatedperson_addresses"
(
    "id"               SERIAL PRIMARY KEY,
    "relatedperson_id" uuid NOT NULL,
    "address_id"       uuid NOT NULL,
    UNIQUE ("relatedperson_id", "address_id")
);

CREATE TABLE IF NOT EXISTS "clients_identifier"
(
    "id"                    uuid PRIMARY KEY NOT NULL,
    "active"                boolean          NOT NULL,
    "created"               timestamp        NOT NULL,
    "created_by"            uuid,
    "updated"               timestamp        NOT NULL,
    "updated_by"            uuid,
    "deleted_at"            timestamp,
    "identifier_type"       varchar(64)      NOT NULL,
    "identifier_value"      text             NOT NULL,
    "identifier_use"        varchar(64)      NOT NULL,
    "description"           text             NOT NULL,
    "valid_from"            timestamp        NOT NULL,
    "valid_to"              timestamp,
    "is_primary_identifier" boolean          NOT NULL,
    "organisation_id"       uuid             NOT NULL,
    "program_id"            uuid             NOT NULL,
    UNIQUE (
            "program_id",
            "identifier_type",
            "identifier_value"
        )
);

CREATE TABLE IF NOT EXISTS "clients_healthdiaryquote"
(
    "id"              uuid PRIMARY KEY NOT NULL,
    "active"          boolean          NOT NULL,
    "created"         timestamp        NOT NULL,
    "created_by"      uuid,
    "updated"         timestamp        NOT NULL,
    "updated_by"      uuid,
    "deleted_at"      timestamp,
    "quote"           text             NOT NULL,
    "by"              text             NOT NULL,
    "program_id"      uuid             NOT NULL,
    "organisation_id" uuid             NOT NULL,
    UNIQUE ("quote", "program_id")
);

CREATE TABLE IF NOT EXISTS "authority_authoritypermission"
(
    "id"              uuid PRIMARY KEY NOT NULL,
    "active"          boolean          NOT NULL,
    "created"         timestamp        NOT NULL,
    "created_by"      uuid,
    "updated"         timestamp        NOT NULL,
    "updated_by"      uuid,
    "deleted_at"      timestamp,
    "name"            varchar(255)     NOT NULL,
    "organisation_id" uuid             NOT NULL
);

CREATE TABLE IF NOT EXISTS "authority_authorityrole"
(
    "id"              uuid PRIMARY KEY NOT NULL,
    "active"          boolean          NOT NULL,
    "created"         timestamp        NOT NULL,
    "created_by"      uuid,
    "updated"         timestamp        NOT NULL,
    "updated_by"      uuid,
    "deleted_at"      timestamp,
    "name"            varchar(255)     NOT NULL,
    "organisation_id" uuid             NOT NULL
);

CREATE TABLE IF NOT EXISTS "authority_authorityrole_permissions"
(
    "id"                     SERIAL PRIMARY KEY,
    "authorityrole_id"       uuid NOT NULL,
    "authoritypermission_id" uuid NOT NULL,
    UNIQUE ("authorityrole_id", "authoritypermission_id")
);

CREATE TABLE IF NOT EXISTS "communities_community"
(
    "id"              uuid PRIMARY KEY NOT NULL,
    "active"          boolean          NOT NULL,
    "created"         timestamp        NOT NULL,
    "created_by"      uuid,
    "updated"         timestamp        NOT NULL,
    "updated_by"      uuid,
    "deleted_at"      timestamp,
    "name"            varchar(64)      NOT NULL,
    "description"     varchar(150)     NOT NULL,
    "client_types"    varchar(64)[]    NOT NULL,
    "gender"          varchar(16)[],
    "invite_only"     boolean          NOT NULL,
    "discoverable"    boolean          NOT NULL,
    "max_age"         integer,
    "min_age"         integer,
    "program_id"      uuid             NOT NULL,
    "organisation_id" uuid             NOT NULL
);

CREATE TABLE IF NOT EXISTS "screeningtools_screeningtoolsquestion"
(
    "id"                uuid PRIMARY KEY NOT NULL,
    "active"            boolean          NOT NULL,
    "created"           timestamp        NOT NULL,
    "created_by"        uuid,
    "updated"           timestamp        NOT NULL,
    "updated_by"        uuid,
    "deleted_at"        timestamp,
    "question"          text             NOT NULL,
    "tool_type"         varchar(32)      NOT NULL,
    "response_choices"  jsonb,
    "response_type"     varchar(32)      NOT NULL,
    "response_category" varchar(32)      NOT NULL,
    "sequence"          integer          NOT NULL,
    "meta"              jsonb,
    "organisation_id"   uuid             NOT NULL,
    "program_id"        uuid             NOT NULL
);

CREATE TABLE IF NOT EXISTS "users_termsofservice"
(
    "id"         SERIAL PRIMARY KEY,
    "text"       text      NOT NULL,
    "valid_from" timestamp NOT NULL,
    "valid_to"   timestamp,
    "active"     boolean   NOT NULL,
    "created"    timestamp NOT NULL,
    "created_by" uuid,
    "updated"    timestamp NOT NULL,
    "updated_by" uuid,
    "deleted_at" timestamp,
    "flavour"    varchar(32)
);

CREATE TABLE IF NOT EXISTS "staff_staff"
(
    "id"                  uuid PRIMARY KEY NOT NULL,
    "active"              boolean          NOT NULL,
    "created"             timestamp        NOT NULL,
    "created_by"          uuid,
    "updated"             timestamp        NOT NULL,
    "updated_by"          uuid,
    "deleted_at"          timestamp,
    "staff_number"        varchar(255),
    "current_facility_id" uuid             NOT NULL,
    "organisation_id"     uuid             NOT NULL,
    "user_id"             uuid             NOT NULL,
    "program_id"          uuid             NOT NULL,
    UNIQUE ("user_id", "program_id")
);

CREATE TABLE IF NOT EXISTS "staff_staff_facilities"
(
    "id"          SERIAL PRIMARY KEY,
    "staff_id"    uuid NOT NULL,
    "facility_id" uuid NOT NULL,
    UNIQUE ("staff_id", "facility_id")
);

CREATE TABLE IF NOT EXISTS "clients_client"
(
    "id"                   uuid PRIMARY KEY NOT NULL,
    "active"               boolean          NOT NULL,
    "created"              timestamp        NOT NULL,
    "created_by"           uuid,
    "updated"              timestamp        NOT NULL,
    "updated_by"           uuid,
    "deleted_at"           timestamp,
    "client_types"         varchar(64)[]    NOT NULL,
    "enrollment_date"      timestamp        NOT NULL,
    "fhir_patient_id"      text,
    "emr_health_record_id" text,
    "counselled"           boolean          NOT NULL,
    "current_facility_id"  uuid             NOT NULL,
    "organisation_id"      uuid             NOT NULL,
    "user_id"              uuid             NOT NULL,
    "program_id"           uuid             NOT NULL,
    UNIQUE ("program_id", "fhir_patient_id", "user_id")
);

CREATE TABLE IF NOT EXISTS "clients_relatedperson_contacts"
(
    "id"               SERIAL PRIMARY KEY,
    "relatedperson_id" uuid NOT NULL,
    "contact_id"       uuid NOT NULL,
    UNIQUE ("contact_id", "relatedperson_id")
);

CREATE TABLE IF NOT EXISTS "clients_clientfacility"
(
    "id"              uuid PRIMARY KEY NOT NULL,
    "active"          boolean          NOT NULL,
    "created"         timestamp        NOT NULL,
    "created_by"      uuid,
    "updated"         timestamp        NOT NULL,
    "updated_by"      uuid,
    "assigned"        timestamp        NOT NULL,
    "transferred_out" timestamp,
    "client_id"       uuid             NOT NULL,
    "facility_id"     uuid             NOT NULL,
    "organisation_id" uuid             NOT NULL,
    "notes"           text             NOT NULL,
    "deleted_at"      timestamp,
    "program_id"      uuid             NOT NULL,
    UNIQUE (
            "organisation_id",
            "program_id",
            "client_id",
            "facility_id"
        )
);

CREATE TABLE IF NOT EXISTS "clients_client_identifiers"
(
    "id"            SERIAL PRIMARY KEY,
    "client_id"     uuid NOT NULL,
    "identifier_id" uuid NOT NULL,
    UNIQUE ("client_id", "identifier_id")
);

CREATE TABLE IF NOT EXISTS "clients_client_related_persons"
(
    "id"               SERIAL PRIMARY KEY,
    "client_id"        uuid NOT NULL,
    "relatedperson_id" uuid NOT NULL,
    UNIQUE ("client_id", "relatedperson_id")
);

CREATE TABLE IF NOT EXISTS "common_securityquestionresponse"
(
    "id"              uuid PRIMARY KEY NOT NULL,
    "active"          boolean          NOT NULL,
    "created"         timestamp        NOT NULL,
    "created_by"      uuid,
    "updated"         timestamp        NOT NULL,
    "updated_by"      uuid,
    "deleted_at"      timestamp,
    "timestamp"       timestamp        NOT NULL,
    "response"        text             NOT NULL,
    "is_correct"      boolean          NOT NULL,
    "organisation_id" uuid             NOT NULL,
    "question_id"     uuid             NOT NULL,
    "user_id"         uuid             NOT NULL,
    UNIQUE ("organisation_id", "question_id", "user_id")
);

CREATE TABLE IF NOT EXISTS "clients_healthdiaryentry"
(
    "id"                       uuid PRIMARY KEY NOT NULL,
    "active"                   boolean          NOT NULL,
    "created"                  timestamp        NOT NULL,
    "created_by"               uuid,
    "updated"                  timestamp        NOT NULL,
    "updated_by"               uuid,
    "deleted_at"               timestamp,
    "mood"                     varchar(16)      NOT NULL,
    "note"                     text,
    "entry_type"               varchar(36)      NOT NULL,
    "share_with_health_worker" boolean          NOT NULL,
    "shared_at"                timestamp,
    "client_id"                uuid             NOT NULL,
    "organisation_id"          uuid             NOT NULL,
    "program_id"               uuid             NOT NULL
);

CREATE TABLE IF NOT EXISTS "clients_servicerequest"
(
    "id"                uuid PRIMARY KEY NOT NULL,
    "active"            boolean          NOT NULL,
    "created"           timestamp        NOT NULL,
    "created_by"        uuid,
    "updated"           timestamp        NOT NULL,
    "updated_by"        uuid,
    "deleted_at"        timestamp,
    "request_type"      varchar(36)      NOT NULL,
    "request"           text             NOT NULL,
    "status"            varchar(36)      NOT NULL,
    "in_progress_at"    timestamp,
    "in_progress_by_id" uuid,
    "resolved_at"       timestamp,
    "resolved_by_id"    uuid,
    "facility_id"       uuid,
    "meta"              jsonb,
    "client_id"         uuid             NOT NULL,
    "organisation_id"   uuid             NOT NULL,
    "program_id"        uuid             NOT NULL
);

CREATE TABLE IF NOT EXISTS "staff_staff_identifiers"
(
    "id"            SERIAL PRIMARY KEY,
    "staff_id"      uuid NOT NULL,
    "identifier_id" uuid NOT NULL,
    UNIQUE ("staff_id", "identifier_id")
);

CREATE TABLE IF NOT EXISTS "appointments_appointment"
(
    "id"                          uuid PRIMARY KEY NOT NULL,
    "active"                      boolean          NOT NULL,
    "created"                     timestamp        NOT NULL,
    "created_by"                  uuid,
    "updated"                     timestamp        NOT NULL,
    "updated_by"                  uuid,
    "deleted_at"                  timestamp,
    "reason"                      text,
    "provider"                    varchar(36),
    "date"                        date,
    "has_rescheduled_appointment" boolean          NOT NULL,
    "external_id"                 varchar(128) UNIQUE,
    "client_id"                   uuid             NOT NULL,
    "organisation_id"             uuid             NOT NULL,
    "staff_id"                    uuid,
    "facility_id"                 uuid,
    "program_id"                  uuid             NOT NULL
);

CREATE TABLE IF NOT EXISTS "authority_authorityrole_users"
(
    "id"               SERIAL PRIMARY KEY,
    "authorityrole_id" uuid NOT NULL,
    "user_id"          uuid NOT NULL,
    UNIQUE ("user_id", "authorityrole_id")
);

CREATE TABLE IF NOT EXISTS "common_notification"
(
    "id"                uuid PRIMARY KEY NOT NULL,
    "active"            boolean          NOT NULL,
    "created"           timestamp        NOT NULL,
    "created_by"        uuid,
    "updated"           timestamp        NOT NULL,
    "updated_by"        uuid,
    "deleted_at"        timestamp,
    "title"             varchar(64)      NOT NULL,
    "body"              text             NOT NULL,
    "notification_type" varchar(32)      NOT NULL,
    "flavour"           varchar(32)      NOT NULL,
    "is_read"           boolean          NOT NULL,
    "facility_id"       uuid,
    "organisation_id"   uuid             NOT NULL,
    "user_id"           uuid,
    "program_id"        uuid             NOT NULL
);

CREATE TABLE IF NOT EXISTS "common_usersurveys"
(
    "id"              uuid PRIMARY KEY NOT NULL,
    "active"          boolean          NOT NULL,
    "created"         timestamp        NOT NULL,
    "created_by"      uuid,
    "updated"         timestamp        NOT NULL,
    "updated_by"      uuid,
    "deleted_at"      timestamp,
    "link"            text             NOT NULL,
    "title"           text             NOT NULL,
    "description"     text,
    "has_submitted"   boolean          NOT NULL,
    "organisation_id" uuid             NOT NULL,
    "user_id"         uuid             NOT NULL,
    "form_id"         text,
    "link_id"         integer,
    "project_id"      integer,
    "token"           text,
    "submitted_at"    timestamp,
    "program_id"      uuid             NOT NULL
);

CREATE TABLE IF NOT EXISTS "common_feedback"
(
    "id"                 uuid PRIMARY KEY NOT NULL,
    "active"             boolean          NOT NULL,
    "created"            timestamp        NOT NULL,
    "created_by"         uuid,
    "updated"            timestamp        NOT NULL,
    "updated_by"         uuid,
    "deleted_at"         timestamp,
    "feedback_type"      varchar(32)      NOT NULL,
    "satisfaction_level" integer,
    "service_name"       varchar(32),
    "feedback"           text,
    "requires_followup"  boolean          NOT NULL,
    "organisation_id"    uuid             NOT NULL,
    "user_id"            uuid             NOT NULL,
    "phone_number"       varchar(32),
    "program_id"         uuid             NOT NULL
);

CREATE TABLE IF NOT EXISTS "communities_community_clients"
(
    "id"           SERIAL PRIMARY KEY,
    "community_id" uuid NOT NULL,
    "client_id"    uuid NOT NULL,
    UNIQUE ("community_id", "client_id")
);

CREATE TABLE IF NOT EXISTS "communities_community_staff"
(
    "id"           SERIAL PRIMARY KEY,
    "community_id" uuid NOT NULL,
    "staff_id"     uuid NOT NULL,
    UNIQUE ("community_id", "staff_id")
);

CREATE TABLE IF NOT EXISTS "screeningtools_screeningtoolsresponse"
(
    "id"              uuid PRIMARY KEY NOT NULL,
    "active"          boolean          NOT NULL,
    "created"         timestamp        NOT NULL,
    "created_by"      uuid,
    "updated"         timestamp        NOT NULL,
    "updated_by"      uuid,
    "deleted_at"      timestamp,
    "response"        text             NOT NULL,
    "client_id"       uuid             NOT NULL,
    "organisation_id" uuid             NOT NULL,
    "question_id"     uuid             NOT NULL,
    "program_id"      uuid             NOT NULL
);

CREATE TABLE IF NOT EXISTS "staff_servicerequest"
(
    "id"              uuid PRIMARY KEY NOT NULL,
    "active"          boolean          NOT NULL,
    "created"         timestamp        NOT NULL,
    "created_by"      uuid,
    "updated"         timestamp        NOT NULL,
    "updated_by"      uuid,
    "deleted_at"      timestamp,
    "request_type"    varchar(36)      NOT NULL,
    "request"         text             NOT NULL,
    "status"          varchar(36)      NOT NULL,
    "resolved_at"     timestamp,
    "meta"            jsonb,
    "organisation_id" uuid             NOT NULL,
    "resolved_by_id"  uuid,
    "staff_id"        uuid             NOT NULL,
    "facility_id"     uuid,
    "program_id"      uuid             NOT NULL
);

CREATE TABLE IF NOT EXISTS "users_metric"
(
    "id"          SERIAL PRIMARY KEY,
    "timestamp"   timestamp   NOT NULL,
    "payload"     jsonb       NOT NULL,
    "metric_type" varchar(32) NOT NULL,
    "active"      boolean     NOT NULL,
    "created"     timestamp   NOT NULL,
    "created_by"  uuid,
    "updated"     timestamp   NOT NULL,
    "updated_by"  uuid,
    "user_id"     uuid        NOT NULL,
    "deleted_at"  timestamp
);

CREATE TABLE IF NOT EXISTS "users_userpin"
(
    "id"              SERIAL PRIMARY KEY,
    "hashed_pin"      text      NOT NULL,
    "valid_from"      timestamp NOT NULL,
    "valid_to"        timestamp NOT NULL,
    "user_type"       varchar(32),
    "active"          boolean   NOT NULL,
    "created"         timestamp NOT NULL,
    "created_by"      uuid,
    "updated"         timestamp NOT NULL,
    "updated_by"      uuid,
    "user_id"         uuid      NOT NULL,
    "flavour"         varchar(32),
    "salt"            text,
    "deleted_at"      timestamp,
    "organisation_id" uuid      NOT NULL
);

CREATE TABLE IF NOT EXISTS "users_userotp"
(
    "id"           SERIAL PRIMARY KEY,
    "is_valid"     boolean     NOT NULL,
    "generated_at" timestamp   NOT NULL,
    "valid_until"  timestamp,
    "channel"      varchar(10) NOT NULL,
    "flavour"      varchar(32),
    "phonenumber"  text        NOT NULL,
    "user_id"      uuid        NOT NULL,
    "created"      timestamp   NOT NULL,
    "created_by"   uuid,
    "updated"      timestamp   NOT NULL,
    "updated_by"   uuid,
    "otp"          varchar(8)  NOT NULL,
    "deleted_at"   timestamp
);

CREATE TABLE IF NOT EXISTS "questionnaires_questionnaire"
(
    "id"              uuid PRIMARY KEY NOT NULL,
    "active"          boolean          NOT NULL,
    "created"         timestamp        NOT NULL,
    "created_by"      uuid,
    "updated"         timestamp        NOT NULL,
    "updated_by"      uuid,
    "deleted_at"      timestamp,
    "organisation_id" uuid             NOT NULL,
    "name"            varchar(64)      NOT NULL,
    "description"     text             NOT NULL,
    "program_id"      uuid             NOT NULL
);

CREATE TABLE IF NOT EXISTS "questionnaires_screeningtool"
(
    "id"               uuid PRIMARY KEY NOT NULL,
    "active"           boolean          NOT NULL,
    "created"          timestamp        NOT NULL,
    "created_by"       uuid,
    "updated"          timestamp        NOT NULL,
    "updated_by"       uuid,
    "deleted_at"       timestamp,
    "organisation_id"  uuid             NOT NULL,
    "questionnaire_id" uuid             NOT NULL,
    "threshold"        integer DEFAULT 0,
    "client_types"     varchar(64)[],
    "genders"          varchar(64)[],
    "min_age"          integer,
    "max_age"          integer,
    "program_id"       uuid             NOT NULL
);

CREATE TABLE IF NOT EXISTS "questionnaires_question"
(
    "id"                  uuid PRIMARY KEY NOT NULL,
    "active"              boolean          NOT NULL,
    "created"             timestamp        NOT NULL,
    "created_by"          uuid,
    "updated"             timestamp        NOT NULL,
    "updated_by"          uuid,
    "deleted_at"          timestamp,
    "organisation_id"     uuid             NOT NULL,
    "questionnaire_id"    uuid             NOT NULL,
    "text"                text             NOT NULL,
    "question_type"       varchar(64)      NOT NULL,
    "response_value_type" varchar(64)      NOT NULL,
    "select_multiple"     boolean,
    "required"            boolean          NOT NULL,
    "sequence"            integer          NOT NULL,
    "program_id"          uuid             NOT NULL
);

CREATE TABLE IF NOT EXISTS "questionnaires_questioninputchoice"
(
    "id"              uuid PRIMARY KEY NOT NULL,
    "active"          boolean          NOT NULL,
    "created"         timestamp        NOT NULL,
    "created_by"      uuid,
    "updated"         timestamp        NOT NULL,
    "updated_by"      uuid,
    "deleted_at"      timestamp,
    "organisation_id" uuid             NOT NULL,
    "question_id"     uuid             NOT NULL,
    "choice"          varchar(64),
    "value"           text             NOT NULL,
    "score"           integer DEFAULT 0,
    "program_id"      uuid             NOT NULL
);

CREATE TABLE IF NOT EXISTS "questionnaires_screeningtoolresponse"
(
    "id"               uuid PRIMARY KEY NOT NULL,
    "active"           boolean          NOT NULL,
    "created"          timestamp        NOT NULL,
    "created_by"       uuid,
    "updated"          timestamp        NOT NULL,
    "updated_by"       uuid,
    "deleted_at"       timestamp,
    "organisation_id"  uuid             NOT NULL,
    "screeningtool_id" uuid             NOT NULL,
    "facility_id"      uuid             NOT NULL,
    "client_id"        uuid             NOT NULL,
    "aggregate_score"  integer DEFAULT 0,
    "program_id"       uuid             NOT NULL
);

CREATE TABLE IF NOT EXISTS "questionnaires_screeningtoolquestionresponse"
(
    "id"                       uuid PRIMARY KEY NOT NULL,
    "active"                   boolean          NOT NULL,
    "created"                  timestamp        NOT NULL,
    "created_by"               uuid,
    "updated"                  timestamp        NOT NULL,
    "updated_by"               uuid,
    "deleted_at"               timestamp,
    "organisation_id"          uuid             NOT NULL,
    "screeningtoolresponse_id" uuid             NOT NULL,
    "question_id"              uuid             NOT NULL,
    "response"                 text             NOT NULL,
    "score"                    integer DEFAULT 0,
    "program_id"               uuid             NOT NULL,
    "facility_id"              uuid             NOT NULL
);

CREATE TABLE IF NOT EXISTS "clients_client_facilities"
(
    "id"          SERIAL PRIMARY KEY,
    "client_id"   uuid NOT NULL,
    "facility_id" uuid NOT NULL,
    UNIQUE ("client_id", "facility_id")
);

CREATE TABLE IF NOT EXISTS "caregivers_caregiver"
(
    "id"               uuid PRIMARY KEY NOT NULL,
    "active"           boolean          NOT NULL,
    "created"          timestamp        NOT NULL,
    "created_by"       uuid,
    "updated"          timestamp        NOT NULL,
    "updated_by"       uuid,
    "deleted_at"       timestamp,
    "organisation_id"  uuid             NOT NULL,
    "user_id"          uuid             NOT NULL,
    "caregiver_number" varchar(255),
    "program_id"       uuid             NOT NULL,
    UNIQUE ("user_id", "program_id")
);

CREATE TABLE IF NOT EXISTS "caregivers_caregiver_client"
(
    "id"                   SERIAL PRIMARY KEY,
    "caregiver_id"         uuid        NOT NULL,
    "client_id"            uuid        NOT NULL,
    "active"               boolean     NOT NULL,
    "created"              timestamp   NOT NULL,
    "created_by"           uuid,
    "updated"              timestamp   NOT NULL,
    "updated_by"           uuid,
    "deleted_at"           timestamp,
    "organisation_id"      uuid        NOT NULL,
    "assigned_by"          uuid        NOT NULL,
    "relationship_type"    varchar(32) NOT NULL,
    "caregiver_consent"    varchar(255),
    "caregiver_consent_at" timestamp,
    "client_consent"       varchar(255),
    "client_consent_at"    timestamp,
    "program_id"           uuid        NOT NULL,
    UNIQUE ("caregiver_id", "client_id")
);

ALTER TABLE
    IF EXISTS "common_organisation"
    ADD
        CONSTRAINT "common_organisation_created_by_fkey" FOREIGN KEY ("created_by") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "common_organisation"
    ADD
        CONSTRAINT "common_organisation_updated_by_fkey" FOREIGN KEY ("updated_by") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "common_facility"
    ADD
        CONSTRAINT "common_facility_created_by_fkey" FOREIGN KEY ("created_by") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "common_facility"
    ADD
        CONSTRAINT "common_facility_updated_by_fkey" FOREIGN KEY ("updated_by") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "common_facility_identifier"
    ADD
        CONSTRAINT "common_facility_identifier_created_by_fkey" FOREIGN KEY ("created_by") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "common_facility_identifier"
    ADD
        CONSTRAINT "common_facility_identifier_updated_by_fkey" FOREIGN KEY ("updated_by") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "common_facility_identifier"
    ADD
        CONSTRAINT "common_facility_identifier_facility_id_fkey" FOREIGN KEY ("facility_id") REFERENCES "common_facility" ("id");

ALTER TABLE
    IF EXISTS "users_user"
    ADD
        CONSTRAINT "users_user_created_by_fkey" FOREIGN KEY ("created_by") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "users_user"
    ADD
        CONSTRAINT "users_user_updated_by_fkey" FOREIGN KEY ("updated_by") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "users_user"
    ADD
        CONSTRAINT "users_user_accepted_terms_of_service_id_fkey" FOREIGN KEY ("accepted_terms_of_service_id") REFERENCES "users_termsofservice" ("id");

ALTER TABLE
    IF EXISTS "users_user"
    ADD
        CONSTRAINT "users_user_current_program_id_fkey" FOREIGN KEY ("current_program_id") REFERENCES "common_program" ("id");

ALTER TABLE
    IF EXISTS "users_user"
    ADD
        CONSTRAINT "users_user_organisation_id_fkey" FOREIGN KEY ("organisation_id") REFERENCES "common_organisation" ("id");

ALTER TABLE
    IF EXISTS "common_contact"
    ADD
        CONSTRAINT "common_contact_created_by_fkey" FOREIGN KEY ("created_by") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "common_contact"
    ADD
        CONSTRAINT "common_contact_updated_by_fkey" FOREIGN KEY ("updated_by") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "common_contact"
    ADD
        CONSTRAINT "common_contact_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "common_contact"
    ADD
        CONSTRAINT "common_contact_organisation_id_fkey" FOREIGN KEY ("organisation_id") REFERENCES "common_organisation" ("id");

ALTER TABLE
    IF EXISTS "common_program"
    ADD
        CONSTRAINT "common_program_created_by_fkey" FOREIGN KEY ("created_by") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "common_program"
    ADD
        CONSTRAINT "common_program_updated_by_fkey" FOREIGN KEY ("updated_by") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "common_program"
    ADD
        CONSTRAINT "common_program_organisation_id_fkey" FOREIGN KEY ("organisation_id") REFERENCES "common_organisation" ("id");

ALTER TABLE
    IF EXISTS "common_program_user"
    ADD
        CONSTRAINT "common_program_user_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "common_program_user"
    ADD
        CONSTRAINT "common_program_user_program_id_fkey" FOREIGN KEY ("program_id") REFERENCES "common_program" ("id");

ALTER TABLE
    IF EXISTS "common_program_facility"
    ADD
        CONSTRAINT "common_program_facility_facility_id_fkey" FOREIGN KEY ("facility_id") REFERENCES "common_facility" ("id");

ALTER TABLE
    IF EXISTS "common_program_facility"
    ADD
        CONSTRAINT "common_program_facility_program_id_fkey" FOREIGN KEY ("program_id") REFERENCES "common_program" ("id");

ALTER TABLE
    IF EXISTS "common_auditlog"
    ADD
        CONSTRAINT "common_auditlog_created_by_fkey" FOREIGN KEY ("created_by") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "common_auditlog"
    ADD
        CONSTRAINT "common_auditlog_updated_by_fkey" FOREIGN KEY ("updated_by") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "common_auditlog"
    ADD
        CONSTRAINT "common_auditlog_organisation_id_fkey" FOREIGN KEY ("organisation_id") REFERENCES "common_organisation" ("id");

ALTER TABLE
    IF EXISTS "common_address"
    ADD
        CONSTRAINT "common_address_created_by_fkey" FOREIGN KEY ("created_by") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "common_address"
    ADD
        CONSTRAINT "common_address_updated_by_fkey" FOREIGN KEY ("updated_by") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "common_address"
    ADD
        CONSTRAINT "common_address_organisation_id_fkey" FOREIGN KEY ("organisation_id") REFERENCES "common_organisation" ("id");

ALTER TABLE
    IF EXISTS "common_securityquestion"
    ADD
        CONSTRAINT "common_securityquestion_created_by_fkey" FOREIGN KEY ("created_by") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "common_securityquestion"
    ADD
        CONSTRAINT "common_securityquestion_updated_by_fkey" FOREIGN KEY ("updated_by") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "clients_relatedperson"
    ADD
        CONSTRAINT "clients_relatedperson_created_by_fkey" FOREIGN KEY ("created_by") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "clients_relatedperson"
    ADD
        CONSTRAINT "clients_relatedperson_updated_by_fkey" FOREIGN KEY ("updated_by") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "clients_relatedperson"
    ADD
        CONSTRAINT "clients_relatedperson_organisation_id_fkey" FOREIGN KEY ("organisation_id") REFERENCES "common_organisation" ("id");

ALTER TABLE
    IF EXISTS "clients_relatedperson"
    ADD
        CONSTRAINT "clients_relatedperson_program_id_fkey" FOREIGN KEY ("program_id") REFERENCES "common_program" ("id");

ALTER TABLE
    IF EXISTS "clients_relatedperson_addresses"
    ADD
        CONSTRAINT "clients_relatedperson_addresses_relatedperson_id_fkey" FOREIGN KEY ("relatedperson_id") REFERENCES "clients_relatedperson" ("id");

ALTER TABLE
    IF EXISTS "clients_relatedperson_addresses"
    ADD
        CONSTRAINT "clients_relatedperson_addresses_address_id_fkey" FOREIGN KEY ("address_id") REFERENCES "common_address" ("id");

ALTER TABLE
    IF EXISTS "clients_identifier"
    ADD
        CONSTRAINT "clients_identifier_created_by_fkey" FOREIGN KEY ("created_by") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "clients_identifier"
    ADD
        CONSTRAINT "clients_identifier_updated_by_fkey" FOREIGN KEY ("updated_by") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "clients_identifier"
    ADD
        CONSTRAINT "clients_identifier_organisation_id_fkey" FOREIGN KEY ("organisation_id") REFERENCES "common_organisation" ("id");

ALTER TABLE
    IF EXISTS "clients_identifier"
    ADD
        CONSTRAINT "clients_identifier_program_id_fkey" FOREIGN KEY ("program_id") REFERENCES "common_program" ("id");

ALTER TABLE
    IF EXISTS "clients_healthdiaryquote"
    ADD
        CONSTRAINT "clients_healthdiaryquote_created_by_fkey" FOREIGN KEY ("created_by") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "clients_healthdiaryquote"
    ADD
        CONSTRAINT "clients_healthdiaryquote_updated_by_fkey" FOREIGN KEY ("updated_by") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "clients_healthdiaryquote"
    ADD
        CONSTRAINT "clients_healthdiaryquote_program_id_fkey" FOREIGN KEY ("program_id") REFERENCES "common_program" ("id");

ALTER TABLE
    IF EXISTS "clients_healthdiaryquote"
    ADD
        CONSTRAINT "clients_healthdiaryquote_organisation_id_fkey" FOREIGN KEY ("organisation_id") REFERENCES "common_organisation" ("id");

ALTER TABLE
    IF EXISTS "authority_authoritypermission"
    ADD
        CONSTRAINT "authority_authoritypermission_created_by_fkey" FOREIGN KEY ("created_by") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "authority_authoritypermission"
    ADD
        CONSTRAINT "authority_authoritypermission_updated_by_fkey" FOREIGN KEY ("updated_by") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "authority_authoritypermission"
    ADD
        CONSTRAINT "authority_authoritypermission_organisation_id_fkey" FOREIGN KEY ("organisation_id") REFERENCES "common_organisation" ("id");

ALTER TABLE
    IF EXISTS "authority_authorityrole"
    ADD
        CONSTRAINT "authority_authorityrole_created_by_fkey" FOREIGN KEY ("created_by") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "authority_authorityrole"
    ADD
        CONSTRAINT "authority_authorityrole_updated_by_fkey" FOREIGN KEY ("updated_by") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "authority_authorityrole"
    ADD
        CONSTRAINT "authority_authorityrole_organisation_id_fkey" FOREIGN KEY ("organisation_id") REFERENCES "common_organisation" ("id");

ALTER TABLE
    IF EXISTS "authority_authorityrole_permissions"
    ADD
        CONSTRAINT "authority_authorityrole_permissions_authorityrole_id_fkey" FOREIGN KEY ("authorityrole_id") REFERENCES "authority_authorityrole" ("id");

ALTER TABLE
    IF EXISTS "authority_authorityrole_permissions"
    ADD
        CONSTRAINT "authority_authorityrole_permissions_authoritypermission_id_fkey" FOREIGN KEY ("authoritypermission_id") REFERENCES "authority_authoritypermission" ("id");

ALTER TABLE
    IF EXISTS "communities_community"
    ADD
        CONSTRAINT "communities_community_created_by_fkey" FOREIGN KEY ("created_by") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "communities_community"
    ADD
        CONSTRAINT "communities_community_updated_by_fkey" FOREIGN KEY ("updated_by") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "communities_community"
    ADD
        CONSTRAINT "communities_community_program_id_fkey" FOREIGN KEY ("program_id") REFERENCES "common_program" ("id");

ALTER TABLE
    IF EXISTS "communities_community"
    ADD
        CONSTRAINT "communities_community_organisation_id_fkey" FOREIGN KEY ("organisation_id") REFERENCES "common_organisation" ("id");

ALTER TABLE
    IF EXISTS "screeningtools_screeningtoolsquestion"
    ADD
        CONSTRAINT "screeningtools_screeningtoolsquestion_created_by_fkey" FOREIGN KEY ("created_by") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "screeningtools_screeningtoolsquestion"
    ADD
        CONSTRAINT "screeningtools_screeningtoolsquestion_updated_by_fkey" FOREIGN KEY ("updated_by") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "screeningtools_screeningtoolsquestion"
    ADD
        CONSTRAINT "screeningtools_screeningtoolsquestion_organisation_id_fkey" FOREIGN KEY ("organisation_id") REFERENCES "common_organisation" ("id");

ALTER TABLE
    IF EXISTS "screeningtools_screeningtoolsquestion"
    ADD
        CONSTRAINT "screeningtools_screeningtoolsquestion_program_id_fkey" FOREIGN KEY ("program_id") REFERENCES "common_program" ("id");

ALTER TABLE
    IF EXISTS "users_termsofservice"
    ADD
        CONSTRAINT "users_termsofservice_created_by_fkey" FOREIGN KEY ("created_by") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "users_termsofservice"
    ADD
        CONSTRAINT "users_termsofservice_updated_by_fkey" FOREIGN KEY ("updated_by") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "staff_staff"
    ADD
        CONSTRAINT "staff_staff_created_by_fkey" FOREIGN KEY ("created_by") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "staff_staff"
    ADD
        CONSTRAINT "staff_staff_updated_by_fkey" FOREIGN KEY ("updated_by") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "staff_staff"
    ADD
        CONSTRAINT "staff_staff_current_facility_id_fkey" FOREIGN KEY ("current_facility_id") REFERENCES "common_facility" ("id");

ALTER TABLE
    IF EXISTS "staff_staff"
    ADD
        CONSTRAINT "staff_staff_organisation_id_fkey" FOREIGN KEY ("organisation_id") REFERENCES "common_organisation" ("id");

ALTER TABLE
    IF EXISTS "staff_staff"
    ADD
        CONSTRAINT "staff_staff_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "staff_staff"
    ADD
        CONSTRAINT "staff_staff_program_id_fkey" FOREIGN KEY ("program_id") REFERENCES "common_program" ("id");

ALTER TABLE
    IF EXISTS "staff_staff_facilities"
    ADD
        CONSTRAINT "staff_staff_facilities_staff_id_fkey" FOREIGN KEY ("staff_id") REFERENCES "staff_staff" ("id");

ALTER TABLE
    IF EXISTS "staff_staff_facilities"
    ADD
        CONSTRAINT "staff_staff_facilities_facility_id_fkey" FOREIGN KEY ("facility_id") REFERENCES "common_facility" ("id");

ALTER TABLE
    IF EXISTS "clients_client"
    ADD
        CONSTRAINT "clients_client_created_by_fkey" FOREIGN KEY ("created_by") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "clients_client"
    ADD
        CONSTRAINT "clients_client_updated_by_fkey" FOREIGN KEY ("updated_by") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "clients_client"
    ADD
        CONSTRAINT "clients_client_current_facility_id_fkey" FOREIGN KEY ("current_facility_id") REFERENCES "common_facility" ("id");

ALTER TABLE
    IF EXISTS "clients_client"
    ADD
        CONSTRAINT "clients_client_organisation_id_fkey" FOREIGN KEY ("organisation_id") REFERENCES "common_organisation" ("id");

ALTER TABLE
    IF EXISTS "clients_client"
    ADD
        CONSTRAINT "clients_client_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "clients_client"
    ADD
        CONSTRAINT "clients_client_program_id_fkey" FOREIGN KEY ("program_id") REFERENCES "common_program" ("id");

ALTER TABLE
    IF EXISTS "clients_relatedperson_contacts"
    ADD
        CONSTRAINT "clients_relatedperson_contacts_relatedperson_id_fkey" FOREIGN KEY ("relatedperson_id") REFERENCES "clients_relatedperson" ("id");

ALTER TABLE
    IF EXISTS "clients_relatedperson_contacts"
    ADD
        CONSTRAINT "clients_relatedperson_contacts_contact_id_fkey" FOREIGN KEY ("contact_id") REFERENCES "common_contact" ("id");

ALTER TABLE
    IF EXISTS "clients_clientfacility"
    ADD
        CONSTRAINT "clients_clientfacility_created_by_fkey" FOREIGN KEY ("created_by") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "clients_clientfacility"
    ADD
        CONSTRAINT "clients_clientfacility_updated_by_fkey" FOREIGN KEY ("updated_by") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "clients_clientfacility"
    ADD
        CONSTRAINT "clients_clientfacility_client_id_fkey" FOREIGN KEY ("client_id") REFERENCES "clients_client" ("id");

ALTER TABLE
    IF EXISTS "clients_clientfacility"
    ADD
        CONSTRAINT "clients_clientfacility_facility_id_fkey" FOREIGN KEY ("facility_id") REFERENCES "common_facility" ("id");

ALTER TABLE
    IF EXISTS "clients_clientfacility"
    ADD
        CONSTRAINT "clients_clientfacility_organisation_id_fkey" FOREIGN KEY ("organisation_id") REFERENCES "common_organisation" ("id");

ALTER TABLE
    IF EXISTS "clients_clientfacility"
    ADD
        CONSTRAINT "clients_clientfacility_program_id_fkey" FOREIGN KEY ("program_id") REFERENCES "common_program" ("id");

ALTER TABLE
    IF EXISTS "clients_client_identifiers"
    ADD
        CONSTRAINT "clients_client_identifiers_client_id_fkey" FOREIGN KEY ("client_id") REFERENCES "clients_client" ("id");

ALTER TABLE
    IF EXISTS "clients_client_identifiers"
    ADD
        CONSTRAINT "clients_client_identifiers_identifier_id_fkey" FOREIGN KEY ("identifier_id") REFERENCES "clients_identifier" ("id");

ALTER TABLE
    IF EXISTS "clients_client_related_persons"
    ADD
        CONSTRAINT "clients_client_related_persons_client_id_fkey" FOREIGN KEY ("client_id") REFERENCES "clients_client" ("id");

ALTER TABLE
    IF EXISTS "clients_client_related_persons"
    ADD
        CONSTRAINT "clients_client_related_persons_relatedperson_id_fkey" FOREIGN KEY ("relatedperson_id") REFERENCES "clients_relatedperson" ("id");

ALTER TABLE
    IF EXISTS "common_securityquestionresponse"
    ADD
        CONSTRAINT "common_securityquestionresponse_created_by_fkey" FOREIGN KEY ("created_by") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "common_securityquestionresponse"
    ADD
        CONSTRAINT "common_securityquestionresponse_updated_by_fkey" FOREIGN KEY ("updated_by") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "common_securityquestionresponse"
    ADD
        CONSTRAINT "common_securityquestionresponse_organisation_id_fkey" FOREIGN KEY ("organisation_id") REFERENCES "common_organisation" ("id");

ALTER TABLE
    IF EXISTS "common_securityquestionresponse"
    ADD
        CONSTRAINT "common_securityquestionresponse_question_id_fkey" FOREIGN KEY ("question_id") REFERENCES "common_securityquestion" ("id");

ALTER TABLE
    IF EXISTS "common_securityquestionresponse"
    ADD
        CONSTRAINT "common_securityquestionresponse_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "clients_healthdiaryentry"
    ADD
        CONSTRAINT "clients_healthdiaryentry_created_by_fkey" FOREIGN KEY ("created_by") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "clients_healthdiaryentry"
    ADD
        CONSTRAINT "clients_healthdiaryentry_updated_by_fkey" FOREIGN KEY ("updated_by") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "clients_healthdiaryentry"
    ADD
        CONSTRAINT "clients_healthdiaryentry_client_id_fkey" FOREIGN KEY ("client_id") REFERENCES "clients_client" ("id");

ALTER TABLE
    IF EXISTS "clients_healthdiaryentry"
    ADD
        CONSTRAINT "clients_healthdiaryentry_organisation_id_fkey" FOREIGN KEY ("organisation_id") REFERENCES "common_organisation" ("id");

ALTER TABLE
    IF EXISTS "clients_healthdiaryentry"
    ADD
        CONSTRAINT "clients_healthdiaryentry_program_id_fkey" FOREIGN KEY ("program_id") REFERENCES "common_program" ("id");

ALTER TABLE
    IF EXISTS "clients_servicerequest"
    ADD
        CONSTRAINT "clients_servicerequest_created_by_fkey" FOREIGN KEY ("created_by") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "clients_servicerequest"
    ADD
        CONSTRAINT "clients_servicerequest_updated_by_fkey" FOREIGN KEY ("updated_by") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "clients_servicerequest"
    ADD
        CONSTRAINT "clients_servicerequest_in_progress_by_id_fkey" FOREIGN KEY ("in_progress_by_id") REFERENCES "staff_staff" ("id");

ALTER TABLE
    IF EXISTS "clients_servicerequest"
    ADD
        CONSTRAINT "clients_servicerequest_resolved_by_id_fkey" FOREIGN KEY ("resolved_by_id") REFERENCES "staff_staff" ("id");

ALTER TABLE
    IF EXISTS "clients_servicerequest"
    ADD
        CONSTRAINT "clients_servicerequest_facility_id_fkey" FOREIGN KEY ("facility_id") REFERENCES "common_facility" ("id");

ALTER TABLE
    IF EXISTS "clients_servicerequest"
    ADD
        CONSTRAINT "clients_servicerequest_client_id_fkey" FOREIGN KEY ("client_id") REFERENCES "clients_client" ("id");

ALTER TABLE
    IF EXISTS "clients_servicerequest"
    ADD
        CONSTRAINT "clients_servicerequest_organisation_id_fkey" FOREIGN KEY ("organisation_id") REFERENCES "common_organisation" ("id");

ALTER TABLE
    IF EXISTS "clients_servicerequest"
    ADD
        CONSTRAINT "clients_servicerequest_program_id_fkey" FOREIGN KEY ("program_id") REFERENCES "common_program" ("id");

ALTER TABLE
    IF EXISTS "staff_staff_identifiers"
    ADD
        CONSTRAINT "staff_staff_identifiers_staff_id_fkey" FOREIGN KEY ("staff_id") REFERENCES "staff_staff" ("id");

ALTER TABLE
    IF EXISTS "staff_staff_identifiers"
    ADD
        CONSTRAINT "staff_staff_identifiers_identifier_id_fkey" FOREIGN KEY ("identifier_id") REFERENCES "clients_identifier" ("id");

ALTER TABLE
    IF EXISTS "appointments_appointment"
    ADD
        CONSTRAINT "appointments_appointment_created_by_fkey" FOREIGN KEY ("created_by") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "appointments_appointment"
    ADD
        CONSTRAINT "appointments_appointment_updated_by_fkey" FOREIGN KEY ("updated_by") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "appointments_appointment"
    ADD
        CONSTRAINT "appointments_appointment_client_id_fkey" FOREIGN KEY ("client_id") REFERENCES "clients_client" ("id");

ALTER TABLE
    IF EXISTS "appointments_appointment"
    ADD
        CONSTRAINT "appointments_appointment_organisation_id_fkey" FOREIGN KEY ("organisation_id") REFERENCES "common_organisation" ("id");

ALTER TABLE
    IF EXISTS "appointments_appointment"
    ADD
        CONSTRAINT "appointments_appointment_staff_id_fkey" FOREIGN KEY ("staff_id") REFERENCES "staff_staff" ("id");

ALTER TABLE
    IF EXISTS "appointments_appointment"
    ADD
        CONSTRAINT "appointments_appointment_facility_id_fkey" FOREIGN KEY ("facility_id") REFERENCES "common_facility" ("id");

ALTER TABLE
    IF EXISTS "appointments_appointment"
    ADD
        CONSTRAINT "appointments_appointment_program_id_fkey" FOREIGN KEY ("program_id") REFERENCES "common_program" ("id");

ALTER TABLE
    IF EXISTS "authority_authorityrole_users"
    ADD
        CONSTRAINT "authority_authorityrole_users_authorityrole_id_fkey" FOREIGN KEY ("authorityrole_id") REFERENCES "authority_authorityrole" ("id");

ALTER TABLE
    IF EXISTS "authority_authorityrole_users"
    ADD
        CONSTRAINT "authority_authorityrole_users_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "common_notification"
    ADD
        CONSTRAINT "common_notification_created_by_fkey" FOREIGN KEY ("created_by") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "common_notification"
    ADD
        CONSTRAINT "common_notification_updated_by_fkey" FOREIGN KEY ("updated_by") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "common_notification"
    ADD
        CONSTRAINT "common_notification_facility_id_fkey" FOREIGN KEY ("facility_id") REFERENCES "common_facility" ("id");

ALTER TABLE
    IF EXISTS "common_notification"
    ADD
        CONSTRAINT "common_notification_organisation_id_fkey" FOREIGN KEY ("organisation_id") REFERENCES "common_organisation" ("id");

ALTER TABLE
    IF EXISTS "common_notification"
    ADD
        CONSTRAINT "common_notification_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "common_notification"
    ADD
        CONSTRAINT "common_notification_program_id_fkey" FOREIGN KEY ("program_id") REFERENCES "common_program" ("id");

ALTER TABLE
    IF EXISTS "common_usersurveys"
    ADD
        CONSTRAINT "common_usersurveys_created_by_fkey" FOREIGN KEY ("created_by") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "common_usersurveys"
    ADD
        CONSTRAINT "common_usersurveys_updated_by_fkey" FOREIGN KEY ("updated_by") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "common_usersurveys"
    ADD
        CONSTRAINT "common_usersurveys_organisation_id_fkey" FOREIGN KEY ("organisation_id") REFERENCES "common_organisation" ("id");

ALTER TABLE
    IF EXISTS "common_usersurveys"
    ADD
        CONSTRAINT "common_usersurveys_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "common_usersurveys"
    ADD
        CONSTRAINT "common_usersurveys_program_id_fkey" FOREIGN KEY ("program_id") REFERENCES "common_program" ("id");

ALTER TABLE
    IF EXISTS "common_feedback"
    ADD
        CONSTRAINT "common_feedback_created_by_fkey" FOREIGN KEY ("created_by") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "common_feedback"
    ADD
        CONSTRAINT "common_feedback_updated_by_fkey" FOREIGN KEY ("updated_by") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "common_feedback"
    ADD
        CONSTRAINT "common_feedback_organisation_id_fkey" FOREIGN KEY ("organisation_id") REFERENCES "common_organisation" ("id");

ALTER TABLE
    IF EXISTS "common_feedback"
    ADD
        CONSTRAINT "common_feedback_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "common_feedback"
    ADD
        CONSTRAINT "common_feedback_program_id_fkey" FOREIGN KEY ("program_id") REFERENCES "common_program" ("id");

ALTER TABLE
    IF EXISTS "communities_community_clients"
    ADD
        CONSTRAINT "communities_community_clients_community_id_fkey" FOREIGN KEY ("community_id") REFERENCES "communities_community" ("id");

ALTER TABLE
    IF EXISTS "communities_community_clients"
    ADD
        CONSTRAINT "communities_community_clients_client_id_fkey" FOREIGN KEY ("client_id") REFERENCES "clients_client" ("id");

ALTER TABLE
    IF EXISTS "communities_community_staff"
    ADD
        CONSTRAINT "communities_community_staff_community_id_fkey" FOREIGN KEY ("community_id") REFERENCES "communities_community" ("id");

ALTER TABLE
    IF EXISTS "communities_community_staff"
    ADD
        CONSTRAINT "communities_community_staff_staff_id_fkey" FOREIGN KEY ("staff_id") REFERENCES "staff_staff" ("id");

ALTER TABLE
    IF EXISTS "screeningtools_screeningtoolsresponse"
    ADD
        CONSTRAINT "screeningtools_screeningtoolsresponse_created_by_fkey" FOREIGN KEY ("created_by") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "screeningtools_screeningtoolsresponse"
    ADD
        CONSTRAINT "screeningtools_screeningtoolsresponse_updated_by_fkey" FOREIGN KEY ("updated_by") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "screeningtools_screeningtoolsresponse"
    ADD
        CONSTRAINT "screeningtools_screeningtoolsresponse_client_id_fkey" FOREIGN KEY ("client_id") REFERENCES "clients_client" ("id");

ALTER TABLE
    IF EXISTS "screeningtools_screeningtoolsresponse"
    ADD
        CONSTRAINT "screeningtools_screeningtoolsresponse_organisation_id_fkey" FOREIGN KEY ("organisation_id") REFERENCES "common_organisation" ("id");

ALTER TABLE
    IF EXISTS "screeningtools_screeningtoolsresponse"
    ADD
        CONSTRAINT "screeningtools_screeningtoolsresponse_question_id_fkey" FOREIGN KEY ("question_id") REFERENCES "screeningtools_screeningtoolsquestion" ("id");

ALTER TABLE
    IF EXISTS "screeningtools_screeningtoolsresponse"
    ADD
        CONSTRAINT "screeningtools_screeningtoolsresponse_program_id_fkey" FOREIGN KEY ("program_id") REFERENCES "common_program" ("id");

ALTER TABLE
    IF EXISTS "staff_servicerequest"
    ADD
        CONSTRAINT "staff_servicerequest_created_by_fkey" FOREIGN KEY ("created_by") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "staff_servicerequest"
    ADD
        CONSTRAINT "staff_servicerequest_updated_by_fkey" FOREIGN KEY ("updated_by") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "staff_servicerequest"
    ADD
        CONSTRAINT "staff_servicerequest_organisation_id_fkey" FOREIGN KEY ("organisation_id") REFERENCES "common_organisation" ("id");

ALTER TABLE
    IF EXISTS "staff_servicerequest"
    ADD
        CONSTRAINT "staff_servicerequest_resolved_by_id_fkey" FOREIGN KEY ("resolved_by_id") REFERENCES "staff_staff" ("id");

ALTER TABLE
    IF EXISTS "staff_servicerequest"
    ADD
        CONSTRAINT "staff_servicerequest_staff_id_fkey" FOREIGN KEY ("staff_id") REFERENCES "staff_staff" ("id");

ALTER TABLE
    IF EXISTS "staff_servicerequest"
    ADD
        CONSTRAINT "staff_servicerequest_facility_id_fkey" FOREIGN KEY ("facility_id") REFERENCES "common_facility" ("id");

ALTER TABLE
    IF EXISTS "staff_servicerequest"
    ADD
        CONSTRAINT "staff_servicerequest_program_id_fkey" FOREIGN KEY ("program_id") REFERENCES "common_program" ("id");

ALTER TABLE
    IF EXISTS "users_metric"
    ADD
        CONSTRAINT "users_metric_created_by_fkey" FOREIGN KEY ("created_by") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "users_metric"
    ADD
        CONSTRAINT "users_metric_updated_by_fkey" FOREIGN KEY ("updated_by") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "users_metric"
    ADD
        CONSTRAINT "users_metric_organisation_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "users_userpin"
    ADD
        CONSTRAINT "users_userpin_created_by_fkey" FOREIGN KEY ("created_by") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "users_userpin"
    ADD
        CONSTRAINT "users_userpin_updated_by_fkey" FOREIGN KEY ("updated_by") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "users_userpin"
    ADD
        CONSTRAINT "users_userpin_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "users_userpin"
    ADD
        CONSTRAINT "users_userpin_organisation_id_fkey" FOREIGN KEY ("organisation_id") REFERENCES "common_organisation" ("id");

ALTER TABLE
    IF EXISTS "users_userotp"
    ADD
        CONSTRAINT "users_userotp_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "users_userotp"
    ADD
        CONSTRAINT "users_userotp_created_by_fkey" FOREIGN KEY ("created_by") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "users_userotp"
    ADD
        CONSTRAINT "users_userotp_updated_by_fkey" FOREIGN KEY ("updated_by") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "questionnaires_questionnaire"
    ADD
        CONSTRAINT "questionnaires_questionnaire_created_by_fkey" FOREIGN KEY ("created_by") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "questionnaires_questionnaire"
    ADD
        CONSTRAINT "questionnaires_questionnaire_updated_by_fkey" FOREIGN KEY ("updated_by") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "questionnaires_questionnaire"
    ADD
        CONSTRAINT "questionnaires_questionnaire_organisation_id_fkey" FOREIGN KEY ("organisation_id") REFERENCES "common_organisation" ("id");

ALTER TABLE
    IF EXISTS "questionnaires_questionnaire"
    ADD
        CONSTRAINT "questionnaires_questionnaire_program_id_fkey" FOREIGN KEY ("program_id") REFERENCES "common_program" ("id");

ALTER TABLE
    IF EXISTS "questionnaires_screeningtool"
    ADD
        CONSTRAINT "questionnaires_screeningtool_created_by_fkey" FOREIGN KEY ("created_by") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "questionnaires_screeningtool"
    ADD
        CONSTRAINT "questionnaires_screeningtool_updated_by_fkey" FOREIGN KEY ("updated_by") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "questionnaires_screeningtool"
    ADD
        CONSTRAINT "questionnaires_screeningtool_organisation_id_fkey" FOREIGN KEY ("organisation_id") REFERENCES "common_organisation" ("id");

ALTER TABLE
    IF EXISTS "questionnaires_screeningtool"
    ADD
        CONSTRAINT "questionnaires_screeningtool_questionnaire_id_fkey" FOREIGN KEY ("questionnaire_id") REFERENCES "questionnaires_questionnaire" ("id");

ALTER TABLE
    IF EXISTS "questionnaires_screeningtool"
    ADD
        CONSTRAINT "questionnaires_screeningtool_program_id_fkey" FOREIGN KEY ("program_id") REFERENCES "common_program" ("id");

ALTER TABLE
    IF EXISTS "questionnaires_question"
    ADD
        CONSTRAINT "questionnaires_question_created_by_fkey" FOREIGN KEY ("created_by") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "questionnaires_question"
    ADD
        CONSTRAINT "questionnaires_question_updated_by_fkey" FOREIGN KEY ("updated_by") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "questionnaires_question"
    ADD
        CONSTRAINT "questionnaires_question_organisation_id_fkey" FOREIGN KEY ("organisation_id") REFERENCES "common_organisation" ("id");

ALTER TABLE
    IF EXISTS "questionnaires_question"
    ADD
        CONSTRAINT "questionnaires_question_questionnaire_id_fkey" FOREIGN KEY ("questionnaire_id") REFERENCES "questionnaires_questionnaire" ("id");

ALTER TABLE
    IF EXISTS "questionnaires_question"
    ADD
        CONSTRAINT "questionnaires_question_program_id_fkey" FOREIGN KEY ("program_id") REFERENCES "common_program" ("id");

ALTER TABLE
    IF EXISTS "questionnaires_questioninputchoice"
    ADD
        CONSTRAINT "questionnaires_questioninputchoice_created_by_fkey" FOREIGN KEY ("created_by") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "questionnaires_questioninputchoice"
    ADD
        CONSTRAINT "questionnaires_questioninputchoice_updated_by_fkey" FOREIGN KEY ("updated_by") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "questionnaires_questioninputchoice"
    ADD
        CONSTRAINT "questionnaires_questioninputchoice_organisation_id_fkey" FOREIGN KEY ("organisation_id") REFERENCES "common_organisation" ("id");

ALTER TABLE
    IF EXISTS "questionnaires_questioninputchoice"
    ADD
        CONSTRAINT "questionnaires_questioninputchoice_question_id_fkey" FOREIGN KEY ("question_id") REFERENCES "questionnaires_question" ("id");

ALTER TABLE
    IF EXISTS "questionnaires_questioninputchoice"
    ADD
        CONSTRAINT "questionnaires_questioninputchoice_program_id_fkey" FOREIGN KEY ("program_id") REFERENCES "common_program" ("id");

ALTER TABLE
    IF EXISTS "questionnaires_screeningtoolresponse"
    ADD
        CONSTRAINT "questionnaires_screeningtoolresponse_created_by_fkey" FOREIGN KEY ("created_by") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "questionnaires_screeningtoolresponse"
    ADD
        CONSTRAINT "questionnaires_screeningtoolresponse_updated_by_fkey" FOREIGN KEY ("updated_by") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "questionnaires_screeningtoolresponse"
    ADD
        CONSTRAINT "questionnaires_screeningtoolresponse_organisation_id_fkey" FOREIGN KEY ("organisation_id") REFERENCES "common_organisation" ("id");

ALTER TABLE
    IF EXISTS "questionnaires_screeningtoolresponse"
    ADD
        CONSTRAINT "questionnaires_screeningtoolresponse_screeningtool_id_fkey" FOREIGN KEY ("screeningtool_id") REFERENCES "questionnaires_screeningtool" ("id");

ALTER TABLE
    IF EXISTS "questionnaires_screeningtoolresponse"
    ADD
        CONSTRAINT "questionnaires_screeningtoolresponse_facility_id_fkey" FOREIGN KEY ("facility_id") REFERENCES "common_facility" ("id");

ALTER TABLE
    IF EXISTS "questionnaires_screeningtoolresponse"
    ADD
        CONSTRAINT "questionnaires_screeningtoolresponse_client_id_fkey" FOREIGN KEY ("client_id") REFERENCES "clients_client" ("id");

ALTER TABLE
    IF EXISTS "questionnaires_screeningtoolresponse"
    ADD
        CONSTRAINT "questionnaires_screeningtoolresponse_program_id_fkey" FOREIGN KEY ("program_id") REFERENCES "common_program" ("id");

ALTER TABLE
    IF EXISTS "questionnaires_screeningtoolquestionresponse"
    ADD
        CONSTRAINT "questionnaires_screeningtoolquestionresponse_created_by_fkey" FOREIGN KEY ("created_by") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "questionnaires_screeningtoolquestionresponse"
    ADD
        CONSTRAINT "questionnaires_screeningtoolquestionresponse_updated_by_fkey" FOREIGN KEY ("updated_by") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "questionnaires_screeningtoolquestionresponse"
    ADD
        CONSTRAINT "questionnes_screeningtoolquestionresponse_organisation_id_fkey" FOREIGN KEY ("organisation_id") REFERENCES "common_organisation" ("id");

ALTER TABLE
    IF EXISTS "questionnaires_screeningtoolquestionresponse"
    ADD
        CONSTRAINT "questionnaires_screeniuestionresponse_screeniesponse_id_fkey" FOREIGN KEY ("screeningtoolresponse_id") REFERENCES "questionnaires_screeningtoolresponse" ("id");

ALTER TABLE
    IF EXISTS "questionnaires_screeningtoolquestionresponse"
    ADD
        CONSTRAINT "questionnaires_screeningtoolquestionresponse_question_id_fkey" FOREIGN KEY ("question_id") REFERENCES "questionnaires_question" ("id");

ALTER TABLE
    IF EXISTS "questionnaires_screeningtoolquestionresponse"
    ADD
        CONSTRAINT "questionnaires_screeningtoolquestionresponse_program_id_fkey" FOREIGN KEY ("program_id") REFERENCES "common_program" ("id");

ALTER TABLE
    IF EXISTS "questionnaires_screeningtoolquestionresponse"
    ADD
        CONSTRAINT "questionnaires_screeningtoolquestionresponse_facility_id_fkey" FOREIGN KEY ("facility_id") REFERENCES "common_facility" ("id");

ALTER TABLE
    IF EXISTS "clients_client_facilities"
    ADD
        CONSTRAINT "clients_client_facilities_client_id_fkey" FOREIGN KEY ("client_id") REFERENCES "clients_client" ("id");

ALTER TABLE
    IF EXISTS "clients_client_facilities"
    ADD
        CONSTRAINT "clients_client_facilities_facility_id_fkey" FOREIGN KEY ("facility_id") REFERENCES "common_facility" ("id");

ALTER TABLE
    IF EXISTS "caregivers_caregiver"
    ADD
        CONSTRAINT "caregivers_caregiver_created_by_fkey" FOREIGN KEY ("created_by") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "caregivers_caregiver"
    ADD
        CONSTRAINT "caregivers_caregiver_updated_by_fkey" FOREIGN KEY ("updated_by") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "caregivers_caregiver"
    ADD
        CONSTRAINT "caregivers_caregiver_organisation_id_fkey" FOREIGN KEY ("organisation_id") REFERENCES "common_organisation" ("id");

ALTER TABLE
    IF EXISTS "caregivers_caregiver"
    ADD
        CONSTRAINT "caregivers_caregiver_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "caregivers_caregiver"
    ADD
        CONSTRAINT "caregivers_caregiver_program_id_fkey" FOREIGN KEY ("program_id") REFERENCES "common_program" ("id");

ALTER TABLE
    IF EXISTS "caregivers_caregiver_client"
    ADD
        CONSTRAINT "caregivers_caregiver_client_caregiver_id_fkey" FOREIGN KEY ("caregiver_id") REFERENCES "caregivers_caregiver" ("id");

ALTER TABLE
    IF EXISTS "caregivers_caregiver_client"
    ADD
        CONSTRAINT "caregivers_caregiver_client_client_id_fkey" FOREIGN KEY ("client_id") REFERENCES "clients_client" ("id");

ALTER TABLE
    IF EXISTS "caregivers_caregiver_client"
    ADD
        CONSTRAINT "caregivers_caregiver_client_created_by_fkey" FOREIGN KEY ("created_by") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "caregivers_caregiver_client"
    ADD
        CONSTRAINT "caregivers_caregiver_client_updated_by_fkey" FOREIGN KEY ("updated_by") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "caregivers_caregiver_client"
    ADD
        CONSTRAINT "caregivers_caregiver_client_organisation_id_fkey" FOREIGN KEY ("organisation_id") REFERENCES "common_organisation" ("id");

ALTER TABLE
    IF EXISTS "caregivers_caregiver_client"
    ADD
        CONSTRAINT "caregivers_caregiver_client_assigned_by_fkey" FOREIGN KEY ("assigned_by") REFERENCES "staff_staff" ("id");

ALTER TABLE
    IF EXISTS "caregivers_caregiver_client"
    ADD
        CONSTRAINT "caregivers_caregiver_client_program_id_fkey" FOREIGN KEY ("program_id") REFERENCES "common_program" ("id");

COMMIT;