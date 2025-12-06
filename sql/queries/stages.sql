-- name: CreateStage :one
INSERT INTO
    stages (
        name,
        description,
        client_type,
        order_index,
        owner_id
    )
VALUES
    ($1, $2, $3, $4, $5)
RETURNING
    *;

-- name: GetStagesByClientType :many
SELECT
    *
FROM
    stages
WHERE
    client_type = $1
    AND owner_id = $2
ORDER BY
    order_index ASC;

-- name: UpdateStage :one
UPDATE
    stages
SET
    name = $2,
    description = $3,
    client_type = $4,
    order_index = $5
WHERE
    id = $1
RETURNING
    *;

-- name: DeleteStage :exec
DELETE FROM
    stages
WHERE
    id = $1;

-- name: GetAllStages :many
SELECT
    *
FROM
    stages
WHERE
    owner_id = $1
ORDER BY
    client_type ASC,
    order_index ASC;

-- name: GetStageByID :one
SELECT
    *
FROM
    stages
WHERE
    id = $1;
