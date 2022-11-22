CREATE DATABASE gobb_dev;
USE gobb_dev;

CREATE TABLE replays (
	id uuid NOT NULL,
	home_team jsonb NOT NULL,
	away_team jsonb NOT NULL
);
