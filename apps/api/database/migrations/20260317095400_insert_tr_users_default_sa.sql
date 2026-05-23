-- +goose Up
-- Default Super Admin user
-- PENTING: Ganti password segera setelah pertama login melalui API reset-password
INSERT INTO tr_users (id, name, phone, password, role_id, created_time, created_by, created_from, modified_time, modified_by, modified_from, deleted_time, deleted_by, deleted_from) VALUES
(1, 'Super Admin', '+6200000000000', '$2a$14$EfziWG3mP9K874kJJGon7uA/BgS6UbQnflLr5wKlndJvryKuAcRJm', 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL);
SELECT setval('tr_users_id_seq', (SELECT COALESCE(MAX(id), 0) FROM tr_users));

-- +goose Down
DELETE FROM tr_users WHERE id = 1;
