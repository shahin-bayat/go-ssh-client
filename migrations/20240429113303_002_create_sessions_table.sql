-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE sessions(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL REFERENCES users(id),
    token VARCHAR(255) NOT NULL,
    expires_at TIMESTAMP NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE sessions;
-- +goose StatementEnd
