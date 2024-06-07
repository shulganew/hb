-- +goose Up
-- +goose StatementBegin


CREATE TABLE IF NOT EXISTS users (
	user_id UUID NOT NULL UNIQUE DEFAULT gen_random_uuid(), 
	login TEXT NOT NULL UNIQUE, 
	password_hash TEXT NOT NULL,
	email TEXT NOT NULL,
	hb DATE NOT NULL
	);


-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
