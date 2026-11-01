begin;

INSERT INTO users (address,birthday,created_at,email,email_verified,first_name,gender,last_name) VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING *;

INSERT INTO posts (user_id,created_at,title,content) VALUES ($1,$2,$3,$4),($5,$6,$7,$8),($9,$10,$11,$12) RETURNING *;

commit;
