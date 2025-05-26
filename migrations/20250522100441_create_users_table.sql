-- +goose Up
-- +goose StatementBegin
create extension IF not exists "uuid-ossp";

create table users (
  id uuid primary key default uuid_generate_v4(),
  first_name text not null,
  last_name text not null,
  email text not null,
  email_verified boolean not null,
  phone_number text not null,
  gender text not null,
  city text not null,
  country text not null,
  street text not null,
  street_number text not null,
  zip_code text not null,
  birthday timestamptz not null,
  created_at timestamptz not null default now(),
  updated_at timestamptz default null,
  constraint unique_email unique (email),
  constraint unique_phone_number unique (phone_number)
);
-- +goose StatementEnd
