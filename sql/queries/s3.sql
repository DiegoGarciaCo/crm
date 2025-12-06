-- name: GetUserImage :one
SELECT
    image
FROM
    users
WHERE
    id = $1;

-- name: UpdateUserImage :exec
UPDATE
    users
SET
    image = $2
WHERE
    id = $1;
