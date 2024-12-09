-- name: GetLimit :one
SELECT `limit` FROM loan_limit
WHERE customer_id = ? AND tenor = ?;