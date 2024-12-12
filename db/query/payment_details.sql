-- name: CreatePayment :exec
INSERT INTO `payment_details` (
    id,
    transaction_id,
    amount,
    due_date,
    is_paid
) VALUES (
    ?, ?, ?, ?, ?
)
