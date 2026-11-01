begin;

UPDATE "users" SET "last_name"=$1,"updated_at"=$2 WHERE id = $3 RETURNING *;

commit;
