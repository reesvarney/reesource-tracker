-- name: GetLocation :one
SELECT *
FROM
    locations
WHERE
    id = ?;

-- name: GetLocations :many
SELECT *
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
    ) AS current_mods_summary,
    users.name AS owner_name
FROM
    samples
LEFT JOIN users ON samples.owner_id = users.id
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
        state,
        owner_id,
        product_issue
    )
VALUES
    (?, ?, ?, ?, ?, ?, ?, ?) ON CONFLICT (id) DO
UPDATE
SET
    location_id = EXCLUDED.location_id,
    product_id = EXCLUDED.product_id,
    last_update = EXCLUDED.last_update,
    owner_id = EXCLUDED.owner_id,
    product_issue = EXCLUDED.product_issue,
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
    *
FROM
    products
ORDER BY
    name;

-- name: GetProductByID :one
SELECT *
FROM
    products
WHERE
    id = ?;

-- name: UpsertProduct :exec
INSERT INTO
    products (id, name, parent_product_id, part_number)
VALUES
    (?, ?, ?, ?) ON CONFLICT (id) DO
UPDATE
SET
    name = EXCLUDED.name,
    parent_product_id = EXCLUDED.parent_product_id,
    part_number = EXCLUDED.part_number;


-- name: GetUserByID :one
SELECT *
FROM
    users
WHERE
    id = ?;

-- name: GetUsers :many
SELECT *
FROM
    users
ORDER BY
    name;

-- name: UpsertUser :exec
INSERT INTO
    users (id, name)
VALUES
    (?, ?) ON CONFLICT (id) DO
UPDATE
SET
    name = EXCLUDED.name;

-- name: DeleteUserByID :exec
DELETE FROM users
WHERE
    id = ?;