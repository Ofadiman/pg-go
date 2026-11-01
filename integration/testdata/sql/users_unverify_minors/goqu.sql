UPDATE "users" SET "email_verified"=FALSE,"updated_at"='2026-01-01T00:00:00Z' WHERE (age('2026-01-01T00:00:00Z', birthday) < interval '18 years') RETURNING *;
