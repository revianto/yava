-- +goose Up
CREATE TABLE tr_cd_roles (
  id BIGSERIAL,
  code varchar(10) NOT NULL,
  name varchar(100) NOT NULL,
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
  UNIQUE (code)
);

-- +goose Down
DROP TABLE IF EXISTS tr_cd_roles;
