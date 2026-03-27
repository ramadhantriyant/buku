-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS urls (
    id          INTEGER  PRIMARY KEY,
    url         TEXT     NOT NULL,
    title       TEXT,
    description TEXT,
    is_pinned   BOOLEAN  NOT NULL DEFAULT 0,
    category_id INTEGER  REFERENCES categories(id) ON DELETE SET NULL,
    user_id     INTEGER  NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at  DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at  DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(url, user_id)
);

CREATE INDEX idx_urls_user_id     ON urls(user_id);
CREATE INDEX idx_urls_category_id ON urls(category_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_urls_category_id;
DROP INDEX IF EXISTS idx_urls_user_id;
DROP TABLE IF EXISTS urls;
-- +goose StatementEnd
