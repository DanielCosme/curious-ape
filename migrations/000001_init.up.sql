CREATE EXTENSION IF NOT EXISTS citext;
CREATE TYPE habit_state AS ENUM ('yes', 'no', 'no_info');
CREATE TYPE habit_type AS ENUM ('sleep', 'fitness', 'work', 'food');

CREATE TABLE IF NOT EXISTS habits (
    id          SERIAL PRIMARY KEY,
    state       habit_state NOT NULL DEFAULT 'no_info',
    "date"      DATE NOT NULL,
    "type"      habit_type NOT NULL,
    origin      TEXT NOT NULL DEFAULT 'unknown',
    UNIQUE("date", "type")
);

CREATE TABLE IF NOT EXISTS users (
    id             bigserial PRIMARY KEY,
    created_at     timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    name           text NOT NULL,
    email          citext UNIQUE NOT NULL,
    password_hash  bytea NOT NULL,
    activated      bool NOT NULL
);

create TABLE sleep_records (
    id               SERIAL PRIMARY KEY,
    "date"           DATE UNIQUE NOT NULL,
    duration         INT NOT NULL,
    start_time       TIMESTAMP NOT NULL,
    end_time         TIMESTAMP,
    minutes_asleep   INT NOT NULL,
    minutes_awake    INT NOT NULL,
    minutes_in_bed   INT NOT NULL,
    provider         TEXT NOT NULL,
    raw_json         JSONB NOT NULL

);
CREATE TABLE auth_tokens (
    service         VARCHAR(50) PRIMARY KEY NOT NULL,
    access_token    TEXT UNIQUE NOT NULL,
    refresh_token   TEXT UNIQUE NOT NULL
);

INSERT INTO auth_tokens (service, access_token, refresh_token)
VALUES
    ('fitbit', 'f', 'f'),
    ('google', 'g', 'g');

CREATE TABLE work_records (
    id              SERIAL PRIMARY KEY,
    "date"          DATE NOT NULL UNIQUE,
    grand_total     INT NOT NULL CHECK(grand_total >= 0),
    raw_json        JSONB NOT NULL,
    provider        TEXT NOT NULL
);

CREATE TABLE fitness_records (
    id                      SERIAL PRIMARY KEY,
    "date"                  DATE UNIQUE NOT NULL,
    start_in_miliseconds    BIGINT NOT NULL,
    end_in_miliseconds      BIGINT NOT NULL,
    provider                TEXT NOT NULL
);
