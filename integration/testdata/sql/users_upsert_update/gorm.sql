begin;

INSERT INTO "users" ("created_at","updated_at","address","birthday","email","email_verified","first_name","gender","last_name") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) ON CONFLICT ("email") DO UPDATE SET "first_name"=excluded.first_name,"last_name"=excluded.last_name,"updated_at"=$10 RETURNING *;

commit;
