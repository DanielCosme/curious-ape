-- schema
CREATE TABLE IF NOT EXISTS days (
    id                  INTEGER primary key,
    "date"              DATE not null UNIQUE,
    deep_work_minutes   INTEGER NOT NULL DEFAULT 0 check ( length(deep_work_minutes) >= 0)
);
