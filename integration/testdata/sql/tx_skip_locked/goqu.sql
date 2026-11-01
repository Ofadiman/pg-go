begin;

SELECT "address", "birthday", "created_at", "email", "email_verified", "first_name", "gender", "id", "last_name", "updated_at" FROM "users" WHERE ("email_verified" IS FALSE) ORDER BY "created_at" ASC, "id" ASC LIMIT 2 FOR UPDATE SKIP LOCKED;

begin;

SELECT "address", "birthday", "created_at", "email", "email_verified", "first_name", "gender", "id", "last_name", "updated_at" FROM "users" WHERE ("email_verified" IS FALSE) ORDER BY "created_at" ASC, "id" ASC LIMIT 2 FOR UPDATE SKIP LOCKED;

commit;

commit;
