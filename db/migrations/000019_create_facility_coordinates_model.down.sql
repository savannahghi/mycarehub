BEGIN;

ALTER TABLE
    IF EXISTS "facility_coordinates"
    DROP CONSTRAINT IF EXISTS "facility_coordinates_facility_id_fkey";

DROP TABLE IF EXISTS "facility_coordinates";

END;