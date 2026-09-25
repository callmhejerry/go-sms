-- name: CreateUser :one
INSERT INTO users (
    email, first_name, last_name, password_hash
)
VALUES ($1, $2, $3, $4)
RETURNING *;


-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 ;

-- name: ListUsersByTenant :many
SELECT u.* FROM users u
INNER JOIN tenant_users tu ON tu.user_id = u.id
WHERE tu.tenant_id = $1
ORDER BY u.created_at DESC;