-- name: CountUsers :one
SELECT COUNT(*) FROM users;

-- name: CountUsersByEmail :one
SELECT COUNT(*) FROM users
WHERE email = ?1;

-- name: CountUsersByUsername :one
SELECT COUNT(*) FROM users
WHERE username = ?1;

-- name: CreateUser :one
INSERT INTO users (
    username, email, password_hash, first_name, last_name
) VALUES (
    ?1, ?2, ?3, ?4, ?5
) RETURNING id, username, email, password_hash, first_name, last_name, is_active, is_admin, created_at, updated_at;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = ?1;

-- name: GetActiveUsers :many
SELECT id, username, email, password_hash, first_name, last_name, is_active, is_admin, created_at, updated_at FROM users
WHERE is_active = TRUE
ORDER BY created_at DESC
LIMIT ?1 OFFSET ?2;

-- name: GetAdminUsers :many
SELECT id, username, email, password_hash, first_name, last_name, is_active, is_admin, created_at, updated_at FROM users
WHERE is_admin = TRUE
ORDER BY created_at DESC;

-- name: GetUser :one
SELECT id, username, email, password_hash, first_name, last_name, is_active, is_admin, created_at, updated_at FROM users
WHERE id = ?1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT id, username, email, password_hash, first_name, last_name, is_active, is_admin, created_at, updated_at FROM users
WHERE email = ?1 LIMIT 1;

-- name: GetUserByUsername :one
SELECT id, username, email, password_hash, first_name, last_name, is_active, is_admin, created_at, updated_at FROM users
WHERE username = ?1 LIMIT 1;

-- name: ListUsers :many
SELECT id, username, email, password_hash, first_name, last_name, is_active, is_admin, created_at, updated_at FROM users
ORDER BY created_at DESC
LIMIT ?1 OFFSET ?2;

-- name: UpdateUser :one
UPDATE users
SET 
    email = COALESCE(?2, email),
    first_name = COALESCE(?3, first_name),
    last_name = COALESCE(?4, last_name),
    is_active = COALESCE(?5, is_active),
    updated_at = CURRENT_TIMESTAMP
WHERE id = ?1
RETURNING id, username, email, password_hash, first_name, last_name, is_active, is_admin, created_at, updated_at;