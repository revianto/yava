-- +goose Up
-- Default recipes (owner_id = NULL, is_default = TRUE)
WITH type_ids AS (
    SELECT id, code FROM yv_cd_recipe_type
),
subtype_ids AS (
    SELECT s.id, s.code AS sub_code, t.code AS type_code
    FROM yv_cd_recipe_subtype s
    JOIN yv_cd_recipe_type t ON t.id = s.type_id
),
r1 AS (
    INSERT INTO yv_recipe (owner_id, type_id, subtype_id, name, description, visibility, is_default, param_dose, param_yield, param_temp, param_grind, param_ratio)
    SELECT NULL, t.id, s.id,
        'V60 Light Roast 15g / 250ml',
        'Pour-over standar untuk light roast. Profil cerah, body ringan, finish clean — disesuaikan untuk bean Kenya AA atau Ethiopia natural.',
        'public', TRUE, '15g', '250ml', '92°C', 'Medium-Fine', '1:16'
    FROM type_ids t, subtype_ids s WHERE t.code = 'v60' AND s.type_code = 'v60' AND s.sub_code = 'regular_drip'
    RETURNING id
),
r2 AS (
    INSERT INTO yv_recipe (owner_id, type_id, subtype_id, name, description, visibility, is_default, param_dose, param_yield, param_temp, param_grind, param_ratio)
    SELECT NULL, t.id, s.id,
        'Aeropress Inverted 17g / 220ml',
        'Metode inverted dengan steep 90 detik. Body lebih kaya dari V60, sweetness terangkat — cocok untuk medium roast.',
        'public', TRUE, '17g', '220ml', '85°C', 'Medium', '1:13'
    FROM type_ids t, subtype_ids s WHERE t.code = 'aeropress' AND s.type_code = 'aeropress' AND s.sub_code = 'inverted'
    RETURNING id
)
-- sessions for V60 default
INSERT INTO yv_recipe_session (recipe_id, sort_order, name, duration_sec, note)
SELECT id, 1, 'Blooming', 45, 'Tuang 45ml air. Gentle swirl. Nikmati aromanya.' FROM r1
UNION ALL
SELECT id, 2, 'First Pour', 100, 'Pour ke 150ml total. Lingkaran konsentris dari tengah.' FROM r1
UNION ALL
SELECT id, 3, 'Second Pour', 100, 'Lanjut ke 250ml. Pour stabil, jangan menyentuh dinding filter.' FROM r1
UNION ALL
SELECT id, 1, 'Steep', 90, 'Pour 220ml, aduk 3x dengan stirrer. Diamkan.' FROM r2
UNION ALL
SELECT id, 2, 'Press', 30, 'Flip dan press perlahan selama 30 detik. Berhenti saat dengar hiss.' FROM r2;

-- +goose Down
DELETE FROM yv_recipe WHERE is_default = TRUE AND owner_id IS NULL;
