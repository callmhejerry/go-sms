-- name: CreateAssessmentType :one
INSERT INTO assessment_types (tenant_id, name, max_score, weight, is_exam)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: ListAssessmentTypes :many
SELECT * FROM assessment_types
WHERE tenant_id = $1
ORDER BY name;

-- name: GetAssessmentTypeByID :one
SELECT * FROM assessment_types
WHERE id = $1 AND tenant_id = $2;