-- name: UpsertScore :one
INSERT INTO scores(
    tenant_id, student_id, subject_id,
    assessment_type_id, class_arm_id, academic_session_id,
    score, recorded_by
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
)
ON CONFLICT (tenant_id, student_id, subject_id, assessment_type_id, academic_session_id)
DO UPDATE SET
    score = EXCLUDED.score,
    recorded_by = EXCLUDED.recorded_by,
    updated_at = NOW()
RETURNING *;


-- name: GetScoresByStudent :many
SELECT s.*, at.name AS assessment_name, sub.name AS subject_name
FROM scores s
JOIN assessment_types at ON at.id = s.assessment_type_id
JOIN subjects sub ON sub.id = s.subject_id
WHERE s.tenant_id = $1 AND s.student_id = $2 AND s.academic_session_id = $3
ORDER BY sub.name, at.name;

-- name: GetScoresByClassArmAndSubject :many
SELECT s.*, st.first_name, st.last_name, st.admission_number
FROM scores s
JOIN students st ON st.id = s.student_id
WHERE s.tenant_id = $1
  AND s.class_arm_id = $2
  AND s.subject_id = $3
  AND s.assessment_type_id = $4
  AND s.academic_session_id = $5
ORDER BY st.last_name, st.first_name;