begin;

INSERT INTO "users" ("created_at","updated_at","address","birthday","email","email_verified","first_name","gender","last_name") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING "id";

INSERT INTO "posts" ("user_id","created_at","updated_at","title","content") VALUES ($1,$2,$3,$4,$5),($6,$7,$8,$9,$10) RETURNING "id";

rollback;
