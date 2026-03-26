-- name: CreateUser :one
INSERT INTO users (username, password, name, is_admin)
VALUES (?, ?, ?, ?)
RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = ?;

-- name: GetUserByUsername :one
SELECT * FROM users
WHERE username = ?;

-- name: UpdateUser :one
UPDATE users
SET name = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ?
RETURNING *;

-- name: UpdateUserPassword :exec
UPDATE users
SET password = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ?;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = ?;

-- name: CountUsers :one
SELECT COUNT(*) FROM users;
