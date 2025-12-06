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
    *
FROM
    contact_notes
WHERE
    contact_id = $1
ORDER BY
    created_at DESC;
