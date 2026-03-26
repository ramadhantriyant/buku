-- name: CreatePasswordResetToken :one
INSERT INTO password_reset_tokens (user_id, token_hash, expires_at)
VALUES (?, ?, ?)
RETURNING *;

-- name: GetPasswordResetTokenByHash :one
SELECT * FROM password_reset_tokens
WHERE token_hash = ? AND used = 0 AND expires_at > CURRENT_TIMESTAMP;

-- name: MarkPasswordResetTokenUsed :exec
UPDATE password_reset_tokens
SET used = 1, updated_at = CURRENT_TIMESTAMP
WHERE id = ?;

-- name: DeleteExpiredPasswordResetTokens :exec
DELETE FROM password_reset_tokens
WHERE expires_at < CURRENT_TIMESTAMP OR used = 1;

-- name: DeleteUserPasswordResetTokens :exec
DELETE FROM password_reset_tokens
WHERE user_id = ?;
