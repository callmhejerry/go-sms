-- name: CreateRole :one
INSERT INTO roles (
    tenant_id, name,
    description
)
VALUES ($1, $2, $3)
RETURNING *;


-- name: GetRoleByName :one
SELECT * FROM roles
WHERE tenant_id = $1 AND name = $2;

-- name: ListRolesByTenant :many
SELECT * FROM roles
WHERE tenant_id = $1
ORDER BY name;


-- name: AssignRoleToUser :exec
INSERT INTO user_roles (
    user_id, role_id
)
VALUES ($1, $2)
ON CONFLICT DO NOTHING;


-- name: RemoveRoleFromUser :exec
DELETE FROM user_roles
WHERE user_id = $1 AND role_id = $2;


-- name: GetUserRoles :many
SELECT r.* FROM roles r
INNER JOIN user_roles ur ON ur.role_id = r.id
WHERE ur.user_id = $1;

-- name: UserHasRole :one
SELECT EXISTS (
    SELECT 1 FROM user_roles ur
    INNER JOIN roles r ON r.id = ur.role_id
    WHERE ur.user_id = $1 AND r.name = $2
)
AS has_role;

