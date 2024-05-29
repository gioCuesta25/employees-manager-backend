CREATE TABLE verification_token
(
  identifier TEXT NOT NULL,
  expires TIMESTAMPTZ NOT NULL,
  token TEXT NOT NULL,
 
  PRIMARY KEY (identifier, token)
);
 
CREATE TABLE accounts
(
  id SERIAL,
  "userId" INTEGER NOT NULL,
  type VARCHAR(255) NOT NULL,
  provider VARCHAR(255) NOT NULL,
  "providerAccountId" VARCHAR(255) NOT NULL,
  refresh_token TEXT,
  access_token TEXT,
  expires_at BIGINT,
  id_token TEXT,
  scope TEXT,
  session_state TEXT,
  token_type TEXT,
 
  PRIMARY KEY (id)
);
 
CREATE TABLE sessions
(
  id SERIAL,
  "userId" INTEGER NOT NULL,
  expires TIMESTAMPTZ NOT NULL,
  "sessionToken" VARCHAR(255) NOT NULL,
 
  PRIMARY KEY (id)
);
 
CREATE TABLE users
(
  id SERIAL,
  name VARCHAR(255),
  email VARCHAR(255),
  "emailVerified" TIMESTAMPTZ,
  image TEXT,
 
  PRIMARY KEY (id)
);
 

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
  "owner" SERIAL NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz
);

-- CREATE TABLE "users" (
--   "id" UUID PRIMARY KEY DEFAULT (uuid_generate_v4()),
--   "full_name" varchar(100) NOT NULL,
--   "email" varchar(100) UNIQUE NOT NULL,
--   "password" varchar NOT NULL,
--   "created_at" timestamptz NOT NULL DEFAULT (now()),
--   "updated_at" timestamptz
-- );

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
