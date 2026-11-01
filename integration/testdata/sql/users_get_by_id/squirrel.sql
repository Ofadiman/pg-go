SELECT id, created_at, updated_at, address, birthday, email, email_verified, first_name, gender, last_name FROM users WHERE id = $1;
