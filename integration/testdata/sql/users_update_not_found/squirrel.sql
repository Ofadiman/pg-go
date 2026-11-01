UPDATE users SET updated_at = $1, last_name = $2 WHERE id = $3 RETURNING *;
