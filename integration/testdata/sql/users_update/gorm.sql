begin;

UPDATE "users" SET "email_verified"=$1,"last_name"=$2,"updated_at"=$3 WHERE id = $4 RETURNING *;

commit;
