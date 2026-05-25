-- +goose Up
CREATE TABLE yv_cd_recipe_subtype (
    id          BIGSERIAL PRIMARY KEY,
    type_id     BIGINT NOT NULL REFERENCES yv_cd_recipe_type(id),
    code        VARCHAR(50) NOT NULL,
    name        VARCHAR(100) NOT NULL,
    sort_order  INT DEFAULT 0 NOT NULL,
    created_at  TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    updated_at  TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    UNIQUE (type_id, code)
);

CREATE INDEX idx_yv_cd_recipe_subtype_type_id ON yv_cd_recipe_subtype(type_id);

-- +goose Down
DROP TABLE IF EXISTS yv_cd_recipe_subtype;
