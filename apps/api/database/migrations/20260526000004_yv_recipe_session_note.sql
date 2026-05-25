-- +goose Up
CREATE TABLE yv_recipe_session (
    id          BIGSERIAL PRIMARY KEY,
    recipe_id   BIGINT NOT NULL REFERENCES yv_recipe(id) ON DELETE CASCADE,
    sort_order  INT NOT NULL,
    name        VARCHAR(255) NOT NULL,
    duration_sec INT NOT NULL,
    note        TEXT,
    created_at  TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    updated_at  TIMESTAMPTZ DEFAULT NOW() NOT NULL
);

CREATE INDEX idx_yv_recipe_session_recipe_id ON yv_recipe_session(recipe_id);

CREATE TABLE yv_recipe_note (
    id          BIGSERIAL PRIMARY KEY,
    recipe_id   BIGINT NOT NULL REFERENCES yv_recipe(id) ON DELETE CASCADE,
    sort_order  INT NOT NULL,
    content     TEXT NOT NULL,
    created_at  TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    updated_at  TIMESTAMPTZ DEFAULT NOW() NOT NULL
);

CREATE INDEX idx_yv_recipe_note_recipe_id ON yv_recipe_note(recipe_id);

-- +goose Down
DROP TABLE IF EXISTS yv_recipe_note;
DROP TABLE IF EXISTS yv_recipe_session;
