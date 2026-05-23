-- +goose Up
CREATE TABLE tr_users (
  id BIGSERIAL,
  name varchar(255) NOT NULL,
  phone varchar(20) NOT NULL,
  password varchar(255) NOT NULL,
  role_id BIGINT NOT NULL,
  created_time timestamp NULL DEFAULT NULL,
  created_by INTEGER DEFAULT NULL,
  created_from varchar(255) DEFAULT NULL,
  modified_time timestamp NULL DEFAULT NULL,
  modified_by INTEGER DEFAULT NULL,
  modified_from varchar(255) DEFAULT NULL,
  deleted_time timestamp NULL DEFAULT NULL,
  deleted_by INTEGER DEFAULT NULL,
  deleted_from varchar(255) DEFAULT NULL,
  PRIMARY KEY (id),
  UNIQUE (phone)
);

CREATE INDEX idx_users_role_id ON tr_users (role_id);

-- +goose Down
DROP TABLE IF EXISTS tr_users;
