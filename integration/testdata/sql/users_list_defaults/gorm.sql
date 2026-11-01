SELECT count(*) FROM "users";

SELECT * FROM "users" ORDER BY "created_at" DESC,"id" LIMIT $1;
