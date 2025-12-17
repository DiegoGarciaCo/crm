-- name: CreateNotification :one
INSERT INTO
    notifications (
        user_id,
        TYPE,
        message,
        contact_id,
        appointment_id,
        task_id
    )
VALUES
    ($1, $2, $3, $4, $5, $6)
RETURNING
    *;

-- name: GetNotificationsByUserID :many
SELECT
    *
FROM
    notifications
WHERE
    user_id = $1
ORDER BY
    created_at DESC
LIMIT
    $2 OFFSET $3;

-- name: MarkNotificationAsRead :exec
UPDATE
    notifications
SET
    READ = TRUE
WHERE
    id = $1;

-- name: DeleteNotification :exec
DELETE FROM
    notifications
WHERE
    id = $1;

-- name: MarkAllNotificationsAsRead :exec
UPDATE
    notifications
SET
    READ = TRUE
WHERE
    user_id = $1;
