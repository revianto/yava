-- +goose Up
CREATE TABLE yv_cd_recipe_type (
    id          BIGSERIAL PRIMARY KEY,
    code        VARCHAR(50) UNIQUE NOT NULL,
    name        VARCHAR(100) NOT NULL,
    sort_order  INT DEFAULT 0 NOT NULL,
    created_at  TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    updated_at  TIMESTAMPTZ DEFAULT NOW() NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS yv_cd_recipe_type;
