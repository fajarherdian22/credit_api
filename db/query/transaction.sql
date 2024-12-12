-- name: CreateTransaction :exec
INSERT INTO `transaction` (
  id,
  customer_id,
  product_name,
  price,
  bunga,
  jumlah_cicilan,
  tenor,
  admin_fee,
  created_at
) VALUES (
  ?, ?, ?, ?, ?, ?, ?, ?, ?
);

-- name: GetTransaction :one
SELECT * FROM `transaction`
WHERE id = ?;

-- name: ListTransaction :many
SELECT * FROM `transaction`
WHERE customer_id = ?
ORDER BY created_at;