CREATE TABLE person(
	id CHAR(32) PRIMARY KEY,
	name CHAR(128),
	pass CHAR(64),
	profile VARCHAR(1024),
	key CHAR(256)
);
CREATE INDEX person_key1 ON person(key);
