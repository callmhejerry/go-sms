-- name: CreateUser :one
INSERT INTO users (
    tenant_id, email, first_name, last_name, password_hash
)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;


-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1 AND tenant_id = $2;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 AND tenant_id = $2;

-- name: ListUsersByTenant :many
SELECT * FROM users
WHERE tenant_id = $1
ORDER BY created_at DESC;