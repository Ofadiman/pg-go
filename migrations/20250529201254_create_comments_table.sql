-- +goose Up
-- +goose StatementBegin
create table comments (
  id uuid default uuid_generate_v4(),
  author_id uuid not null references users (id),
  post_id uuid not null references posts (id),
  parent_id uuid references comments (id) default null,

  created_at timestamptz not null default now(),
  updated_at timestamptz default null,

  content text not null,

  primary key (id)
);
-- +goose StatementEnd
