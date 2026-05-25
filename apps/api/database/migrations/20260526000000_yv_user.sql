-- +goose Up
CREATE TABLE yv_user (
    id          BIGSERIAL PRIMARY KEY,
    google_id   VARCHAR(255) UNIQUE NOT NULL,
    email       VARCHAR(255) UNIQUE NOT NULL,
    name        VARCHAR(255) NOT NULL,
    avatar_url  TEXT,
    created_at  TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    updated_at  TIMESTAMPTZ DEFAULT NOW() NOT NULL
);

CREATE INDEX idx_yv_user_google_id ON yv_user(google_id);
CREATE INDEX idx_yv_user_email ON yv_user(email);

-- +goose Down
DROP TABLE IF EXISTS yv_user;
