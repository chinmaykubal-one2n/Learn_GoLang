-- Modify "teachers" table
ALTER TABLE "public"."teachers" DROP CONSTRAINT "uni_teachers_username", ALTER COLUMN "id" DROP DEFAULT, ALTER COLUMN "id" TYPE text, DROP COLUMN "created_at", DROP COLUMN "updated_at", DROP COLUMN "deleted_at", ALTER COLUMN "username" DROP NOT NULL, ALTER COLUMN "password" DROP NOT NULL, ALTER COLUMN "role" DROP NOT NULL, ADD COLUMN "email" text NULL;
-- Drop sequence used by serial column "id"
DROP SEQUENCE IF EXISTS "public"."teachers_id_seq";
