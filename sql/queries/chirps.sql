-- name: CreateChirp :one
INSERT INTO chirps (user_id, body)
VALUES ($1, $2)
RETURNING *;
