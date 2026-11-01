SELECT COUNT(*) AS "count" FROM "users" LIMIT 1;

SELECT "address", "birthday", "created_at", "email", "email_verified", "first_name", "gender", "id", "last_name", "updated_at" FROM "users" ORDER BY "created_at" DESC, "id" ASC LIMIT 20;
