BEGIN;

DROP TABLE IF EXISTS users_userotp;

DROP TABLE IF EXISTS users_userpin;

DROP TABLE IF EXISTS users_metric;

DROP TABLE IF EXISTS staff_servicerequest;

DROP TABLE IF EXISTS screeningtools_screeningtoolsresponse;

DROP TABLE IF EXISTS communities_community_staff;

DROP TABLE IF EXISTS communities_community_clients;

DROP TABLE IF EXISTS common_feedback;

DROP TABLE IF EXISTS common_usersurveys;

DROP TABLE IF EXISTS common_notification;

DROP TABLE IF EXISTS authority_authorityrole_users;

DROP TABLE IF EXISTS appointments_appointment;

DROP TABLE IF EXISTS staff_staff_identifiers;

DROP TABLE IF EXISTS staff_staff_contacts;

DROP TABLE IF EXISTS staff_staff_addresses;

DROP TABLE IF EXISTS clients_servicerequest;

DROP TABLE IF EXISTS clients_healthdiaryattachment;

DROP TABLE IF EXISTS clients_healthdiaryentry;

DROP TABLE IF EXISTS clients_securityquestionresponse;

DROP TABLE IF EXISTS clients_client_related_persons;

DROP TABLE IF EXISTS clients_client_identifiers;

DROP TABLE IF EXISTS clients_clientfacility;

DROP TABLE IF EXISTS clients_relatedperson_contacts;

DROP TABLE IF EXISTS clients_client_contacts;

DROP TABLE IF EXISTS clients_client_addresses;

DROP TABLE IF EXISTS clients_client;

DROP TABLE IF EXISTS staff_staff_facilities;

DROP TABLE IF EXISTS staff_staff;

DROP TABLE IF EXISTS common_contact;

DROP TABLE IF EXISTS users_user;

DROP TABLE IF EXISTS users_termsofservice;

DROP TABLE IF EXISTS screeningtools_screeningtoolsquestion;


DROP TABLE IF EXISTS communities_community;

DROP TABLE IF EXISTS authority_authorityrole_permissions;

DROP TABLE IF EXISTS authority_authorityrole;

DROP TABLE IF EXISTS authority_authoritypermission;

DROP TABLE IF EXISTS clients_caregiver;

DROP TABLE IF EXISTS clients_healthdiaryquote;

DROP TABLE IF EXISTS clients_identifier;

DROP TABLE IF EXISTS clients_relatedperson_addresses;

DROP TABLE IF EXISTS clients_relatedperson;

DROP TABLE IF EXISTS clients_securityquestion;

DROP TABLE IF EXISTS common_address;

DROP TABLE IF EXISTS common_auditlog;

DROP TABLE IF EXISTS common_facilityattachment;

DROP TABLE IF EXISTS common_facility;

DROP TABLE IF EXISTS common_organisation;

COMMIT;