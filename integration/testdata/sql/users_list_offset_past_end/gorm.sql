SELECT count(*) FROM "users" WHERE first_name LIKE $1 AND gender = $2;

SELECT * FROM "users" WHERE first_name LIKE $1 AND gender = $2 ORDER BY "created_at" DESC,"id" LIMIT $3 OFFSET $4;
