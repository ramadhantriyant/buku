-- name: CreateCategory :one
INSERT INTO categories (name, description, user_id)
VALUES (?, ?, ?)
RETURNING *;

-- name: GetCategoryByID :one
SELECT * FROM categories
WHERE id = ?;

-- name: ListCategoriesByUser :many
SELECT * FROM categories
WHERE user_id = ?
ORDER BY name ASC;

-- name: UpdateCategory :one
UPDATE categories
SET name = ?, description = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ? AND user_id = ?
RETURNING *;

-- name: DeleteCategory :exec
DELETE FROM categories
WHERE id = ? AND user_id = ?;
