-- name: AddCollaborator :exec
INSERT INTO
    collaborators (contact_id, user_id, role)
VALUES
    ($1, $2, $3);

-- name: RemoveCollaborator :exec
DELETE FROM
    collaborators
WHERE
    contact_id = $1
    AND user_id = $2;

-- name: ListCollaborators :many
SELECT
    u.id,
    u.name,
    c.role
FROM
    collaborators c
    JOIN users u ON c.user_id = u.id
WHERE
    c.contact_id = $1;
