-- name: CreateTenant :one
INSERT INTO tenants(name, slug)
VALUES ($1, $2)
RETURNING *;

-- name: GetTenantById :one
SELECT * FROM tenants
WHERE id = $1;

-- name: GetTenantBySlug :one
SELECT * FROM tenants
WHERE slug = $1;


-- name: ListTenants :many
SELECT * FROM tenants
ORDER BY created_at DESC;


-- name: UpdateTenantStatus :one
UPDATE tenants
SET status = $2, updated_at = NOW()
WHERE id = $1
RETURNING *; 