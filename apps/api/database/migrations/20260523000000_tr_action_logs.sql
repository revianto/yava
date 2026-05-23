-- +goose Up
CREATE TABLE IF NOT EXISTS tr_action_logs (
    id          BIGSERIAL PRIMARY KEY,
    user_id     BIGINT       NOT NULL,
    action_text TEXT         NOT NULL,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_action_logs_user_id    ON tr_action_logs (user_id);
CREATE INDEX idx_action_logs_created_at ON tr_action_logs (created_at);

-- +goose Down
DROP TABLE IF EXISTS tr_action_logs;
