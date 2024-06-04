-- schema
CREATE TABLE IF NOT EXISTS days (
        id                  INTEGER PRIMARY KEY,
        "date"              DATE NOT NULL UNIQUE
);

CREATE INDEX IF NOT EXISTS days_date_idx ON days ("date");

CREATE TABLE IF NOT EXISTS habit_categories (
        id                  INTEGER PRIMARY KEY,

        name                TEXT UNIQUE NOT NULL                    CHECK (LENGTH(name) > 0),
        type                TEXT UNIQUE NOT NULL DEFAULT "dynamic"  CHECK (LENGTH(type) > 0),
        description         TEXT NOT NULL DEFAULT ""

        CHECK(LENGTH(id) > 0)
);

CREATE TABLE IF NOT EXISTS habits (
        id                  INTEGER PRIMARY KEY,
        day_id              INTEGER NOT NULL,
        habit_category_id   INTEGER NOT NULL,

        state              TEXT NOT NULL CHECK (state IN ('no_info', 'not_done', 'done')),

        FOREIGN KEY (habit_category_id) REFERENCES habit_categories (id) ON DELETE CASCADE,
        FOREIGN KEY (day_id) REFERENCES "days" (id) ON DELETE CASCADE,
        UNIQUE (day_id, habit_category_id),
        CHECK(length(id) > 0 AND length(day_id) > 0 AND length(habit_category_id) > 0)
);

CREATE TABLE IF NOT EXISTS habit_logs (
        id                  INTEGER PRIMARY KEY,
        habit_id            INTEGER NOT NULL,

        origin              TEXT NOT NULL DEFAULT "unknown" CHECK (length(origin) > 0),
        success             BOOLEAN NOT NULL DEFAULT false,
        is_automated        BOOLEAN NOT NULL DEFAULT false,

        FOREIGN KEY (habit_id) REFERENCES habits (id) ON DELETE CASCADE,
        UNIQUE (habit_id, origin),
        CHECK(length(id) > 0 AND length(habit_id) > 0)
);

CREATE TABLE IF NOT EXISTS auths (
        id                  INTEGER primary key,

        provider            TEXT NOT NULL UNIQUE,
        access_token        TEXT NOT NULL UNIQUE,
        refresh_token       TEXT,
        token_type          TEXT,
        expiration          DATE
);

CREATE TABLE IF NOT EXISTS users (
        id              INTEGER PRIMARY KEY,

        username        TEXT NOT NULL UNIQUE,
        password        TEXT NOT NULL UNIQUE,
        role            TEXT CHECK (role IN ('admin', 'user', 'guest')) NOT NULL,
        email           TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS sleep_logs (
        id                  INTEGER PRIMARY KEY,
        day_id              INTEGER NOT NULL,

        "date"              DATE NOT NULL,
        start_time          DATE NOT NULL,
        end_time            DATE NOT NULL,
        is_main_sleep       BOOLEAN DEFAULT TRUE,
        is_automated        BOOLEAN DEFAULT FALSE,
        origin              TEXT NOT NULL CHECK (length(origin) > 1),
        total_time_in_bed   INTEGER DEFAULT 0,
        minutes_asleep      INTEGER DEFAULT 0,
        minutes_rem         INTEGER DEFAULT 0,
        minutes_deep        INTEGER DEFAULT 0,
        minutes_light       INTEGER DEFAULT 0,
        minutes_awake       INTEGER DEFAULT 0,
        raw                 TEXT,

        FOREIGN KEY (day_id) REFERENCES "days" (id) ON DELETE CASCADE,
        UNIQUE (day_id, is_main_sleep)
);

CREATE TABLE IF NOT EXISTS fitness_logs (
        id                  INTEGER PRIMARY KEY,
        day_id              INTEGER NOT NULL,

        "date"              DATE NOT NULL,
        start_time          DATE NOT NULL,
        end_time            DATE NOT NULL,
        "type"              TEXT NOT NULL DEFAULT '',
        title               TEXT NOT NULL DEFAULT '',
        origin              TEXT NOT NULL CHECK (length(origin) > 1),
        note                TEXT DEFAULT '',
        raw                 TEXT,

        FOREIGN KEY (day_id) REFERENCES "days" (id) ON DELETE CASCADE,
        UNIQUE (day_id, start_time)
);

CREATE TABLE IF NOT EXISTS sessions (
    token CHAR(43) PRIMARY KEY,
    data BLOB NOT NULL,
    expiry TIMESTAMP(6) NOT NULL
);

CREATE INDEX IF NOT EXISTS sessions_expiry_idx ON sessions (expiry);
