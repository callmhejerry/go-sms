-- name: CreateAuditLog :one
INSERT INTO audit_logs (
    tenant_id, user_id, action, resource_type, resource_id,
    metadata, ip_address, user_agent, request_id
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
) RETURNING *;

-- name: ListAuditLogs :many
SELECT * FROM audit_logs
WHERE tenant_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;