-- +goose Up
CREATE TABLE IF NOT EXISTS passwords
(
    user_id VARCHAR(36) NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT (NOW() AT TIME ZONE 'Europe/Moscow'),
    expires_at TIMESTAMP NOT NULL,
    PRIMARY KEY (user_id)
);

-- +goose Down
DROP TABLE passwords;