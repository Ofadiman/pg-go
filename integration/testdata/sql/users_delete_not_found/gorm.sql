begin;

DELETE FROM "users" WHERE id = $1 RETURNING *;

commit;
