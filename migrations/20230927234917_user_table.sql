-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE users (
  "id" uuid DEFAULT uuid_generate_v1(),
  "name" varchar(128) NOT NULL,
  "surname" varchar(128) NOT NULL,
  "email" varchar(128) NOT NULL UNIQUE,
  "birthday" date NULL,
  "biography" text NULL,
  "city" varchar(128) NOT NULL,
  "password" text NOT NULL,
  "created_at" timestamp NOT NULL,
  "updated_at" timestamp NOT NULL,
   PRIMARY KEY (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
