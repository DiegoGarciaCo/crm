-- name: LogContact :one
INSERT INTO
    contact_logs (
        contact_id,
        contact_method,
        created_by,
        note
    )
VALUES
    ($1, $2, $3, $4)
RETURNING
    *;

-- name: GetContactLogsByContactID :many
SELECT
    *
FROM
    contact_logs
WHERE
    contact_id = $1
ORDER BY
    created_at DESC;
