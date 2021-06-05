CREATE TABLE fitness_records (
    id SERIAL PRIMARY KEY,
    "date" DATE UNIQUE NOT NULL,
    start_in_miliseconds BIGINT NOT NULL,
    end_in_miliseconds BIGINT NOT NULL,
    provider TEXT NOT NULL
);

