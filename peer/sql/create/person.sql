CREATE TABLE person(
	id VARCHAR(32) PRIMARY KEY,
	name VARCHAR(128),
	pass VARCHAR(64),
	profile VARCHAR(1024),
	key VARCHAR(256)
);
CREATE INDEX person_key1 ON person(key);
