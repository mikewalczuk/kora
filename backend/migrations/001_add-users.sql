-- +goose Up
CREATE TABLE users (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username      TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at    TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Default user: admin / admin
INSERT INTO users (username, password_hash)
VALUES ('admin', '$2a$10$Ybv.w9Gk3f2UsHDSVqkKROHvSdCFTCtwT0L/B13zBr3YoGJ2foWu6');

-- +goose Down
DROP TABLE users;