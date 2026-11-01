begin;

INSERT INTO "users" ("address", "birthday", "created_at", "email", "email_verified", "first_name", "gender", "last_name", "updated_at") VALUES ('{"city":"Gdansk","country":"Poland","street":"Long Market","street_number":"1","zip_code":"80-001"}', '1991-02-02T00:00:00Z', '2026-01-01T00:00:00Z', 'new-user-1@example.com', FALSE, 'First1', 'male', 'Last1', NULL) RETURNING *;

SAVEPOINT sp_1;

INSERT INTO "users" ("address", "birthday", "created_at", "email", "email_verified", "first_name", "gender", "last_name", "updated_at") VALUES ('{"city":"Gdansk","country":"Poland","street":"Long Market","street_number":"2","zip_code":"80-002"}', '1992-03-03T00:00:00Z', '2026-01-01T00:00:00Z', 'new-user-2@example.com', TRUE, 'First2', 'female', 'Last2', NULL) RETURNING *;

ROLLBACK TO SAVEPOINT sp_1;

commit;
