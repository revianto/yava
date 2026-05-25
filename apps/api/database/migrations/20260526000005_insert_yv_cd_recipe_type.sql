-- +goose Up
INSERT INTO yv_cd_recipe_type (code, name, sort_order) VALUES
    ('espresso',   'Espresso',   1),
    ('v60',        'V60',        2),
    ('aeropress',  'Aeropress',  3),
    ('cold_brew',  'Cold Brew',  4),
    ('moka_pot',   'Moka Pot',   5),
    ('chemex',     'Chemex',     6),
    ('french_press','French Press',7);

-- +goose Down
DELETE FROM yv_cd_recipe_type WHERE code IN ('espresso','v60','aeropress','cold_brew','moka_pot','chemex','french_press');
