-- schema
CREATE TABLE IF NOT EXISTS days (
    id               INTEGER primary key,
    "date"           TEXT not null UNIQUE
);
