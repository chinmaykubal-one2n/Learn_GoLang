-- Create "teachers" table
CREATE TABLE "public"."teachers" ("id" bigserial NOT NULL, "created_at" timestamptz NULL, "updated_at" timestamptz NULL, "deleted_at" timestamptz NULL, "username" text NOT NULL, "password" text NOT NULL, "role" text NOT NULL, PRIMARY KEY ("id"), CONSTRAINT "uni_teachers_username" UNIQUE ("username"));
-- Create index "idx_teachers_deleted_at" to table: "teachers"
CREATE INDEX "idx_teachers_deleted_at" ON "public"."teachers" ("deleted_at");
