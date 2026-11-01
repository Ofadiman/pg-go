SELECT count(*) FROM "users" WHERE email_verified = $1 AND birthday > $2;

SELECT * FROM "users" WHERE email_verified = $1 AND birthday > $2 ORDER BY "first_name" DESC,"created_at","id" LIMIT $3;
