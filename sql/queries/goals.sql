-- name: SetGoal :one
INSERT INTO
    goals (
        user_id,
        year,
        MONTH,
        income_goal,
        transaction_goal,
        estimated_average_sale_price,
        estimated_average_commission_rate
    )
VALUES
    (
        $1,
        $2,
        $3,
        $4,
        $5,
        $6,
        $7
    )
RETURNING
    *;

-- name: GetGoalByUserAndYear :one
SELECT
    *
FROM
    goals
WHERE
    user_id = $1
    AND year = $2;

-- name: UpdateGoal :one
UPDATE
    goals
SET
    income_goal = $2,
    transaction_goal = $3,
    estimated_average_sale_price = $4,
    estimated_average_commission_rate = $5
WHERE
    id = $1
RETURNING
    *;

-- name: DeleteGoal :one
DELETE FROM
    goals
WHERE
    id = $1
RETURNING
    *;
