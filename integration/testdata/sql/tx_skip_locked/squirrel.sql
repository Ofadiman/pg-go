begin;

SELECT id, created_at, updated_at, address, birthday, email, email_verified, first_name, gender, last_name FROM users WHERE email_verified = $1 ORDER BY created_at ASC, id ASC LIMIT 2 FOR UPDATE SKIP LOCKED;

begin;

SELECT id, created_at, updated_at, address, birthday, email, email_verified, first_name, gender, last_name FROM users WHERE email_verified = $1 ORDER BY created_at ASC, id ASC LIMIT 2 FOR UPDATE SKIP LOCKED;

commit;

commit;
