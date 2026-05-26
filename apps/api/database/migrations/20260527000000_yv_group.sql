-- +goose Up

CREATE TABLE yv_group (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    avatar_url TEXT,
    invite_code VARCHAR(32) NOT NULL UNIQUE,
    created_by BIGINT NOT NULL REFERENCES yv_user(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE yv_group_member (
    id BIGSERIAL PRIMARY KEY,
    group_id BIGINT NOT NULL REFERENCES yv_group(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL REFERENCES yv_user(id) ON DELETE CASCADE,
    role VARCHAR(20) NOT NULL DEFAULT 'member' CHECK (role IN ('founder', 'admin', 'member')),
    joined_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (group_id, user_id)
);

CREATE TABLE yv_group_recipe (
    id BIGSERIAL PRIMARY KEY,
    group_id BIGINT NOT NULL REFERENCES yv_group(id) ON DELETE CASCADE,
    recipe_id BIGINT NOT NULL REFERENCES yv_recipe(id) ON DELETE CASCADE,
    submitted_by BIGINT NOT NULL REFERENCES yv_user(id),
    status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'approved', 'rejected')),
    reviewed_by BIGINT REFERENCES yv_user(id),
    reviewed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (group_id, recipe_id)
);

CREATE INDEX idx_yv_group_invite_code ON yv_group(invite_code);
CREATE INDEX idx_yv_group_member_group ON yv_group_member(group_id);
CREATE INDEX idx_yv_group_member_user ON yv_group_member(user_id);
CREATE INDEX idx_yv_group_recipe_group ON yv_group_recipe(group_id);
CREATE INDEX idx_yv_group_recipe_status ON yv_group_recipe(group_id, status);

-- +goose Down
DROP TABLE IF EXISTS yv_group_recipe;
DROP TABLE IF EXISTS yv_group_member;
DROP TABLE IF EXISTS yv_group;
