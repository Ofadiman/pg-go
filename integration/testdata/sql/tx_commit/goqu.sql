begin;

INSERT INTO "users" ("address", "birthday", "created_at", "email", "email_verified", "first_name", "gender", "last_name", "updated_at") VALUES ('{"city":"Gdansk","country":"Poland","street":"Long Market","street_number":"1","zip_code":"80-001"}', '1991-02-02T00:00:00Z', '2026-01-01T00:00:00Z', 'new-user-1@example.com', FALSE, 'First1', 'male', 'Last1', NULL) RETURNING *;

INSERT INTO "posts" ("content", "created_at", "title", "updated_at", "user_id") VALUES ('Content of transaction post 1.', '2026-01-01T00:00:00Z', 'Transaction post 1', NULL, '<generated-1>'), ('Content of transaction post 2.', '2026-01-01T00:00:00Z', 'Transaction post 2', NULL, '<generated-1>'), ('Content of transaction post 3.', '2026-01-01T00:00:00Z', 'Transaction post 3', NULL, '<generated-1>') RETURNING *;

commit;
