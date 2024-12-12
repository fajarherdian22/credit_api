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

-- name: UpdatePayment :exec
UPDATE payment_details
SET is_paid = true
WHERE id = ? AND is_paid = false;