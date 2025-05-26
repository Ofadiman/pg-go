-- +goose Up
-- +goose StatementBegin
CREATE TABLE posts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

    author_id UUID NOT NULL REFERENCES users(id),

    title TEXT NOT NULL,
    content TEXT NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at timestamptz DEFAULT NULL
);
-- +goose StatementEnd
