-- +goose Up
-- +goose StatementBegin
create table reactions (
  id uuid default uuid_generate_v4(),
  author_id uuid not null references users (id),
  comment_id uuid not null references comments (id),

  created_at timestamptz not null default now(),

  emoji text not null,

  primary key (id)
);
-- +goose StatementEnd
