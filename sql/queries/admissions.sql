-- name: CreateAdmission :one
INSERT INTO admissions (
    tenant_id, academic_session_id,
    first_name, last_name, middle_name,
    gender, date_of_birth,
    preferred_class_id,
    parent_first_name, parent_last_name,
    parent_phone_number, parent_email,
    parent_relationship
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
RETURNING *;


-- name: GetAdmissionById :one
SELECT * FROM admissions
WHERE id = $1 AND tenant_id = $2;


-- name: ListAdmissions :many
SELECT * FROM admissions
WHERE tenant_id = $1
ORDER BY created_at DESC;

-- name: ListAdmissionsByStatus :many
SELECT * FROM admissions
WHERE tenant_id = $1 AND status = $2
ORDER BY created_at DESC;


-- name: UpdateAdmissionStatus :one
UPDATE admissions
SET
    status = $3,
    reviewed_by = $4,
    reviewed_at = NOW(),
    rejection_reason = $5,
    admission_number = $6,
    student_id = $7,
    updated_at = NOW()
WHERE id = $1 AND tenant_id = $2
RETURNING *;
