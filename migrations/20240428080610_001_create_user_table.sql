-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE IF NOT EXISTS users
(
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    username   TEXT NOT NULL,
    password   TEXT NOT NULL,
    role       TEXT NOT NULL DEFAULT 'user',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
