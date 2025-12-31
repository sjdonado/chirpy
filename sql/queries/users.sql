-- name: CreateUser :one
INSERT INTO users (email, hashed_password) VALUES ($1, $2) RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1 LIMIT 1;

-- name: UpdateUser :one
UPDATE users SET email = $2, hashed_password = $3 WHERE id = $1 RETURNING *;

-- name: UpdateUserChirpyRed :one
UPDATE users SET is_chirpy_red = $2 WHERE id = $1 RETURNING *;

-- name: DeleteAllUsers :exec
DELETE FROM users;
