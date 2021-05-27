create TABLE sleep_records (
    id SERIAL PRIMARY KEY,
    "date" DATE UNIQUE NOT NULL,
    duration INT NOT NULL,
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP,
    minutes_asleep INT NOT NULL,
    minutes_awake INT NOT NULL,
    minutes_in_bed INT NOT NULL,
    provider TEXT NOT NULL,
    raw_json JSONB NOT NULL
)
