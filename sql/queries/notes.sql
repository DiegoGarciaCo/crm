-- name: CreateNote :one
INSERT INTO
    contact_notes (
        contact_id,
        note,
        created_by
    )
VALUES
    ($1, $2, $3)
RETURNING
    *;

-- name: GetNotesByContactID :many
SELECT
    cn.*,
    u.name AS created_by_name
FROM
    contact_notes cn
    JOIN users u ON u.id = cn.created_by
WHERE
    cn.contact_id = $1
ORDER BY
    cn.created_at DESC;
