CREATE TABLE IF NOT EXISTS sleep_log (
        id                  INTEGER PRIMARY KEY,
        day_id              INTEGER NOT NULL,
        title               TEXT NOT NULL CHECK (length(origin) > 2),
        "date"              DATE NOT NULL,
        start_time          DATE NOT NULL,
        end_time            DATE NOT NULL,
        note                TEXT NOT NULL,

        is_main_sleep       BOOLEAN NOT NULL DEFAULT TRUE,
        minutes_in_bed      INTEGER NOT NULL DEFAULT 0,
        minutes_asleep      INTEGER NOT NULL DEFAULT 0,

        origin              TEXT NOT NULL CHECK (length(origin) > 3),
        raw                 TEXT NOT NULL,

        FOREIGN KEY (day_id) REFERENCES "day" (id) ON DELETE CASCADE,
        CONSTRAINT unique_day_main_sleep UNIQUE (day_id, is_main_sleep)
);

CREATE TABLE IF NOT EXISTS fitness_log (
        id                  INTEGER PRIMARY KEY,
        day_id              INTEGER NOT NULL,
        title               TEXT NOT NULL CHECK (length(origin) > 2),
        "date"              DATE NOT NULL,
        start_time          DATE NOT NULL,
        end_time            DATE NOT NULL,
        note                TEXT NOT NULL,

        "type"              TEXT NOT NULL DEFAULT '',
  
        origin              TEXT NOT NULL CHECK (length(origin) > 3),
        raw                 TEXT NOT NULL,

        FOREIGN KEY (day_id) REFERENCES "day" (id) ON DELETE CASCADE,
        CONSTRAINT unique_day_start_time UNIQUE (day_id, start_time)
);

CREATE TABLE IF NOT EXISTS deep_work_log (
        id                  INTEGER PRIMARY KEY,
        day_id              INTEGER NOT NULL,
        title               TEXT NOT NULL CHECK (length(origin) > 2),
        "date"              DATE NOT NULL,
        start_time          DATE NOT NULL,
        -- end_time            DATE NOT NULL,
        -- note                TEXT NOT NULL,

        seconds             INTEGER NOT NULL CHECK (length(seconds) < 60),

        origin              TEXT NOT NULL CHECK (length(origin) > 3),
        raw                 TEXT NOT NULL,

        UNIQUE (day_id, origin),
        FOREIGN KEY (day_id) REFERENCES "day" (id) ON DELETE CASCADE
        CONSTRAINT unique_day_start_time UNIQUE (day_id, start_time)
);

