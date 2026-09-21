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

-- name: ListStudentsPage :many
SELECT * FROM students
WHERE tenant_id = $1
ORDER BY last_name, first_name, id
LIMIT sqlc.arg('limit')
OFFSET sqlc.arg('offset');

-- name: ListStudentsCursor :many
SELECT * FROM students
WHERE tenant_id = $1
    AND (
    sqlc.narg('cursor_last_name')::text IS NULL
    OR (last_name, first_name, id) > (sqlc.arg('cursor_last_name'), sqlc.arg('cursor_first_name'), sqlc.arg('cursor_id'))
)
ORDER BY last_name, first_name, id
LIMIT sqlc.arg('limit');


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


-- name: GetStudentProfile :one
SELECT s.*, 
    c.name as class_name,
    ca.name as class_arm_name,
    sess.name as admission_session_name
FROM students s
LEFT JOIN class_arms ca ON ca.id = s.current_class_arm_id
LEFT JOIN classes c ON c.id = ca.class_id
LEFT JOIN academic_sessions sess ON sess.id = s.academic_session_id
WHERE s.id = $1 AND s.tenant_id = $2;


-- name: SearchStudentsPage :many
SELECT * FROM students
WHERE tenant_id = $1
    AND (
        sqlc.narg('search')::text IS NULL OR
        first_name ILIKE '%' || sqlc.narg('search') || '%' OR 
        last_name ILIKE '%' || sqlc.narg('search') || '%' OR
        admission_number ILIKE '%' || sqlc.narg('search') || '%'
    )
    AND (
        sqlc.narg('current_class_arm')::uuid IS NULL
        OR current_class_arm_id = sqlc.narg('current_class_arm')::uuid
    )
    AND (
        sqlc.narg('status')::text IS NULL
        OR status = sqlc.narg('status')::text
    )
ORDER BY last_name, first_name, id
LIMIT sqlc.arg('limit')
OFFSET sqlc.arg('offset');

-- name: CountSearchStudents :one
SELECT COUNT(*) FROM students
WHERE tenant_id = $1
    AND (
        sqlc.narg('search')::text IS NULL OR
        first_name ILIKE '%' || sqlc.narg('search') || '%' OR 
        last_name ILIKE '%' || sqlc.narg('search') || '%' OR
        admission_number ILIKE '%' || sqlc.narg('search') || '%'
    )
    AND (
        sqlc.narg('current_class_arm')::uuid IS NULL
        OR current_class_arm_id = sqlc.narg('current_class_arm')::uuid
    )
    AND (
        sqlc.narg('status')::text IS NULL
        OR status = sqlc.narg('status')::text
);

-- name: SearchStudentsCursor :many
SELECT * FROM students
WHERE tenant_id = $1
    AND (
        sqlc.narg('search')::text IS NULL OR
        first_name ILIKE '%' || sqlc.narg('search') || '%' OR 
        last_name ILIKE '%' || sqlc.narg('search') || '%' OR
        admission_number ILIKE '%' || sqlc.narg('search') || '%'
    )
    AND (
        sqlc.narg('current_class_arm')::uuid IS NULL
        OR current_class_arm_id = sqlc.narg('current_class_arm')::uuid
    )
    AND (
        sqlc.narg('status')::text IS NULL
        OR status = sqlc.narg('status')::text
    )
    AND (
        sqlc.narg('cursor_last_name')::text IS NULL
        OR (last_name, first_name, id) > (sqlc.narg('cursor_last_name')::text, sqlc.narg('cursor_first_name')::text, sqlc.narg('cursor_id')::uuid)
    )
ORDER BY last_name, first_name, id
LIMIT sqlc.arg('limit');


-- name: UpdateStudent :one
UPDATE students
SET 
    first_name = COALESCE(sqlc.narg('first_name'), first_name),
    last_name = COALESCE(sqlc.narg('last_name'), last_name),
    middle_name = COALESCE(sqlc.narg('middle_name'), middle_name),
    gender = COALESCE(sqlc.narg('gender'), gender),
    date_of_birth = COALESCE(sqlc.narg('date_of_birth'), date_of_birth),
    current_class_arm_id = COALESCE(sqlc.narg('current_class_arm_id'), current_class_arm_id),
    status = COALESCE(sqlc.narg('status'), status),
    updated_at = NOW()
WHERE id = $1 AND tenant_id = $2
RETURNING *;