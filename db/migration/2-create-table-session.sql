-- +migrate Up notransaction

CREATE TABLE IF NOT EXISTS "sessions" (
    id BIGINT PRIMARY KEY,
    username TEXT NOT NULL,
    access_token TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    expired_at TIMESTAMP NOT NULL
);

ALTER TABLE "sessions" ADD FOREIGN KEY (username) REFERENCES "users" ("username");

-- +migrate Down

DROP TABLE IF EXISTS "sessions";