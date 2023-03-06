BEGIN;

ALTER TABLE
    IF EXISTS "authority_authorityrole_users"
    DROP CONSTRAINT IF EXISTS "authority_authorityrole_users_authorityrole_id_fkey",
    DROP CONSTRAINT IF EXISTS "authority_authorityrole_users_user_id_fkey";

DROP TABLE IF EXISTS "authority_authorityrole_users";

ALTER TABLE
    IF EXISTS "authority_authorityrole"
    ADD COLUMN IF NOT EXISTS "description" TEXT NOT NULL,
    ADD COLUMN IF NOT EXISTS "is_system_role" BOOLEAN NOT NULL,
    ADD COLUMN IF NOT EXISTS "user_type" varchar(255) NOT NULL;


CREATE TABLE IF NOT EXISTS "authority_authorityrole_staff" (
    "id" SERIAL PRIMARY KEY,
    "authorityrole_id" uuid NOT NULL,
    "staff_id" uuid NOT NULL,
    UNIQUE ("staff_id", "authorityrole_id"),
    CONSTRAINT "authority_authorityrole_staff_staff_authorityrole_id_fkey" FOREIGN KEY ("authorityrole_id") REFERENCES "authority_authorityrole" ("id"),
    CONSTRAINT "authority_authorityrole_staff_staff_id_fkey" FOREIGN KEY ("staff_id") REFERENCES "staff_staff" ("id")
);
   
CREATE TABLE IF NOT EXISTS "authority_authorityrole_clients" (
    "id" SERIAL PRIMARY KEY,
    "authorityrole_id" uuid NOT NULL,
    "client_id" uuid NOT NULL,
    UNIQUE ("client_id", "authorityrole_id"),
    CONSTRAINT "authority_authorityrole_client_authorityrole_id_fkey" FOREIGN KEY ("authorityrole_id") REFERENCES "authority_authorityrole" ("id"),
    CONSTRAINT "authority_authorityrole_clients_client_id_fkey" FOREIGN KEY ("client_id") REFERENCES "clients_client" ("id")
);

CREATE TABLE IF NOT EXISTS "authority_authorityrole_caregivers" (
  "id" SERIAL PRIMARY KEY,
  "authorityrole_id" uuid NOT NULL,
  "caregiver_id" uuid NOT NULL,
   UNIQUE ("caregiver_id", "authorityrole_id"),
   CONSTRAINT "authority_authorityrole_caregiver_authorityrole_id_fkey" FOREIGN KEY ("authorityrole_id") REFERENCES "authority_authorityrole" ("id"),
   CONSTRAINT "authority_authorityrole_caregivers_caregiver_id_fkey" FOREIGN KEY ("caregiver_id") REFERENCES "caregivers_caregiver" ("id")
);

COMMIT;