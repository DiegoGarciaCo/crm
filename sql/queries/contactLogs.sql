-- name: LogContact :one
WITH inserted AS (
    INSERT INTO
        contact_logs(contact_id, contact_method, created_by, note)
    VALUES
        ($1, $2, $3, $4)
    RETURNING
        *
),
updated AS (
    UPDATE
        contacts
    SET
        last_contacted_at = NOW()
    WHERE
        id = $1
    RETURNING
        id
)
SELECT
    *
FROM
    inserted;

-- name: GetContactLogsByContactID :many
SELECT
    cl.*,
    u.name AS created_by_name
FROM
    contact_logs cl
    JOIN users u ON u.id = cl.created_by
WHERE
    cl.contact_id = $1
ORDER BY
    cl.created_at DESC;
