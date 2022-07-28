CREATE TABLE IF NOT EXISTS fitness_logs (
    id                  INTEGER primary key,
    day_id              INTEGER not null,

    "date"              DATE not null,
    start_time          DATE not null,
    end_time            DATE not null,
    "type"              TEXT not null default '',
    title               TEXT not null default '',
    origin              TEXT not null CHECK (length(origin) > 1),
    note                TEXT default '',
    raw                 TEXT,

    FOREIGN KEY (day_id) REFERENCES "days" (id) ON DELETE CASCADE,
    UNIQUE (day_id, start_time)
);
