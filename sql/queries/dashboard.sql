-- name: NewContactsThisMonth :one
SELECT
    count(*) AS new_contacts
FROM
    contacts
WHERE
    created_at >= date_trunc('month', current_date)
    AND owner_id = $1;

-- name: AppointmentsThisWeek :one
SELECT
    count(*) AS appointments_this_week
FROM
    appointments
WHERE
    scheduled_at >= date_trunc('week', current_date)
    AND scheduled_at < date_trunc('week', current_date) + INTERVAL '7 days'
    AND outcome = 'no-outcome'
    AND assigned_to_id = $1;

-- name: TasksDueTodayCount :one
SELECT
    count(*) AS tasks_due_today
FROM
    tasks
WHERE
    date = current_date
    AND STATUS = 'pending'
    AND assigned_to_id = $1;

-- name: Get5NewestContacts :many
SELECT
    c.id,
    c.first_name,
    c.last_name,
    c.birthdate,
    c.source,
    c.status,
    c.address,
    c.city,
    c.state,
    c.zip_code,
    c.lender,
    c.price_range,
    c.timeframe,
    c.owner_id,
    c.created_at,
    c.updated_at,
    c.last_contacted_at,
    -- Aggregate emails
    coalesce(
        jsonb_agg(
            DISTINCT jsonb_build_object(
                'id',
                e.id,
                'email_address',
                e.email_address,
                'type',
                e.type,
                'is_primary',
                e.is_primary
            )
        ) filter (
            WHERE
                e.id IS NOT NULL
        ),
        '[]'::jsonb
    ) AS emails,
    -- Aggregate phones
    coalesce(
        jsonb_agg(
            DISTINCT jsonb_build_object(
                'id',
                p.id,
                'phone_number',
                p.phone_number,
                'type',
                p.type,
                'is_primary',
                p.is_primary
            )
        ) filter (
            WHERE
                p.id IS NOT NULL
        ),
        '[]'::jsonb
    ) AS phone_numbers
FROM
    contacts c
    LEFT JOIN emails e ON e.contact_id = c.id
    LEFT JOIN phone_numbers p ON p.contact_id = c.id
WHERE
    c.owner_id = $1
GROUP BY
    c.id
ORDER BY
    c.created_at DESC
LIMIT
    5;

-- name: GetUpcomingAppointments :many
SELECT
    *
FROM
    appointments a
WHERE
    a.scheduled_at >= current_date
    AND a.status = 'scheduled'
    AND a.assigned_to_id = $1
ORDER BY
    a.scheduled_at ASC
LIMIT
    5;

-- name: ContactsCount :one
SELECT
    count(*) AS total_contacts
FROM
    contacts
WHERE
    owner_id = $1;

-- name: ContactsBySource :many
SELECT
    source,
    count(*) AS contact_count
FROM
    contacts
WHERE
    owner_id = $1
GROUP BY
    source
ORDER BY
    contact_count DESC;
