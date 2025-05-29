-- +goose Up
-- +goose StatementBegin
create table posts (
  id uuid default uuid_generate_v4(),
  author_id uuid not null references users (id),

  created_at timestamptz not null default now(),
  updated_at timestamptz default null,

  content text not null,
  title text not null,

  primary key (id)
);
-- +goose StatementEnd
