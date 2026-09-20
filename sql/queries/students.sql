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


-- name: GetStudentById :one
SELECT s.*, 
    c.name as class_name,
    ca.name as class_arm_name,
    sess.name as admission_session_name
FROM students s
LEFT JOIN class_arms ca ON ca.id = s.current_class_arm_id
LEFT JOIN classes c ON c.id = ca.class_id
LEFT JOIN academic_sessions sess ON sess.id = s.academic_session_id
WHERE s.id = $1 AND s.tenant_id = $2;


-- name: SearchStudents :many
SELECT * FROM students
WHERE tenant_id = $1
    AND (
        $2::text IS NULL OR
        first_name ILIKE '%' || $2 || '%' OR 
        last_name ILIKE '%' || $2 || '%' OR
        admission_number ILIKE '%' || $2 || '%'
    )
    AND ($3::uuid IS NULL OR current_class_arm_id = $3)
    AND ($4::text IS NULL OR status = $4)
ORDER BY last_name, first_name;


-- name: UpdateStudent :one
UPDATE students
SET 
    first_name = COALESCE($3, first_name),
    last_name = COALESCE($4, last_name),
    middle_name = COALESCE($5, middle_name),
    gender = COALESCE($6, gender),
    date_of_birth = COALESCE($7, date_of_birth),
    current_class_arm_id = COALESCE($8, current_class_arm_id),
    status = COALESCE($9, status),
    updated_at = NOW()
WHERE id = $1 AND tenant_id = $2
RETURNING *;