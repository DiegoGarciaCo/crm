-- name: CreateTag :one
INSERT INTO
    tags (
        name,
        description,
        user_id
    )
VALUES
    ($1, $2, $3)
RETURNING
    *;

-- name: DeleteTag :exec
DELETE FROM
    tags
WHERE
    id = $1
RETURNING
    *;

-- name: GetAllTags :many
SELECT
    *
FROM
    tags
WHERE
    user_id = $1;

-- name: AssignTagToContact :one
INSERT INTO
    contact_tags(tag_id, contact_id)
VALUES
    ($1, $2)
RETURNING
    *;

-- name: RemoveTagFromContact :exec
DELETE FROM
    contact_tags
WHERE
    tag_id = $1
    AND contact_id = $2;

-- name: AssignTagsToContact :exec
WITH input_tags AS (
    SELECT
        unnest($2::text []) AS tag_name
),
inserted_tags AS (
    INSERT INTO
        tags(name, user_id)
    SELECT
        tag_name,
        $1
    FROM
        input_tags ON conflict ON CONSTRAINT unique_user_tag DO nothing
    RETURNING
        id,
        name
),
all_tags AS (
    SELECT
        id,
        name
    FROM
        inserted_tags
    UNION
    SELECT
        id,
        name
    FROM
        tags t
    WHERE
        t.name IN (
            SELECT
                tag_name
            FROM
                input_tags
        )
)
INSERT INTO
    contact_tags(contact_id, tag_id)
SELECT
    $3,
    id
FROM
    all_tags ON conflict DO nothing;

-- name: TestBulkAssignTagsToContacts :exec
INSERT INTO
    contact_tags (contact_id, user_id, tag)
SELECT
    t.contact_id,
    t.user_id,
    t.tag
FROM
    jsonb_to_recordset($1::jsonb) AS t(
        contact_id uuid,
        user_id uuid,
        tag text
    );
