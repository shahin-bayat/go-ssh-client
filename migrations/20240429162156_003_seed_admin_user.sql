-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
INSERT INTO users (username, password, role)
VALUES ('admin', '$2a$10$3XIDakUeV9mvfEIcMmDpAO4VeZOHwnZFOoa1uza2bW/MbnXU..e32', 'admin');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DELETE FROM users WHERE username = 'admin';
-- +goose StatementEnd
