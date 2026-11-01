SELECT * FROM "users" WHERE email_verified = $1 AND ((gender = $2 AND age($3, birthday) >= interval '21 years') OR (gender = $4 AND age($5, birthday) >= interval '18 years')) ORDER BY id;
