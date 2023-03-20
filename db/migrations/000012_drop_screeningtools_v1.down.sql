BEGIN;

CREATE TABLE IF NOT EXISTS "screeningtools_screeningtoolsquestion" (
  "id" uuid PRIMARY KEY NOT NULL,
  "active" boolean NOT NULL,
  "created" timestamp NOT NULL,
  "created_by" uuid,
  "updated" timestamp NOT NULL,
  "updated_by" uuid,
  "deleted_at" timestamp,
  "question" text NOT NULL,
  "tool_type" varchar(32) NOT NULL,
  "response_choices" jsonb,
  "response_type" varchar(32) NOT NULL,
  "response_category" varchar(32) NOT NULL,
  "sequence" integer NOT NULL,
  "meta" jsonb,
  "organisation_id" uuid NOT NULL,
  "program_id" uuid NOT NULL,
  CONSTRAINT "screeningtools_screeningtoolsquestion_created_by_fkey" FOREIGN KEY ("created_by") REFERENCES "users_user" ("id"),
  CONSTRAINT "screeningtools_screeningtoolsquestion_updated_by_fkey" FOREIGN KEY ("updated_by") REFERENCES "users_user" ("id"),
  CONSTRAINT "screeningtools_screeningtoolsquestion_organisation_id_fkey" FOREIGN KEY ("organisation_id") REFERENCES "common_organisation" ("id"),
  CONSTRAINT "screeningtools_screeningtoolsquestion_program_id_fkey" FOREIGN KEY ("program_id") REFERENCES "common_program" ("id")
);


CREATE TABLE IF NOT EXISTS "screeningtools_screeningtoolsresponse" (
  "id" uuid PRIMARY KEY NOT NULL,
  "active" boolean NOT NULL,
  "created" timestamp NOT NULL,
  "created_by" uuid,
  "updated" timestamp NOT NULL,
  "updated_by" uuid,
  "deleted_at" timestamp,
  "response" text NOT NULL,
  "client_id" uuid NOT NULL,
  "organisation_id" uuid NOT NULL,
  "question_id" uuid NOT NULL,
  "program_id" uuid NOT NULL,
  CONSTRAINT "screeningtools_screeningtoolsresponse_created_by_fkey" FOREIGN KEY ("created_by") REFERENCES "users_user" ("id"),
  CONSTRAINT "screeningtools_screeningtoolsresponse_updated_by_fkey" FOREIGN KEY ("updated_by") REFERENCES "users_user" ("id"),
  CONSTRAINT "screeningtools_screeningtoolsresponse_client_id_fkey" FOREIGN KEY ("client_id") REFERENCES "clients_client" ("id"),
  CONSTRAINT "screeningtools_screeningtoolsresponse_organisation_id_fkey" FOREIGN KEY ("organisation_id") REFERENCES "common_organisation" ("id"),
  CONSTRAINT "screeningtools_screeningtoolsresponse_question_id_fkey" FOREIGN KEY ("question_id") REFERENCES "screeningtools_screeningtoolsquestion" ("id"),
  CONSTRAINT "screeningtools_screeningtoolsresponse_program_id_fkey" FOREIGN KEY ("program_id") REFERENCES "common_program" ("id")
);

COMMIT;