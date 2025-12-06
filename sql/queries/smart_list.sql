-- name: CreateSmartList :one
INSERT INTO
    smart_lists (
        name,
        description,
        user_id
    )
VALUES
    ($1, $2, $3)
RETURNING
    *;

-- name: UpdateSmartList :one
UPDATE
    smart_lists
SET
    name = $2,
    description = $3,
    filter_criteria = $4,
    updated_at = CURRENT_TIMESTAMP
WHERE
    id = $1
RETURNING
    *;

-- name: GetAllSmartLists :many
SELECT
    *
FROM
    smart_lists
WHERE
    user_id = $1;

-- name: SetSmartListFilterCriteria :one
UPDATE
    smart_lists
SET
    filter_criteria = $2,
    updated_at = CURRENT_TIMESTAMP
WHERE
    id = $1
RETURNING
    *;
