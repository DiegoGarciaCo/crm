-- name: GetOrganizationMembers :many
SELECT
    m.id,
    m."userId",
    u.name,
    m.role
FROM
    member m
    JOIN users u ON u.id = m."userId"
WHERE
    m."organizationId" = $1;
