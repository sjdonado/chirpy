-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (user_id, token) VALUES ($1, $2) RETURNING token;

-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens SET revoked_at = NOW() WHERE token = $1;

-- name: GetUserFromRefreshToken :one
SELECT * FROM users WHERE id = (SELECT user_id FROM refresh_tokens WHERE token = $1) LIMIT 1;
