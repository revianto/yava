-- +goose Up
CREATE TABLE tr_oauth_tokens (
	id varchar(255) DEFAULT NULL,
	refresh_id varchar(255) DEFAULT NULL,
	client_id INTEGER DEFAULT NULL,
	user_id BIGINT DEFAULT NULL,
	revoked INTEGER DEFAULT NULL,
	expiary_time timestamp DEFAULT NULL,
	refresh_expiary_time timestamp DEFAULT NULL,
	created_time timestamp DEFAULT NULL,
	created_from varchar(20) DEFAULT NULL,
	modified_time timestamp DEFAULT NULL,
	modified_from varchar(20) DEFAULT NULL
);

CREATE INDEX idx_ot_cd_oauth_token_client_id ON tr_oauth_tokens (client_id);
CREATE INDEX idx_ot_cd_oauth_token_refresh_id ON tr_oauth_tokens (refresh_id);

-- +goose Down
DROP TABLE IF EXISTS tr_oauth_tokens;
