-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at timestamptz NOT NULL,
    updated_at timestamptz DEFAULT NULL,
    first_name text NOT NULL,
    last_name text NOT NULL,
    email text NOT NULL,
    email_verified boolean NOT NULL,
    phone_number text NOT NULL,
    gender text NOT NULL,
    city text NOT NULL,
    country text NOT NULL,
    street text NOT NULL,
    street_number text NOT NULL,
    zip_code text NOT NULL,
    birthday timestamptz NOT NULL,
    CONSTRAINT unique_email UNIQUE (email),
    CONSTRAINT unique_phone_number UNIQUE (phone_number)
);
-- +goose StatementEnd
