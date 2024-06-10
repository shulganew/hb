-- +goose Up
-- +goose StatementBegin


CREATE TABLE IF NOT EXISTS users (
	user_id UUID NOT NULL UNIQUE DEFAULT gen_random_uuid(), 
	tg_user TEXT NOT NULL UNIQUE, 
	name TEXT NOT NULL UNIQUE, 
	password_hash TEXT NOT NULL,
	hb DATE NOT NULL
	);

CREATE TABLE IF NOT EXISTS subscription (
	tg_user TEXT NOT NULL REFERENCES users(tg_user),
	subscribed TEXT NOT NULL REFERENCES users(tg_user),
	chat_id BIGINT NOT NULL,
	UNIQUE (tg_user, subscribed)
	);

-- insert test users data
INSERT INTO users (tg_user, name, password_hash, hb) VALUES ('shulgaigor', 'Igor', '$2a$10$VQdSYvDhYtzjTE72FwtxU.6Pm5aT50baHWKVI3mtb8arQG18O5X9i', '1983-03-11');
INSERT INTO users (tg_user, name, password_hash, hb) VALUES ('Oleg', 'Oleg', '$2a$10$VQdSYvDhYtzjTE72FwtxU.6Pm5aT50baHWKVI3mtb8arQG18O5X9i', '1975-10-20');
INSERT INTO users (tg_user, name, password_hash, hb) VALUES ('Anya', 'Anya', '$2a$10$VQdSYvDhYtzjTE72FwtxU.6Pm5aT50baHWKVI3mtb8arQG18O5X9i', '1999-04-14');
INSERT INTO users (tg_user, name, password_hash, hb) VALUES ('Marina', 'Marina', '$2a$10$VQdSYvDhYtzjTE72FwtxU.6Pm5aT50baHWKVI3mtb8arQG18O5X9i', '2000-01-11');
INSERT INTO users (tg_user, name, password_hash, hb) VALUES ('MrToday', 'MrToday', '$2a$10$VQdSYvDhYtzjTE72FwtxU.6Pm5aT50baHWKVI3mtb8arQG18O5X9i', '2000-06-10');
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd




-- INSERT INTO subscription VALUES ('shulgaigor', 'Oleg');