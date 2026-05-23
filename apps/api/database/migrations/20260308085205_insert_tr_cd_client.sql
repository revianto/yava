-- +goose Up
INSERT INTO tr_cd_client
(id, name, secret, created_time, created_by, created_from, modified_time, modified_by, modified_from, deleted_time, deleted_by, deleted_from)
VALUES(1, 'Web', '$2a$14$WGL7YvLyLf1oZgxiOiWaMO0Vsl4RauJOMLf6IB5SxfIHF9e73FNri', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL);
SELECT setval('tr_cd_client_id_seq', (SELECT COALESCE(MAX(id), 0) FROM tr_cd_client));

-- +goose Down
DELETE FROM tr_cd_client WHERE id = 1;
