-- Add migration script here
CREATE TABLE habits(
    id serial NOT NULL,
    PRIMARY KEY (id),
    name TEXT NOT NULL UNIQUE,
    description TEXT DEFAULT '',
    created_at timestamptz NOT NULL
);
