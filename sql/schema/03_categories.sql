-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS categories (
    id          INTEGER  PRIMARY KEY,
    name        TEXT     NOT NULL,
    description TEXT,
    color       TEXT,
    user_id     INTEGER  NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at  DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at  DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(name, user_id)
);

CREATE INDEX idx_categories_name ON categories(name);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_categories_name;
DROP TABLE IF EXISTS categories;
-- +goose StatementEnd
