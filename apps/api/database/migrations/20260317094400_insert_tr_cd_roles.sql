-- +goose Up
INSERT INTO tr_cd_roles (id, name, code, created_time, created_by, created_from, modified_time, modified_by, modified_from, deleted_time, deleted_by, deleted_from) VALUES
(1, 'Super Admin', 'SA', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(2, 'Admin', 'AD', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(3, 'Bagian Gudang', 'BG', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(4, 'Kepala Tim Produksi', 'KTP', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(5, 'Tim Produksi', 'TP', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL),
(6, 'Tim Packing', 'TPC', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL);
SELECT setval('tr_cd_roles_id_seq', (SELECT COALESCE(MAX(id), 0) FROM tr_cd_roles));

-- +goose Down
DELETE FROM tr_cd_roles WHERE id IN (1, 2, 3, 4, 5, 6);
