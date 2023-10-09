BEGIN;

CREATE TABLE IF NOT EXISTS "facility_coordinates" (
    "id" uuid PRIMARY KEY NOT NULL,
    "created_by" uuid,
    "created" timestamp NOT NULL,
    "updated_by" uuid,
    "updated" timestamp NOT NULL,
    "deleted_at" timestamp,
    "active" boolean NOT NULL,
    "lat" float NOT NULL,
    "lng" float NOT NULL,
    "facility_id" uuid NOT NULL
);

ALTER TABLE
    IF EXISTS "facility_coordinates"
ADD
    CONSTRAINT "facility_coordinates_facility_id_fkey" FOREIGN KEY ("facility_id") REFERENCES "common_facility" ("id");

END;