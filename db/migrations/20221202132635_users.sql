-- migrate:up
CREATE TABLE IF NOT EXISTS users(id UUID, name varchar);

-- migrate:down
DROP TABLE users;

postgres://username:password@127.0.0.1:5432/database_name?sslmode=disable