-- schema
CREATE TABLE IF NOT EXISTS days (
    id                  INTEGER PRIMARY KEY,
    "date"              DATE NOT NULL UNIQUE
);