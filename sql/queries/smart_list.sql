-- name: CreateSmartList :one
INSERT INTO
    smart_lists (
        name,
        description,
        user_id
    )
VALUES
    ($1, $2, $3)
RETURNING
    *;

-- name: UpdateSmartList :one
UPDATE
    smart_lists
SET
    name = $2,
    description = $3,
    filter_criteria = $4,
    updated_at = CURRENT_TIMESTAMP
WHERE
    id = $1
RETURNING
    *;

-- name: GetAllSmartLists :many
SELECT
    *
FROM
    smart_lists
WHERE
    user_id = $1;

-- name: SetSmartListFilterCriteria :one
UPDATE
    smart_lists
SET
    filter_criteria = $2,
    updated_at = CURRENT_TIMESTAMP
WHERE
    id = $1
RETURNING
    *;

-- name: GetAllSmartListsWithCounts :many
SELECT
    s.*,
    count(c.id) AS total_contacts,
    count(c.id) filter (
        WHERE
            -- first_name
            (
                coalesce(s.filter_criteria, '{}'::jsonb) ->> 'first_name' IS NULL
                OR c.first_name ilike '%' || (
                    coalesce(s.filter_criteria, '{}'::jsonb) ->> 'first_name'
                ) || '%'
            )
            -- last_name
            AND (
                coalesce(s.filter_criteria, '{}'::jsonb) ->> 'last_name' IS NULL
                OR c.last_name ilike '%' || (
                    coalesce(s.filter_criteria, '{}'::jsonb) ->> 'last_name'
                ) || '%'
            )
            -- birthdate
            AND (
                coalesce(s.filter_criteria, '{}'::jsonb) ->> 'birthdate' IS NULL
                OR c.birthdate = (
                    coalesce(s.filter_criteria, '{}'::jsonb) ->> 'birthdate'
                )::date
            )
            -- source
            AND (
                coalesce(s.filter_criteria, '{}'::jsonb) ->> 'source' IS NULL
                OR c.source ilike '%' || (
                    coalesce(s.filter_criteria, '{}'::jsonb) ->> 'source'
                ) || '%'
            )
            -- status
            AND (
                coalesce(s.filter_criteria, '{}'::jsonb) ->> 'status' IS NULL
                OR c.status = (
                    coalesce(s.filter_criteria, '{}'::jsonb) ->> 'status'
                )
            )
            -- address
            AND (
                coalesce(s.filter_criteria, '{}'::jsonb) ->> 'address' IS NULL
                OR c.address ilike '%' || (
                    coalesce(s.filter_criteria, '{}'::jsonb) ->> 'address'
                ) || '%'
            )
            -- city
            AND (
                coalesce(s.filter_criteria, '{}'::jsonb) ->> 'city' IS NULL
                OR c.city ilike '%' || (
                    coalesce(s.filter_criteria, '{}'::jsonb) ->> 'city'
                ) || '%'
            )
            -- state
            AND (
                coalesce(s.filter_criteria, '{}'::jsonb) ->> 'state' IS NULL
                OR c.state ilike '%' || (
                    coalesce(s.filter_criteria, '{}'::jsonb) ->> 'state'
                ) || '%'
            )
            -- zip_code
            AND (
                coalesce(s.filter_criteria, '{}'::jsonb) ->> 'zip_code' IS NULL
                OR c.zip_code = (
                    coalesce(s.filter_criteria, '{}'::jsonb) ->> 'zip_code'
                )
            )
            -- lender
            AND (
                coalesce(s.filter_criteria, '{}'::jsonb) ->> 'lender' IS NULL
                OR c.lender ilike '%' || (
                    coalesce(s.filter_criteria, '{}'::jsonb) ->> 'lender'
                ) || '%'
            )
            -- price_range
            AND (
                coalesce(s.filter_criteria, '{}'::jsonb) ->> 'price_range' IS NULL
                OR c.price_range = (
                    coalesce(s.filter_criteria, '{}'::jsonb) ->> 'price_range'
                )
            )
            -- timeframe
            AND (
                coalesce(s.filter_criteria, '{}'::jsonb) ->> 'timeframe' IS NULL
                OR c.timeframe = (
                    coalesce(s.filter_criteria, '{}'::jsonb) ->> 'timeframe'
                )
            )
            -- owner_id
            AND (
                coalesce(s.filter_criteria, '{}'::jsonb) ->> 'owner_id' IS NULL
                OR c.owner_id = (
                    coalesce(s.filter_criteria, '{}'::jsonb) ->> 'owner_id'
                )::uuid
            )
            -- contact has all required tags
            AND (
                coalesce(s.filter_criteria, '{}'::jsonb) -> 'tag_ids' IS NULL
                OR (
                    SELECT
                        count(DISTINCT required_tag)
                    FROM
                        jsonb_array_elements_text(
                            coalesce(s.filter_criteria, '{}'::jsonb) -> 'tag_ids'
                        ) AS required(required_tag)
                        JOIN contact_tags ct2 ON ct2.tag_id = required.required_tag::uuid
                        AND ct2.contact_id = c.id
                ) = jsonb_array_length(
                    coalesce(s.filter_criteria, '{}'::jsonb) -> 'tag_ids'
                )
            )
            -- last_contacted_days
            AND (
                coalesce(s.filter_criteria, '{}'::jsonb) ->> 'last_contacted_days' IS NULL
                OR c.last_contacted_at IS NULL
                OR c.last_contacted_at <= NOW() - (
                    (
                        coalesce(s.filter_criteria, '{}'::jsonb) ->> 'last_contacted_days'
                    ) || ' days'
                )::INTERVAL
            )
    ) AS contact_count
FROM
    smart_lists s
    LEFT JOIN contacts c ON (
        c.owner_id = s.user_id
        OR EXISTS (
            SELECT
                1
            FROM
                collaborators col
            WHERE
                col.contact_id = c.id
                AND col.user_id = s.user_id
        )
    )
WHERE
    s.user_id = $1
GROUP BY
    s.id
ORDER BY
    s.list_index ASC;

-- name: ReorderSmartLists :exec
UPDATE
    smart_lists
SET
    list_index = data.new_index,
    updated_at = CURRENT_TIMESTAMP
FROM
    (
        SELECT
            unnest($1::uuid []) AS id,
            unnest($2::int []) AS new_index
    ) AS data
WHERE
    smart_lists.id = data.id;

-- name: BumpSmartListIndexes :exec
UPDATE
    smart_lists
SET
    list_index = list_index + 1000
WHERE
    user_id = $1;
