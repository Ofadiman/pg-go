-- +goose Up
-- +goose StatementBegin
create table posts (
  id uuid primary key default uuid_generate_v4(),
  author_id UUID not null references users (id),
  title text not null,
  content text not null,
  created_at timestamptz not null default now(),
  updated_at timestamptz default null
);
-- +goose StatementEnd
