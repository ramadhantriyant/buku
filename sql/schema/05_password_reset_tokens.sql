-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS password_reset_tokens (
    id         INTEGER  PRIMARY KEY,
    user_id    INTEGER  NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash TEXT     NOT NULL UNIQUE,
    expires_at DATETIME NOT NULL,
    used       BOOLEAN  NOT NULL DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_prt_user_id ON password_reset_tokens(user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_prt_user_id;
DROP TABLE IF EXISTS password_reset_tokens;
-- +goose StatementEnd
