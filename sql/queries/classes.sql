-- name: CreateClass :one
INSERT INTO classes (
    tenant_id, name , level_order
)
VALUES ($1, $2, $3)
RETURNING *;


-- name: ListClasses :many
SELECT * FROM classes
WHERE tenant_id = $1
ORDER BY level_order, name;


-- name: GetClassByID :one
SELECT * FROM classes
WHERE id = $1 AND tenant_id = $2;


-- name: CreateClassArm :one
INSERT INTO class_arms (
    tenant_id, class_id, name
)
VALUES ($1, $2, $3)
RETURNING *;

-- name: ListClassArms :many
SELECT * FROM class_arms
WHERE tenant_id = $1 AND class_id = $2
ORDER BY name;


-- name: ListAllClassArms :many
SELECT ca.* , c.name as class_name
FROM class_arms ca
JOIN classes c ON c.id = ca.class_id
WHERE ca.tenant_id = $1
ORDER by c.level_order, ca.name;
