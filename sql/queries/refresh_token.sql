-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (refresh_token_hash, user_id, expires_at)
VALUES ($1, $2, $3)
RETURNING *;

-- name: DeleteRefreshToken :exec
DELETE FROM refresh_tokens
WHERE id = $1;

-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens
SET
    revoked_at = NOW()
WHERE id = $1;

-- name: GetRefreshToken :one
SELECT * FROM refresh_tokens
WHERE refresh_token_hash = $1;