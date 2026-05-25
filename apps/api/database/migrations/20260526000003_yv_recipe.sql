-- +goose Up
CREATE TABLE yv_recipe (
    id            BIGSERIAL PRIMARY KEY,
    owner_id      BIGINT REFERENCES yv_user(id) ON DELETE SET NULL,
    type_id       BIGINT NOT NULL REFERENCES yv_cd_recipe_type(id),
    subtype_id    BIGINT REFERENCES yv_cd_recipe_subtype(id),
    name          VARCHAR(255) NOT NULL,
    description   TEXT,
    visibility    VARCHAR(20) DEFAULT 'private' NOT NULL,
    is_default    BOOLEAN DEFAULT FALSE NOT NULL,
    is_archived   BOOLEAN DEFAULT FALSE NOT NULL,
    saves_count   INT DEFAULT 0 NOT NULL,
    param_dose    VARCHAR(50),
    param_yield   VARCHAR(50),
    param_temp    VARCHAR(50),
    param_grind   VARCHAR(50),
    param_ratio   VARCHAR(50),
    last_brewed_at TIMESTAMPTZ,
    created_at    TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    updated_at    TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    deleted_at    TIMESTAMPTZ
);

-- §12.3 indexes
CREATE INDEX idx_yv_recipe_owner_id   ON yv_recipe(owner_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_yv_recipe_visibility ON yv_recipe(visibility) WHERE deleted_at IS NULL;
CREATE INDEX idx_yv_recipe_type_id    ON yv_recipe(type_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_yv_recipe_is_default ON yv_recipe(is_default) WHERE deleted_at IS NULL;

-- +goose Down
DROP TABLE IF EXISTS yv_recipe;
