-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
    id         INTEGER  PRIMARY KEY,
    username   TEXT     NOT NULL UNIQUE,
    password   TEXT     NOT NULL,
    name       TEXT     NOT NULL,
    is_admin   BOOLEAN  NOT NULL DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
