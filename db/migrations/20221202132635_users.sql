-- migrate:up
CREATE TABLE IF NOT EXISTS users(id UUID, name varchar);

-- migrate:down
DROP TABLE users;
