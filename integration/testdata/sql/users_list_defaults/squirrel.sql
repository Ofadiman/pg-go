SELECT count(*) FROM users;

SELECT id, created_at, updated_at, address, birthday, email, email_verified, first_name, gender, last_name FROM users ORDER BY created_at DESC, id ASC LIMIT 20 OFFSET 0;
