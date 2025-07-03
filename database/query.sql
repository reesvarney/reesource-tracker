-- name: ListSamples :many
SELECT 
    samples.*, 
    COALESCE((
        SELECT GROUP_CONCAT(sample_mods.name, ', ')
        FROM sample_mods 
        WHERE sample_mods.sample_id = samples.id AND sample_mods.time_removed IS NULL
    ), '') AS current_mods_summary
FROM samples
ORDER BY time_registered;

-- name: ListSampleMods :many
SELECT * FROM sample_mods
WHERE sample_mods.sample_id = ?
ORDER BY time_added;

-- name: ListLocations :many
SELECT * FROM locations
ORDER BY name;

-- name: GetProducts :many
WITH RECURSIVE item_tree(id, name, full_name) AS (
    SELECT id, name, name as full_name FROM products
    UNION ALL
    SELECT p.id, p.name, item_tree.full_name || ' > ' || p.name
    FROM products p
    JOIN item_tree ON p.parent_product_id = item_tree.id
)
SELECT
    item_tree.id AS item_id,
    item_tree.name AS item_name,
    item_tree.full_name AS combined_name
FROM item_tree;

-- name: UpdateOrCreateSample :one
INSERT INTO samples (id, location_id, product_id, time_registered, last_update, state)
VALUES (?, ?, ?, ?, ?, ?)
ON CONFLICT(id) DO UPDATE SET
    location_id = EXCLUDED.location_id,
    product_id = EXCLUDED.product_id,
    last_update = EXCLUDED.last_update,
    state = EXCLUDED.state
RETURNING *;

-- name: GetSampleById :one
SELECT * FROM samples WHERE id = ?;

-- name: ListProducts :many
SELECT id, name, parent_product_id FROM products ORDER BY name;

-- name: GetProductByID :one
SELECT id, name, parent_product_id FROM products WHERE id = ?;

-- name: UpsertProduct :exec
INSERT INTO products (id, name, parent_product_id)
VALUES (?, ?, ?)
ON CONFLICT(id) DO UPDATE SET
    name = EXCLUDED.name,
    parent_product_id = EXCLUDED.parent_product_id;