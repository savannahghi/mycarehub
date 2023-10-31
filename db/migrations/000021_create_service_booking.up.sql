BEGIN;

CREATE TABLE IF NOT EXISTS "service_booking" (
  "id" uuid PRIMARY KEY NOT NULL,
  "active" boolean NOT NULL,
  "created" timestamp NOT NULL,
  "created_by" uuid,
  "updated" timestamp NOT NULL,
  "updated_by" uuid,
  "deleted_at" timestamp,
  "services" varchar(64)[],
  "date" timestamp,
  "verification_code" varchar(64) NOT NULL,
  "client_id" uuid NOT NULL,
  "organisation_id" uuid NOT NULL,
  "facility_id" uuid,
  "program_id" uuid NOT NULL
);

ALTER TABLE
    IF EXISTS "service_booking"
    ADD
        CONSTRAINT "service_booking_created_by_fkey" FOREIGN KEY ("created_by") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "service_booking"
    ADD
        CONSTRAINT "service_booking_updated_by_fkey" FOREIGN KEY ("updated_by") REFERENCES "users_user" ("id");

ALTER TABLE
    IF EXISTS "service_booking"
    ADD
        CONSTRAINT "service_booking_client_id_fkey" FOREIGN KEY ("client_id") REFERENCES "clients_client" ("id");

ALTER TABLE
    IF EXISTS "service_booking"
    ADD
        CONSTRAINT "service_booking_organisation_id_fkey" FOREIGN KEY ("organisation_id") REFERENCES "common_organisation" ("id");

ALTER TABLE
    IF EXISTS "service_booking"
    ADD
        CONSTRAINT "service_booking_program_id_fkey" FOREIGN KEY ("program_id") REFERENCES "common_program" ("id");

COMMIT;