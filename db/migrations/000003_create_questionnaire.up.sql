BEGIN;

-- Questionnaires 
CREATE TABLE IF NOT EXISTS questionnaires_questionnaire (
    id uuid NOT NULL PRIMARY KEY,
    active boolean NOT NULL,
    created timestamp WITH time zone NOT NULL,
    created_by uuid,
    updated timestamp WITH time zone NOT NULL,
    updated_by uuid,
    deleted_at timestamp WITH time zone,
    organisation_id uuid NOT NULL CONSTRAINT questionnaires_questionnaire_organisation_id REFERENCES common_organisation deferrable initially deferred,
    name varchar(64) NOT NULL UNIQUE,
    description text NOT NULL
);
CREATE INDEX IF NOT EXISTS questionnaires_questionnaire_name ON questionnaires_questionnaire (name varchar_pattern_ops);

CREATE TABLE IF NOT EXISTS questionnaires_screeningtool (
    id uuid NOT NULL PRIMARY KEY,
    active boolean NOT NULL,
    created timestamp WITH time zone NOT NULL,
    created_by uuid,
    updated timestamp WITH time zone NOT NULL,
    updated_by uuid,
    deleted_at timestamp WITH time zone,
    organisation_id uuid NOT NULL CONSTRAINT questionnaires_screeningtool_organisation_id REFERENCES common_organisation deferrable initially deferred,
    questionnaire_id uuid NOT NULL CONSTRAINT questionnaires_screeningtool_questionnaire_id REFERENCES questionnaires_questionnaire deferrable initially deferred,
    threshold integer DEFAULT 0,
    client_types varchar(64) [] NULL,
    genders varchar(64) [] NULL,
    min_age integer NULL,
    max_age integer NULL
);
CREATE INDEX IF NOT EXISTS questionnaires_screeningtool_questionnaire_id ON questionnaires_screeningtool (questionnaire_id);
CREATE INDEX IF NOT EXISTS questionnaires_screeningtool_genders ON questionnaires_screeningtool (genders);
CREATE INDEX IF NOT EXISTS questionnaires_screeningtool_client_types ON questionnaires_screeningtool(client_types);
CREATE INDEX IF NOT EXISTS questionnaires_screeningtool_min_age ON questionnaires_screeningtool(min_age);
CREATE INDEX IF NOT EXISTS questionnaires_screeningtool_max_age ON questionnaires_screeningtool(max_age);

-- Questions
CREATE TABLE IF NOT EXISTS questionnaires_question (
    id uuid NOT NULL PRIMARY KEY,
    active boolean NOT NULL,
    created timestamp WITH time zone NOT NULL,
    created_by uuid,
    updated timestamp WITH time zone NOT NULL,
    updated_by uuid,
    deleted_at timestamp WITH time zone,
    organisation_id uuid NOT NULL CONSTRAINT questionnaires_questionnaire_organisation_id REFERENCES common_organisation deferrable initially deferred,
    questionnaire_id uuid NOT NULL CONSTRAINT questionnaires_questionnaire_questionnaire_id REFERENCES questionnaires_questionnaire deferrable initially deferred,
    text text NOT NULL UNIQUE,
    question_type varchar(64) NOT NULL,
    response_value_type varchar(64) NOT NULL,
    select_multiple boolean NULL,
    required boolean NOT NULL,
    sequence integer NOT NULL
);
CREATE INDEX IF NOT EXISTS questionnaires_question_questionnaire_id ON questionnaires_question (questionnaire_id);

CREATE TABLE IF NOT EXISTS questionnaires_questioninputchoice (
    id uuid NOT NULL PRIMARY KEY,
    active boolean NOT NULL,
    created timestamp WITH time zone NOT NULL,
    created_by uuid,
    updated timestamp WITH time zone NOT NULL,
    updated_by uuid,
    deleted_at timestamp WITH time zone,
    organisation_id uuid NOT NULL CONSTRAINT questionnaires_questioninputchoice_organisation_id REFERENCES common_organisation deferrable initially deferred,
    question_id uuid NOT NULL CONSTRAINT questionnaires_questioninputchoice_question_id REFERENCES questionnaires_question deferrable initially deferred,
    choice varchar(64) UNIQUE,
    value text NOT NULL,
    score integer DEFAULT 0
);
CREATE INDEX IF NOT EXISTS questionnaires_questioninputchoice_question_id ON questionnaires_questioninputchoice (question_id);
CREATE INDEX IF NOT EXISTS questionnaires_questioninputchoice_choice ON questionnaires_questioninputchoice (choice varchar_pattern_ops);

--  Responses
CREATE TABLE IF NOT EXISTS questionnaires_screeningtoolresponse (
    id uuid NOT NULL PRIMARY KEY,
    active boolean NOT NULL,
    created timestamp WITH time zone NOT NULL,
    created_by uuid,
    updated timestamp WITH time zone NOT NULL,
    updated_by uuid,
    deleted_at timestamp WITH time zone,
    organisation_id uuid NOT NULL CONSTRAINT questionnaires_screeningtoolresponse_organisation_id REFERENCES common_organisation deferrable initially deferred,
    screeningtool_id uuid NOT NULL CONSTRAINT questionnaires_screeningtoolresponse_screeningtool_id REFERENCES questionnaires_screeningtool deferrable initially deferred,
    facility_id uuid NOT NULL CONSTRAINT questionnaires_screeningtoolresponse_facility_id REFERENCES common_facility deferrable initially deferred,
    client_id uuid NOT NULL CONSTRAINT questionnaires_screeningtoolresponse_client_id REFERENCES clients_client deferrable initially deferred,
    aggregate_score integer DEFAULT 0
);
CREATE INDEX IF NOT EXISTS questionnaires_screeningtoolresponse_screeningtool_id ON questionnaires_screeningtoolresponse (screeningtool_id);
CREATE INDEX IF NOT EXISTS questionnaires_screeningtoolresponse_facility_id ON questionnaires_screeningtoolresponse (facility_id);
CREATE INDEX IF NOT EXISTS questionnaires_screeningtoolresponse_client_id ON questionnaires_screeningtoolresponse (client_id);

CREATE TABLE IF NOT EXISTS questionnaires_screeningtoolquestionresponse (
    id uuid NOT NULL PRIMARY KEY,
    active boolean NOT NULL,
    created timestamp WITH time zone NOT NULL,
    created_by uuid,
    updated timestamp WITH time zone NOT NULL,
    updated_by uuid,
    deleted_at timestamp WITH time zone,
    organisation_id uuid NOT NULL CONSTRAINT questionnaires_screeningtoolquestionresponse_organisation_id REFERENCES common_organisation deferrable initially deferred,
    screeningtoolresponse_id uuid NOT NULL CONSTRAINT questionnaires_screeningtoolquestionresponse_screeningtoolresponse_id REFERENCES questionnaires_screeningtoolresponse deferrable initially deferred,
    question_id uuid NOT NULL CONSTRAINT questionnaires_screeningtoolquestionresponse_question_id REFERENCES questionnaires_question deferrable initially deferred,
    response text NOT NULL,
    score int DEFAULT 0
);
CREATE INDEX IF NOT EXISTS questionnaires_screeningtoolquestionresponse_screeningtoolresponse_id ON questionnaires_screeningtoolquestionresponse (screeningtoolresponse_id);
CREATE INDEX IF NOT EXISTS questionnaires_screeningtoolquestionresponse_question_id ON questionnaires_screeningtoolquestionresponse (question_id);
CREATE INDEX IF NOT EXISTS questionnaires_screeningtoolquestionresponse_response ON questionnaires_screeningtoolquestionresponse (response text_pattern_ops);

COMMIT;