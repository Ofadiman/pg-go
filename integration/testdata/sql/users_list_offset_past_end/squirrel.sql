SELECT count(*) FROM users WHERE first_name LIKE $1 AND gender = $2;

SELECT id, created_at, updated_at, address, birthday, email, email_verified, first_name, gender, last_name FROM users WHERE first_name LIKE $1 AND gender = $2 ORDER BY created_at DESC, id ASC LIMIT 20 OFFSET 50;
