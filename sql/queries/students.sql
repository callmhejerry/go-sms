-- name: CreateStudent :one
INSERT INTO students(
    tenant_id, admission_number,
    first_name, last_name, middle_name,
    gender, date_of_birth, status,
    current_class_arm_id, admission_session_id
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING *;

-- name: GetStudentByID :one
SELECT * FROM students
WHERE id = $1 AND tenant_id = $2;

-- name: ListStudents :many
SELECT * FROM students
WHERE tenant_id = $1
ORDER BY last_name, first_name;


-- name: CreateParent :one
INSERT INTO parents(
    tenant_id, first_name, last_name,
    email, phone_number, address
)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;


-- name: GetParentById :one
SELECT * FROM parents
WHERE id = $1 AND tenant_id = $2;

-- name: LinkStudentParent :exec
INSERT INTO student_parents(
    student_id, parent_id, relationship,
    is_primary
)
VALUES ($1, $2, $3, $4)
ON CONFLICT DO NOTHING;

-- name: GetParentsByStudent :many
SELECT p.*, sp.relationship, sp.is_primary FROM parents p
INNER JOIN student_parents sp ON sp.parent_id = p.id
WHERE sp.student_id = $1 AND p.tenant_id = $2;

-- name: GetStudentsByParent :many
SELECT s.* , sp.relationship, sp.is_primary FROM students s
INNER JOIN student_parents sp ON sp.student_id = s.id
WHERE sp.parent_id = $1 AND s.tenant_id = $2;