SELECT count(*) FROM "users" WHERE first_name LIKE $1 AND gender = $2 AND email_verified = $3 AND birthday > $4 AND birthday < $5;

SELECT * FROM "users" WHERE first_name LIKE $1 AND gender = $2 AND email_verified = $3 AND birthday > $4 AND birthday < $5 ORDER BY "last_name","birthday" DESC,"id" LIMIT $6 OFFSET $7;
