-- name: CreateURL :one
INSERT INTO urls (url, title, description, is_pinned, category_id, user_id)
VALUES (?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: GetURLByID :one
SELECT * FROM urls
WHERE id = ?;

-- name: ListURLsByUser :many
SELECT * FROM urls
WHERE user_id = ?
ORDER BY is_pinned DESC, created_at DESC;

-- name: ListURLsByCategory :many
SELECT * FROM urls
WHERE category_id = ? AND user_id = ?
ORDER BY is_pinned DESC, created_at DESC;

-- name: UpdateURL :one
UPDATE urls
SET url = ?, title = ?, description = ?, is_pinned = ?, category_id = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ? AND user_id = ?
RETURNING *;

-- name: DeleteURL :exec
DELETE FROM urls
WHERE id = ? AND user_id = ?;

-- name: SearchURLs :many
SELECT * FROM urls
WHERE user_id = ? AND (
    url LIKE ? OR
    title LIKE ? OR
    description LIKE ?
)
ORDER BY is_pinned DESC, created_at DESC;
