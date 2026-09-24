-- name: CreateInventoryCategory :one
INSERT INTO inventory_categories (tenant_id, name)
VALUES ($1, $2)
RETURNING *;

-- name: ListInventoryCategories :many
SELECT * FROM inventory_categories
WHERE tenant_id = $1
ORDER BY name;

-- name: CreateInventoryItem :one
INSERT INTO inventory_items(
    tenant_id, category_id, name, code, description,
    unit, quantity_in_stock, reorder_level, unit_cost_kobo
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;

-- name: ListInventoryItems :many
SELECT i.*, c.name AS category_name
FROM inventory_items i
JOIN inventory_categories c ON c.id = i.category_id
WHERE i.tenant_id = $1
ORDER BY i.name;

-- name: GetInventoryItemByID :one
SELECT * FROM inventory_items
WHERE id = $1 AND tenant_id = $2;

-- name: CreateStockMovement :one
INSERT INTO stock_movements (
    tenant_id, item_id, movement_type, quantity, reason, reference, performed_by, notes
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING *;

-- name: UpdateItemStock :one
UPDATE inventory_items
SET 
    quantity_in_stock = $3,
    updated_at = NOW()
WHERE id = $1 AND tenant_id = $2
RETURNING *;

-- name: ListStockMovementsByItem :many
SELECT * FROM stock_movements
WHERE tenant_id = $1 AND item_id = $2
ORDER BY created_at DESC;

-- name: GetInventoryItemForUpdate :one
SELECT * FROM inventory_items
WHERE id = $1 AND tenant_id = $2
FOR UPDATE;

-- name: CreateIssuance :one
INSERT INTO inventory_issuances(
    tenant_id, item_id, quantity, issued_to_type,
    issued_to_id, issued_by, academic_session_id, notes
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING *;

-- name: ListIssuancesByItem :many
SELECT * FROM inventory_issuances
WHERE tenant_id = $1 AND item_id = $2
ORDER BY created_at DESC;

-- name: ListIssuancesByRecipient :many
SELECT * FROM inventory_issuances
WHERE tenant_id = $1 AND issued_to_type = $2 AND issued_to_id = $3
ORDER BY created_at DESC;