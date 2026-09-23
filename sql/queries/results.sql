-- name: UpsertResult :one
INSERT INTO results(
    tenant_id, student_id, subject_id, class_arm_id,
    academic_session_id, total_score, max_total, percentage,
    grade, remark
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
ON CONFLICT (tenant_id, student_id, subject_id, academic_session_id)
DO UPDATE SET
    total_score = EXCLUDED.total_score,
    max_total = EXCLUDED.max_total,
    percentage = EXCLUDED.percentage,
    grade = EXCLUDED.grade,
    remark = EXCLUDED.remark,
    updated_at = NOW()
RETURNING *;

-- name: GetResultsByStudent :many
SELECT r.*, sub.name AS subject_name
FROM results r
JOIN subjects sub ON sub.id = r.subject_id
WHERE r.tenant_id = $1 AND r.student_id = $2 AND r.academic_session_id = $3
ORDER BY sub.name;


-- name: GetScoresForComputation :many
SELECT
    s.student_id,
    s.subject_id,
    s.class_arm_id,
    s.score,
    at.max_score,
    at.weight
FROM scores s
JOIN assessment_types at ON at.id = s.assessment_type_id
WHERE s.tenant_id = $1
    AND s.class_arm_id = $2
    AND s.academic_session_id = $3;
