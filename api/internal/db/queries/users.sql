-- name: GetUserByEmail :one
SELECT id, email, name, password_hash, role
FROM users
WHERE email = $1;

-- name: GetUserByID :one
SELECT id, email, name, password_hash, role
FROM users
WHERE id = $1;

-- name: CreateUser :one
INSERT INTO users (email, name, password_hash, role)
VALUES ($1, $2, $3, 'user')
RETURNING *;
