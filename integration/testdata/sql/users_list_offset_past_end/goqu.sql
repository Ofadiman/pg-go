SELECT COUNT(*) AS "count" FROM "users" WHERE (("first_name" LIKE 'An%') AND ("gender" = 'female')) LIMIT 1;

SELECT "address", "birthday", "created_at", "email", "email_verified", "first_name", "gender", "id", "last_name", "updated_at" FROM "users" WHERE (("first_name" LIKE 'An%') AND ("gender" = 'female')) ORDER BY "created_at" DESC, "id" ASC LIMIT 20 OFFSET 50;
