-- +goose Up
ALTER TABLE tr_users
  ADD COLUMN reset_code varchar(6) DEFAULT NULL,
  ADD COLUMN reset_code_expired_at timestamp NULL DEFAULT NULL;

-- +goose Down
ALTER TABLE tr_users
  DROP COLUMN reset_code,
  DROP COLUMN reset_code_expired_at;
