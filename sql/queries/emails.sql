-- name: EnterEmail :exec
INSERT INTO
    emails (contact_id, email_address, TYPE, is_primary)
VALUES
    ($1, $2, $3, $4);

-- name: BulkEnterEmails :exec
INSERT INTO
    emails (contact_id, email, TYPE, is_primary)
SELECT
    unnest(@contact_ids::uuid []),
    unnest(@emails::text []),
    unnest(@types::text []),
    unnest(@is_primary::boolean []);

-- name: VerifyEmail :exec
UPDATE
    emails
SET
    is_verified = TRUE,
    is_subscribed = TRUE
WHERE
    email_address = $1
RETURNING
    *;

-- name: TestBulkInsertEmails :exec
INSERT INTO
    emails (contact_id, email_address, TYPE, is_primary)
SELECT
    e.contact_id,
    e.email_address,
    e.type,
    e.is_primary
FROM
    jsonb_to_recordset($1::jsonb) AS e(
        contact_id uuid,
        email_address text,
        TYPE text,
        is_primary boolean
    );
