-- name: GetAPIKeyByHash :one
SELECT
    name,
    "userId",
    permissions,
    "expiresAt",
    enabled,
    metadata
FROM
    apikey
WHERE
    KEY = $1;
