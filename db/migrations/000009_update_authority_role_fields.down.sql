BEGIN;

CREATE TABLE IF NOT EXISTS "authority_authorityrole_users" (
  "id" SERIAL PRIMARY KEY,
  "authorityrole_id" uuid NOT NULL,
  "user_id" uuid NOT NULL,
   UNIQUE ("user_id", "authorityrole_id"),
   CONSTRAINT "authority_authorityrole_users_authorityrole_id_fkey" FOREIGN KEY ("authorityrole_id") REFERENCES "authority_authorityrole" ("id"),
   CONSTRAINT "authority_authorityrole_users_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users_user" ("id")
);

ALTER TABLE
    IF EXISTS "authority_authorityrole"
    DROP COLUMN IF EXISTS "description",
    DROP COLUMN IF EXISTS "is_system_role",
    DROP COLUMN IF EXISTS "user_type";

ALTER TABLE
    IF EXISTS "authority_authorityrole_staff"
    DROP CONSTRAINT IF EXISTS "authority_authorityrole_staff_authorityrole_id_fkey",
    DROP CONSTRAINT IF EXISTS "authority_authorityrole_staff_staff_id_fkey";

DROP TABLE IF EXISTS "authority_authorityrole_staff";

ALTER TABLE
    IF EXISTS "authority_authorityrole_clients"
    DROP CONSTRAINT IF EXISTS "authority_authorityrole_client_authorityrole_id_fkey",
    DROP CONSTRAINT IF EXISTS "authority_authorityrole_clients_client_id_fkey";

DROP TABLE IF EXISTS "authority_authorityrole_clients";

ALTER TABLE
    IF EXISTS "authority_authorityrole_caregivers"
    DROP CONSTRAINT IF EXISTS "authority_authorityrole_caregiver_authorityrole_id_fkey",
    DROP CONSTRAINT IF EXISTS "authority_authorityrole_caregivers_caregiver_id_fkey";

DROP TABLE IF EXISTS "authority_authorityrole_caregivers";

COMMIT;