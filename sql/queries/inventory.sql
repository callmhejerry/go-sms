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