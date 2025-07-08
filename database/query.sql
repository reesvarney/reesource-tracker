-- name: GetLocation :one
SELECT
    id,
    name,
    description,
    parent_location_id
FROM
    locations
WHERE
    id = ?;

-- name: GetLocations :many
SELECT
    id,
    name,
    description,
    parent_location_id
FROM
    locations
ORDER BY
    name;

-- name: UpsertLocation :exec
INSERT INTO
    locations (id, name, description, parent_location_id)
VALUES
    (?, ?, ?, ?) ON CONFLICT (id) DO
UPDATE
SET
    name = EXCLUDED.name,
    description = EXCLUDED.description,
    parent_location_id = EXCLUDED.parent_location_id;

-- name: ListSamples :many
SELECT
    samples.*,
    COALESCE(
        (
            SELECT
                GROUP_CONCAT (sample_mods.name, ', ')
            FROM
                sample_mods
            WHERE
                sample_mods.sample_id = samples.id
                AND sample_mods.time_removed IS NULL
        ),
        ''
    ) AS current_mods_summary
FROM
    samples
ORDER BY
    time_registered;

-- name: ListSampleMods :many
SELECT
    *
FROM
    sample_mods
WHERE
    sample_mods.sample_id = ?
ORDER BY
    time_added;

-- name: AddSampleMod :exec
INSERT INTO
    sample_mods (id, sample_id, name, time_added, time_removed)
VALUES
    (?, ?, ?, ?, NULL);

-- name: RemoveSampleMod :exec
UPDATE sample_mods
SET
    time_removed = ?
WHERE
    id = ?;

-- name: DeleteProductByID :exec
DELETE FROM products
WHERE
    id = ?;

-- name: DeleteLocationByID :exec
DELETE FROM locations
WHERE
    id = ?;

-- name: GetProducts :many
SELECT
    *
FROM
    products
ORDER BY
    name;

-- name: UpdateOrCreateSample :one
INSERT INTO
    samples (
        id,
        location_id,
        product_id,
        time_registered,
        last_update,
        state
    )
VALUES
    (?, ?, ?, ?, ?, ?) ON CONFLICT (id) DO
UPDATE
SET
    location_id = EXCLUDED.location_id,
    product_id = EXCLUDED.product_id,
    last_update = EXCLUDED.last_update,
    state = EXCLUDED.state RETURNING *;

-- name: GetSampleById :one
SELECT
    *
FROM
    samples
WHERE
    id = ?;

-- name: ListProducts :many
SELECT
    id,
    name,
    parent_product_id
FROM
    products
ORDER BY
    name;

-- name: GetProductByID :one
SELECT
    id,
    name,
    parent_product_id
FROM
    products
WHERE
    id = ?;

-- name: UpsertProduct :exec
INSERT INTO
    products (id, name, parent_product_id)
VALUES
    (?, ?, ?) ON CONFLICT (id) DO
UPDATE
SET
    name = EXCLUDED.name,
    parent_product_id = EXCLUDED.parent_product_id;