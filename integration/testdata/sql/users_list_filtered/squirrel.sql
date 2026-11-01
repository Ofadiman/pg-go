SELECT count(*) FROM users WHERE first_name LIKE $1 AND gender = $2 AND email_verified = $3 AND birthday > $4 AND birthday < $5;

SELECT id, created_at, updated_at, address, birthday, email, email_verified, first_name, gender, last_name FROM users WHERE first_name LIKE $1 AND gender = $2 AND email_verified = $3 AND birthday > $4 AND birthday < $5 ORDER BY last_name ASC, birthday DESC, id ASC LIMIT 3 OFFSET 3;
