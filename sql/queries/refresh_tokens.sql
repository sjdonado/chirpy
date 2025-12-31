-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (user_id, token, expires_at) VALUES ($1, $2, $3) RETURNING token;

-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens SET revoked_at = NOW(), updated_at = NOW() WHERE token = $1;

-- name: GetRefreshtoken :one
SELECT * FROM refresh_tokens WHERE token = $1 LIMIT 1;
