-- name: EnterPhoneNumber :exec
INSERT INTO
    phone_numbers (contact_id, phone_number, TYPE, is_primary)
VALUES
    (@contact_id, @phone_number, @TYPE, @is_primary);

-- name: DeletePhoneNumber :exec
DELETE FROM
    phone_numbers
WHERE
    id = @id;

-- name: UpdatePhoneNumber :exec
UPDATE
    phone_numbers
SET
    phone_number = @phone_number,
    TYPE = @TYPE,
    is_primary = @is_primary,
    updated_at = CURRENT_TIMESTAMP
WHERE
    id = @id;

-- name: GetPhoneNumbersByContactID :many
SELECT
    *
FROM
    phone_numbers
WHERE
    contact_id = $1
ORDER BY
    is_primary DESC,
    created_at DESC;

-- name: GetPhoneNumberByID :one
SELECT
    *
FROM
    phone_numbers
WHERE
    id = $1;

-- name: BulkEnterPhoneNumbers :exec
INSERT INTO
    phone_numbers (contact_id, phone_number, TYPE, is_primary)
SELECT
    unnest(@contact_ids::uuid []),
    unnest(@phone_numbers::text []),
    unnest(@types::text []),
    unnest(@is_primary::boolean []);

-- name: TestBulkInsertPhoneNumbers :exec
INSERT INTO
    phone_numbers (contact_id, number, TYPE, is_primary)
SELECT
    p.contact_id,
    p.number,
    p.type,
    p.is_primary
FROM
    jsonb_to_recordset($1::jsonb) AS p(
        contact_id uuid,
        number text,
        TYPE text,
        is_primary boolean
    );
