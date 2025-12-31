-- name: CreateChirp :one
INSERT INTO chirps (user_id, body) VALUES ($1, $2) RETURNING *;

-- name: FilterChirps :many
SELECT * FROM chirps WHERE user_id = $1 ORDER BY created_at ASC;

-- name: GetChirp :one
SELECT * FROM chirps WHERE id = $1 LIMIT 1;

-- name: DeleteChirp :exec
DELETE FROM chirps WHERE id = $1;
