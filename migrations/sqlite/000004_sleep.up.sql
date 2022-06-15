CREATE TABLE IF NOT EXISTS sleep_logs (
    id                  INTEGER primary key,
    day_id              INTEGER not null,

    start_time          DATE not null,
    end_time            DATE not null,
    is_main_sleep       BOOLEAN default true,
    is_automated        BOOLEAN default false,
    origin              TEXT not null CHECK (length(origin) > 1),
    total_time_in_bed   INTEGER default 0,
    minutes_asleep      INTEGER default 0,
    minutes_rem         INTEGER default 0,
    minutes_deep        INTEGER default 0,
    minutes_light       INTEGER default 0,
    minutes_awake       INTEGER default 0,
    raw                 TEXT,

    FOREIGN KEY (day_id) REFERENCES "days" (id) ON DELETE CASCADE,
    UNIQUE (day_id, is_main_sleep)
);
