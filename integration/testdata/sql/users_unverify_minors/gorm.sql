begin;

UPDATE "users" SET "email_verified"=$1,"updated_at"=$2 WHERE age($3, birthday) < interval '18 years' RETURNING *;

commit;
