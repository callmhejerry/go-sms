-- name: CreateAcademicSession :one
INSERT INTO academic_sessions (
    tenant_id, name, start_date, end_date, is_current
)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;


-- name: GetAcademicSessionById :one
SELECT * FROM academic_sessions
WHERE id = $1 AND tenant_id = $2;

-- name: ListAcademicSessions :many
SELECT * FROM academic_sessions
WHERE tenant_id = $1
ORDER BY start_date DESC;

-- name: GetCurrentAcademicSession :one
SELECT * FROM academic_sessions
WHERE tenant_id = $1 AND is_current = true;


-- name: SetCurrentAcademicSession :exec
UPDATE academic_sessions
SET is_current = false, updated_at = NOW()
WHERE tenant_id = $1 AND is_current = true;


-- name: UpdateAcademicSession :one
UPDATE academic_sessions
SET 
    name = COALESCE($3, name),
    start_date = COALESCE($4, start_date),
    end_date = COALESCE($5, end_date),
    is_current = COALESCE($6, is_current),
    updated_at = NOW()
WHERE id = $1 AND tenant_id = $2
RETURNING *;




