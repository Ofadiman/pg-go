-- +goose Up
-- +goose StatementBegin
create extension IF not exists "uuid-ossp";

create table users (
  id uuid default uuid_generate_v4(),

  created_at timestamptz not null default now(),
  updated_at timestamptz default null,

  address jsonb not null,
  birthday timestamptz not null,
  email text not null,
  email_verified boolean not null,
  first_name text not null,
  gender text not null,
  last_name text not null,

  primary key (id),
  constraint unique_email unique (email)
);
-- +goose StatementEnd
