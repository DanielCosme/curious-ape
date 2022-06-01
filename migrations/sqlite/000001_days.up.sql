-- schema
CREATE TABLE IF NOT EXISTS days (
    id               INTEGER primary key,
    "date"           DATE not null UNIQUE
);
