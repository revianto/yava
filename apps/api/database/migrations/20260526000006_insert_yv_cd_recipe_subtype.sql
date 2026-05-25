-- +goose Up
INSERT INTO yv_cd_recipe_subtype (type_id, code, name, sort_order)
SELECT t.id, s.code, s.name, s.sort_order FROM yv_cd_recipe_type t
JOIN (VALUES
    ('espresso',  'manual',      'Manual (Flair/Cafelat)', 1),
    ('espresso',  'machine',     'Machine',                2),
    ('espresso',  'milk_based',  'Milk-based',             3),
    ('v60',       'regular_drip','Regular Drip',           1),
    ('v60',       'chemex',      'Chemex',                 2),
    ('v60',       'kalita',      'Kalita Wave',            3),
    ('aeropress', 'regular',     'Regular',                1),
    ('aeropress', 'inverted',    'Inverted',               2),
    ('cold_brew', 'slow_drip',   'Slow Drip',              1),
    ('cold_brew', 'immersion',   'Immersion',              2),
    ('moka_pot',  'stovetop',    'Stovetop',               1),
    ('chemex',    'classic',     'Classic',                1),
    ('french_press','classic',   'Classic',                1)
) AS s(type_code, code, name, sort_order) ON t.code = s.type_code;

-- +goose Down
DELETE FROM yv_cd_recipe_subtype;
