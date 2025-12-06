-- name: LandingPageEmails :one
INSERT INTO
    contacts (first_name, last_name, source, owner_id)
VALUES
    ($1, $2, $3, $4)
RETURNING
    *;
