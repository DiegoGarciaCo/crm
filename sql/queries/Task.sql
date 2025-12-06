-- name: CreateTask :one
INSERT INTO
    tasks (
        contact_id,
        assigned_to_id,
        title,
        TYPE,
        date,
        STATUS,
        priority,
        note
    )
VALUES
    ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING
    *;

-- name: GetTasksByContactID :many
SELECT
    *
FROM
    tasks
WHERE
    contact_id = $1
ORDER BY
    date DESC;

-- name: GetTaskByAssignedToID :many
SELECT
    *
FROM
    tasks
WHERE
    assigned_to_id = $1
ORDER BY
    date DESC;

-- name: GetTaskDueToday :many
SELECT
    *
FROM
    tasks
WHERE
    date::date = current_date
    AND assigned_to_id = $1
    AND STATUS != 'completed'
ORDER BY
    date DESC;

-- name: GetOverdueTasks :many
SELECT
    *
FROM
    tasks
WHERE
    date::date < current_date
    AND STATUS != 'completed'
    AND assigned_to_id = $1
ORDER BY
    date DESC;

-- name: UpdateTaskStatus :one
UPDATE
    tasks
SET
    STATUS = $2,
    updated_at = CURRENT_TIMESTAMP
WHERE
    id = $1
RETURNING
    *;

-- name: DeleteTask :exec
DELETE FROM
    tasks
WHERE
    id = $1;

-- name: GetTaskByID :one
SELECT
    *
FROM
    tasks
WHERE
    id = $1;

-- name: UpdateTask :one
UPDATE
    tasks
SET
    contact_id = $2,
    assigned_to_id = $3,
    title = $4,
    TYPE = $5,
    date = $6,
    STATUS = $7,
    priority = $8,
    note = $9,
    updated_at = CURRENT_TIMESTAMP
WHERE
    id = $1
RETURNING
    *;
