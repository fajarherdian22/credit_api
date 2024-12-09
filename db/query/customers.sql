-- name: CreateCustomers :exec
INSERT INTO customers (
  id,
  nik,
  hashed_password,
  email,
  full_name,
  legal_name,
  tempat_lahir,
  tanggal_lahir,
  gaji,
  photo_ktp,
  foto_selfie
) VALUES (
  ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
);

-- name: GetCustomers :one
SELECT * FROM customers
WHERE email = ? LIMIT 1;
