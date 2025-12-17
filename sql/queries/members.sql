-- name: GetOrganizationMembers :many
SELECT
    id,
    "userId",
    role
FROM
    member
WHERE
    "organizationId" = $1;
