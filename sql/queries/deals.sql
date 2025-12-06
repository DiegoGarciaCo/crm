-- name: CreateDeal :one
INSERT INTO
    deals (
        contact_id,
        assigned_to_id,
        title,
        price,
        closing_date,
        earnest_money_due_date,
        mutual_acceptance_date,
        inspection_date,
        appraisal_date,
        final_walkthrough_date,
        possession_date,
        closed_date,
        commission,
        commission_split,
        property_address,
        property_city,
        property_state,
        property_zip_code,
        description,
        stage_id
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
        $13,
        $14,
        $15,
        $16,
        $17,
        $18,
        $19,
        $20
    )
RETURNING
    *;

-- name: GetDealById :one
SELECT
    *
FROM
    deals
WHERE
    id = $1;

-- name: UpdateDeal :one
UPDATE
    deals
SET
    contact_id = $2,
    assigned_to_id = $3,
    title = $4,
    price = $5,
    closing_date = $6,
    earnest_money_due_date = $7,
    mutual_acceptance_date = $8,
    inspection_date = $9,
    appraisal_date = $10,
    final_walkthrough_date = $11,
    possession_date = $12,
    commission = $13,
    commission_split = $14,
    property_address = $15,
    property_city = $16,
    property_state = $17,
    property_zip_code = $18,
    description = $19,
    stage_id = $20,
    closed_date = $21
WHERE
    id = $1
RETURNING
    *;

-- name: DeleteDeal :exec
DELETE FROM
    deals
WHERE
    id = $1
RETURNING
    *;

-- name: ListDeals :many
SELECT
    *
FROM
    deals
WHERE
    assigned_to_id = $3
ORDER BY
    created_at DESC
LIMIT
    $1 OFFSET $2;

-- name: CountDeals :one
SELECT
    count(*)
FROM
    deals
WHERE
    assigned_to_id = $1;

-- name: ListDealsByStage :many
SELECT
    *
FROM
    deals
WHERE
    stage_id = $1
    AND assigned_to_id = $4
ORDER BY
    created_at DESC
LIMIT
    $2 OFFSET $3;

-- name: CountDealsByStage :one
SELECT
    count(*)
FROM
    deals
WHERE
    stage_id = $1
    AND assigned_to_id = $2;

-- name: ListDealsByContactID :many
SELECT
    *
FROM
    deals
WHERE
    contact_id = $1
    AND assigned_to_id = $2;
