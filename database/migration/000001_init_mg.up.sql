
CREATE TABLE "id_types" (
  "id" UUID PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "name" varchar(100) NOT NULL,
  "code" varchar(100) NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz
);

INSERT INTO id_types (name, code) VALUES ('Cédula de Ciudadanía', 'CC');

CREATE TABLE "positions" (
  "id" UUID PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "company_id" UUID NOT NULL,
  "department_id" UUID NOT NULL,
  "name" varchar(100) NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz
);

CREATE TABLE "departments" (
  "id" UUID PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "name" varchar(100) NOT NULL,
  "company_id" UUID NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz
);

CREATE TABLE "employees" (
  "id" UUID PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "name" varchar(100) NOT NULL,
  "last_name" varchar(100) NOT NULL,
  "phone_number" varchar(100) NOT NULL,
  "email" varchar(100) NOT NULL,
  "id_type" UUID NOT NULL,
  "id_number" varchar(100) NOT NULL,
  "admission_date" timestamptz NOT NULL,
  "salary" bigint,
  "position_id" UUID NOT NULL,
  "departament_id" UUID NOT NULL,
  "company_id" UUID NOT NULL,
  "picture_url" varchar,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz
);

CREATE TABLE "companies" (
  "id" UUID PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "name" varchar(100) NOT NULL,
  "owner" UUID NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz
);

CREATE TABLE "users" (
  "id" UUID PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "full_name" varchar(100) NOT NULL,
  "email" varchar(100) UNIQUE NOT NULL,
  "password" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz
);

CREATE INDEX ON "positions" ("company_id");

CREATE INDEX ON "positions" ("department_id");

CREATE INDEX ON "employees" ("position_id");

CREATE INDEX ON "employees" ("departament_id");

CREATE INDEX ON "employees" ("company_id");

ALTER TABLE "positions" ADD FOREIGN KEY ("company_id") REFERENCES "companies" ("id");

ALTER TABLE "positions" ADD FOREIGN KEY ("department_id") REFERENCES "departments" ("id");

ALTER TABLE "departments" ADD FOREIGN KEY ("company_id") REFERENCES "companies" ("id");

ALTER TABLE "employees" ADD FOREIGN KEY ("id_type") REFERENCES "id_types" ("id");

ALTER TABLE "employees" ADD FOREIGN KEY ("position_id") REFERENCES "positions" ("id");

ALTER TABLE "employees" ADD FOREIGN KEY ("departament_id") REFERENCES "departments" ("id");

ALTER TABLE "employees" ADD FOREIGN KEY ("company_id") REFERENCES "companies" ("id");

ALTER TABLE "companies" ADD FOREIGN KEY ("owner") REFERENCES "users" ("id");
