-- +goose Up
ALTER TABLE tr_users
  ADD COLUMN is_active BOOLEAN NOT NULL DEFAULT TRUE;

-- +goose Down
ALTER TABLE tr_users DROP COLUMN is_active;
