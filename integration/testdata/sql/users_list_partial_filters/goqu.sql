SELECT COUNT(*) AS "count" FROM "users" WHERE (("email_verified" IS FALSE) AND ("birthday" > '2000-01-01T00:00:00Z')) LIMIT 1;

SELECT "address", "birthday", "created_at", "email", "email_verified", "first_name", "gender", "id", "last_name", "updated_at" FROM "users" WHERE (("email_verified" IS FALSE) AND ("birthday" > '2000-01-01T00:00:00Z')) ORDER BY "first_name" DESC, "created_at" ASC, "id" ASC LIMIT 5;
