-- schema
CREATE TABLE IF NOT EXISTS day (
        id                  INTEGER PRIMARY KEY,
        "date"              DATE NOT NULL UNIQUE
);

CREATE INDEX IF NOT EXISTS day_date_idx ON day ("date");

CREATE TABLE IF NOT EXISTS habit_category (
        id                  INTEGER PRIMARY KEY,

        name                TEXT UNIQUE NOT NULL CHECK (LENGTH(name) > 0),
        kind                TEXT UNIQUE NOT NULL CHECK (LENGTH(kind) > 0),
        description         TEXT NOT NULL DEFAULT "",

        CHECK(LENGTH(id) > 0)
);

CREATE TABLE IF NOT EXISTS habit (
        id                  INTEGER PRIMARY KEY,
        day_id              INTEGER NOT NULL,
        habit_category_id   INTEGER NOT NULL,

        state              TEXT NOT NULL CHECK (state IN ('done', 'not_done','no_info')),
        automated          BOOLEAN NOT NULL DEFAULT false,

        FOREIGN KEY (habit_category_id) REFERENCES habit_category (id) ON DELETE CASCADE,
        FOREIGN KEY (day_id) REFERENCES "day" (id) ON DELETE CASCADE,
        UNIQUE (day_id, habit_category_id),
        CHECK(length(id) > 0 AND length(day_id) > 0 AND length(habit_category_id) > 0)
);

CREATE TABLE IF NOT EXISTS oauth_token (
        id                  INTEGER primary key,

        provider            TEXT NOT NULL UNIQUE,
        access_token        TEXT NOT NULL UNIQUE,
        refresh_token       TEXT NOT NULL DEFAULT "",
        token_type          TEXT NOT NULL DEFAULT "",
        expiration          DATE NOT NULL DEFAULT '1992-01-21'
);

CREATE TABLE IF NOT EXISTS user (
        id              INTEGER PRIMARY KEY,

        username        TEXT NOT NULL UNIQUE,
        password        TEXT NOT NULL UNIQUE,
        role            TEXT NOT NULL CHECK (role IN ('admin', 'user', 'guest')),
        email           TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS sessions (
    token TEXT PRIMARY KEY,
    data BLOB NOT NULL,
    expiry REAL NOT NULL
);

CREATE INDEX IF NOT EXISTS sessions_expiry_idx ON sessions (expiry);
