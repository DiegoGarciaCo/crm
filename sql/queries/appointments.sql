-- name: CreateAppointment :one
INSERT INTO
    appointments (
        assigned_to_id,
        contact_id,
        title,
        scheduled_at,
        location,
        TYPE,
        outcome,
        note
    )
VALUES
    (
        $1,
        $2,
        $3,
        $4,
        $5,
        $6,
        $7,
        $8
    )
RETURNING
    *;

-- name: GetAppointmentById :one
SELECT
    *
FROM
    appointments
WHERE
    id = $1;

-- name: UpdateAppointment :one
UPDATE
    appointments
SET
    contact_id = $2,
    title = $3,
    scheduled_at = $4,
    location = $5,
    TYPE = $6,
    outcome = $7,
    note = $8
WHERE
    id = $1
RETURNING
    *;

-- name: DeleteAppointment :exec
DELETE FROM
    appointments
WHERE
    id = $1;

-- name: ListTodaysAppointments :many
SELECT
    *
FROM
    appointments
WHERE
    scheduled_at::date = current_date
    AND assigned_to_id = $1
ORDER BY
    scheduled_at ASC;

-- name: ListAppointments :many
SELECT
    *
FROM
    appointments
WHERE
    assigned_to_id = $1
ORDER BY
    scheduled_at ASC;

-- name: ListUpcomingAppointments :many
SELECT
    *
FROM
    appointments
WHERE
    scheduled_at > NOW()
    AND assigned_to_id = $1
ORDER BY
    scheduled_at ASC;

-- name: ListPastAppointments :many
SELECT
    *
FROM
    appointments
WHERE
    scheduled_at < NOW()
    AND assigned_to_id = $1
ORDER BY
    scheduled_at DESC;

-- name: ListAppointmentsByContactId :many
SELECT
    *
FROM
    appointments
WHERE
    contact_id = $1
ORDER BY
    scheduled_at ASC;
