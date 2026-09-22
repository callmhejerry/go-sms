
-- name: CreateSubject :one
INSERT INTO subjects (tenant_id, name, code)
VALUES ($1,$2,$3)
RETURNING *;


-- name: ListSubjects :many
SELECT * FROM subjects
WHERE tenant_id = $1
ORDER BY name;


-- name: AddSubjectToClass :one
INSERT INTO class_subjects (tenant_id, class_id, subject_id)
VALUES ($1, $2, $3)
RETURNING *;


-- name: ListSubjectsByClass :many
SELECT s.* FROM subjects s
INNER JOIN class_subjects cs ON cs.subject_id = s.id
WHERE cs.tenant_id = $1 AND cs.class_id = $2
ORDER BY s.name;


-- name: AssignTeacher :one
INSERT INTO teacher_assignments (
    tenant_id, user_id, subject_id, academic_session_id, class_id, class_arm_id
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: ListTeacherAssignments :many
SELECT * FROM teacher_assignments
WHERE tenant_id = $1 AND academic_session_id = $2
ORDER BY created_at;