-- +migrate Up notransaction

CREATE TABLE IF NOT EXISTS "users" (
    username TEXT PRIMARY KEY,
    password TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP DEFAULT NULL
);

-- +migrate Down

DROP TABLE IF EXISTS "users";