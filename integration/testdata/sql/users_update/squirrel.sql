UPDATE users SET updated_at = $1, last_name = $2, email_verified = $3 WHERE id = $4 RETURNING *;
