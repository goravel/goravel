CREATE TABLE users (
  id SERIAL PRIMARY KEY NOT NULL,
  name varchar,
  avatar varchar,
  created_at timestamp NOT NULL,
  updated_at timestamp NOT NULL,
  deleted_at timestamp
);
