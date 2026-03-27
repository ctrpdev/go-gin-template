-- name: CreateUser :one
INSERT INTO users (email, password_hash, role)
VALUES ($1, $2, COALESCE($3, 'user'))
RETURNING id, email, role, verified, created_at, updated_at;

-- name: GetUserByEmail :one
SELECT id, email, password_hash, role, verified, created_at, updated_at
FROM users
WHERE email = $1 LIMIT 1;

-- name: VerifyUser :execrows
UPDATE users
SET verified = TRUE, updated_at = NOW()
WHERE id = $1;

-- name: GetUserByID :one
SELECT id, email, role, verified, created_at, updated_at
FROM users
WHERE id = $1 LIMIT 1;
