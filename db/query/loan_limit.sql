-- name: GetLimit :one
SELECT `limit` FROM loan_limit
WHERE customer_id = ? AND tenor = ?;

-- name: GenerateLimit :exec
INSERT INTO loan_limit (
    id,
    customer_id,
    tenor,
    `limit`
) VALUES (
    ?, ?, ?, ?
);

-- name: GetCustomerLimit :many
SELECT tenor, `limit` FROM loan_limit
where customer_id  = ?
order by tenor;

-- name: ReduceLimit :exec
UPDATE loan_limit
SET `limit` = `limit` - ?
WHERE customer_id = ? AND `limit` >= ?;

-- name: IncreaseLimit :exec
UPDATE loan_limit
SET `limit` = `limit` + ?
WHERE customer_id = ?