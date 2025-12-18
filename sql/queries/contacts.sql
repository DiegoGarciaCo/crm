-- name: BulkInsertContacts :many
INSERT INTO
    contacts (
        first_name,
        last_name,
        birthdate,
        source,
        STATUS,
        address,
        city,
        state,
        zip_code,
        lender,
        price_range,
        timeframe,
        owner_id
    )
SELECT
    unnest(@first_names::text []),
    unnest(@last_names::text []),
    unnest(@birthdates::date []),
    unnest(@sources::text []),
    unnest(@statuses::text []),
    unnest(@addresses::text []),
    unnest(@cities::text []),
    unnest(@states::text []),
    unnest(@zip_codes::text []),
    unnest(@lenders::text []),
    unnest(@price_ranges::text []),
    unnest(@timeframes::text []),
    unnest(@owner_ids::uuid [])
RETURNING
    *;

-- name: CreateContact :one
INSERT INTO
    contacts (
        first_name,
        last_name,
        birthdate,
        source,
        STATUS,
        address,
        city,
        state,
        zip_code,
        lender,
        price_range,
        timeframe,
        owner_id
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
        $8,
        $9,
        $10,
        $11,
        $12,
        $13
    )
RETURNING
    *;

-- name: GetContactWithDetails :one
SELECT
    c.*,
    coalesce(
        (
            SELECT
                json_agg(e.*)::text
            FROM
                emails e
            WHERE
                e.contact_id = c.id
        ),
        '[]'
    ) AS emails,
    coalesce(
        (
            SELECT
                json_agg(p.*)::text
            FROM
                phone_numbers p
            WHERE
                p.contact_id = c.id
        ),
        '[]'
    ) AS phone_numbers,
    coalesce(
        (
            SELECT
                json_agg(t.*)::text
            FROM
                tags t
                JOIN contact_tags ct ON ct.tag_id = t.id
            WHERE
                ct.contact_id = c.id
        ),
        '[]'
    ) AS tags,
    coalesce(
        (
            SELECT
                json_agg(
                    json_build_object('id', u.id, 'name', u.name, 'role', cb.role)
                )::text
            FROM
                collaborators cb
                JOIN users u ON cb.user_id = u.id
            WHERE
                cb.contact_id = c.id
        ),
        '[]'
    ) AS collaborators
FROM
    contacts c
WHERE
    c.id = $1;

-- name: GetAllContacts :many
SELECT
    c.*,
    coalesce(
        (
            SELECT
                json_agg(p.*)::text
            FROM
                phone_numbers p
            WHERE
                p.contact_id = c.id
        ),
        '[]'
    ) AS phone_numbers,
    count(*) over () AS total_count
FROM
    contacts c
WHERE
    -- user owns the contact
    c.owner_id = $3
    -- OR user is a collaborator on the contact
    OR EXISTS (
        SELECT
            1
        FROM
            collaborators col
        WHERE
            col.contact_id = c.id
            AND col.user_id = $3
    )
ORDER BY
    c.created_at DESC
LIMIT
    $1 OFFSET $2;

-- name: CreateContactWithDetails :one
WITH new_contact AS (
    INSERT INTO
        contacts(
            first_name,
            last_name,
            birthdate,
            source,
            STATUS,
            address,
            city,
            state,
            zip_code,
            lender,
            price_range,
            timeframe,
            owner_id
        )
    VALUES
        (
            @first_name,
            @last_name,
            @birthdate,
            @source,
            @status,
            @address,
            @city,
            @state,
            @zip_code,
            @lender,
            @price_range,
            @timeframe,
            @owner_id
        )
    RETURNING
        id
),
insert_phones AS (
    INSERT INTO
        phone_numbers(contact_id, phone_number, TYPE, is_primary)
    SELECT
        id,
        phone_number,
        TYPE,
        is_primary
    FROM
        new_contact,
        jsonb_to_recordset(@phones::jsonb) AS p(
            phone_number text,
            TYPE text,
            is_primary boolean
        )
),
insert_emails AS (
    INSERT INTO
        emails(contact_id, email_address, TYPE, is_primary)
    SELECT
        id,
        email_address,
        TYPE,
        is_primary
    FROM
        new_contact,
        jsonb_to_recordset(@emails::jsonb) AS e(
            email_address text,
            TYPE text,
            is_primary boolean
        )
)
SELECT
    id
FROM
    new_contact;

-- name: SearchContacts :many
SELECT
    *
FROM
    contacts
WHERE
    owner_id = $1
    AND (
        first_name ilike $2
        OR last_name ilike $2
        OR concat(first_name, ' ', last_name) ilike $2
        OR address ilike $2
        OR city ilike $2
        OR state ilike $2
        OR lender ilike $2
        OR source ilike $2
    )
ORDER BY
    last_name,
    first_name
LIMIT
    50;

-- name: GetContactsBySmartList :many
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
    c.last_contacted_at,
    c.created_at,
    c.updated_at,
    count(*) over () AS total_count
FROM
    contacts c
    LEFT JOIN contact_tags ct ON ct.contact_id = c.id
    LEFT JOIN tags t ON t.id = ct.tag_id
    JOIN smart_lists s ON s.id = $1
WHERE
    (
        c.owner_id = $4
        OR EXISTS (
            SELECT
                1
            FROM
                collaborators col
            WHERE
                col.contact_id = c.id
                AND col.user_id = $4
        )
    )
    -- 🔽 ALL YOUR EXISTING FILTERS 🔽
    AND (
        coalesce(s.filter_criteria, '{}'::jsonb) ->> 'first_name' IS NULL
        OR c.first_name ilike '%' || (
            coalesce(s.filter_criteria, '{}'::jsonb) ->> 'first_name'
        ) || '%'
    )
    AND (
        coalesce(s.filter_criteria, '{}'::jsonb) ->> 'last_name' IS NULL
        OR c.last_name ilike '%' || (
            coalesce(s.filter_criteria, '{}'::jsonb) ->> 'last_name'
        ) || '%'
    )
    AND (
        coalesce(s.filter_criteria, '{}'::jsonb) ->> 'birthdate' IS NULL
        OR c.birthdate = (
            coalesce(s.filter_criteria, '{}'::jsonb) ->> 'birthdate'
        )::date
    )
    AND (
        coalesce(s.filter_criteria, '{}'::jsonb) ->> 'source' IS NULL
        OR c.source ilike '%' || (
            coalesce(s.filter_criteria, '{}'::jsonb) ->> 'source'
        ) || '%'
    )
    AND (
        coalesce(s.filter_criteria, '{}'::jsonb) ->> 'status' IS NULL
        OR c.status = (
            coalesce(s.filter_criteria, '{}'::jsonb) ->> 'status'
        )
    )
    AND (
        coalesce(s.filter_criteria, '{}'::jsonb) ->> 'address' IS NULL
        OR c.address ilike '%' || (
            coalesce(s.filter_criteria, '{}'::jsonb) ->> 'address'
        ) || '%'
    )
    AND (
        coalesce(s.filter_criteria, '{}'::jsonb) ->> 'city' IS NULL
        OR c.city ilike '%' || (
            coalesce(s.filter_criteria, '{}'::jsonb) ->> 'city'
        ) || '%'
    )
    AND (
        coalesce(s.filter_criteria, '{}'::jsonb) ->> 'state' IS NULL
        OR c.state ilike '%' || (
            coalesce(s.filter_criteria, '{}'::jsonb) ->> 'state'
        ) || '%'
    )
    AND (
        coalesce(s.filter_criteria, '{}'::jsonb) ->> 'zip_code' IS NULL
        OR c.zip_code = (
            coalesce(s.filter_criteria, '{}'::jsonb) ->> 'zip_code'
        )
    )
    AND (
        coalesce(s.filter_criteria, '{}'::jsonb) ->> 'lender' IS NULL
        OR c.lender ilike '%' || (
            coalesce(s.filter_criteria, '{}'::jsonb) ->> 'lender'
        ) || '%'
    )
    AND (
        coalesce(s.filter_criteria, '{}'::jsonb) ->> 'price_range' IS NULL
        OR c.price_range = (
            coalesce(s.filter_criteria, '{}'::jsonb) ->> 'price_range'
        )
    )
    AND (
        coalesce(s.filter_criteria, '{}'::jsonb) ->> 'timeframe' IS NULL
        OR c.timeframe = (
            coalesce(s.filter_criteria, '{}'::jsonb) ->> 'timeframe'
        )
    )
    AND (
        coalesce(s.filter_criteria, '{}'::jsonb) ->> 'owner_id' IS NULL
        OR c.owner_id = (
            coalesce(s.filter_criteria, '{}'::jsonb) ->> 'owner_id'
        )::uuid
    )
    AND (
        coalesce(s.filter_criteria, '{}'::jsonb) ->> 'tag_id' IS NULL
        OR t.id = (
            coalesce(s.filter_criteria, '{}'::jsonb) ->> 'tag_id'
        )::uuid
    )
    AND (
        coalesce(s.filter_criteria, '{}'::jsonb) ->> 'last_contacted_days' IS NULL
        OR c.last_contacted_at <= NOW() - (
            (
                coalesce(s.filter_criteria, '{}'::jsonb) ->> 'last_contacted_days'
            ) || ' days'
        )::INTERVAL
    )
ORDER BY
    c.last_contacted_at ASC nulls FIRST,
    c.created_at DESC
LIMIT
    $2 OFFSET $3;

-- name: TestBulkInsertContacts :many
INSERT INTO
    contacts (
        first_name,
        last_name,
        birthdate,
        source,
        STATUS,
        address,
        city,
        state,
        zip_code,
        lender,
        price_range,
        timeframe,
        owner_id
    )
SELECT
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
    c.owner_id
FROM
    jsonb_to_recordset($1::jsonb) AS c(
        first_name text,
        last_name text,
        birthdate date,
        source text,
        STATUS text,
        address text,
        city text,
        state text,
        zip_code text,
        lender text,
        price_range text,
        timeframe text,
        owner_id uuid
    )
RETURNING
    id;

-- name: UpdateContact :one
UPDATE
    contacts
SET
    first_name = $2,
    last_name = $3,
    birthdate = $4,
    source = $5,
    STATUS = $6,
    address = $7,
    city = $8,
    state = $9,
    zip_code = $10,
    lender = $11,
    price_range = $12,
    timeframe = $13,
    updated_at = NOW()
WHERE
    id = $1
RETURNING
    *;
