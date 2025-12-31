-- name: CreateChirp :one
INSERT INTO chirps (user_id, body) VALUES ($1, $2) RETURNING *;

-- name: QueryChirps :many
SELECT * FROM chirps
WHERE (sqlc.narg('user_id')::uuid IS NULL OR user_id = sqlc.narg('user_id')::uuid)
ORDER BY
  CASE WHEN sqlc.arg('sort') = 'desc' THEN created_at END DESC,
  CASE WHEN sqlc.arg('sort') = 'asc' THEN created_at END ASC;

-- name: GetChirp :one
SELECT * FROM chirps WHERE id = $1 LIMIT 1;

-- name: DeleteChirp :exec
DELETE FROM chirps WHERE id = $1;
