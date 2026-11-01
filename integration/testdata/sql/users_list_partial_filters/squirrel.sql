SELECT count(*) FROM users WHERE email_verified = $1 AND birthday > $2;

SELECT id, created_at, updated_at, address, birthday, email, email_verified, first_name, gender, last_name FROM users WHERE email_verified = $1 AND birthday > $2 ORDER BY first_name DESC, created_at ASC, id ASC LIMIT 5 OFFSET 0;
