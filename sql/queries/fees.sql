-- name: CreateFeeType :one
INSERT INTO fee_types(
    tenant_id, name, description,
    is_optional
)
VALUES ($1, $2, $3, $4)
RETURNING *;


-- name: ListFeeTypes :many
SELECT * FROM fee_types
WHERE tenant_id = $1
ORDER BY name;


-- name: CreateFeeStructure :one
INSERT INTO fee_structures (
    tenant_id, fee_type_id, academic_session_id,
    class_id, amount_kobo, due_date
)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;


-- name: ListFeeStructures :many
SELECT fs.*, ft.name AS fee_type_name
FROM fee_structures fs
INNER JOIN fee_types ft ON ft.id = fs.fee_type_id
WHERE fs.tenant_id = $1;


-- name: ListFeeStructuresBySession :many
SELECT fs.*, ft.name AS fee_type_name
FROM fee_structures fs
JOIN fee_types ft ON ft.id = fs.fee_type_id
WHERE fs.tenant_id = $1 AND fs.academic_session_id = $2
ORDER BY ft.name;


-- name: GetFeeStructureById :one
SELECT * FROM fee_structures
WHERE id = $1 AND tenant_id = $2;