begin;

SELECT * FROM "users" WHERE email_verified = $1 ORDER BY created_at ASC, id ASC LIMIT $2 FOR UPDATE SKIP LOCKED;

begin;

SELECT * FROM "users" WHERE email_verified = $1 ORDER BY created_at ASC, id ASC LIMIT $2 FOR UPDATE SKIP LOCKED;

commit;

commit;
