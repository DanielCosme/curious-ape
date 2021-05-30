CREATE TYPE habit_state AS ENUM ('yes', 'no', 'no_info');
CREATE TYPE habit_type AS ENUM ('sleep', 'fitness', 'work', 'food');

CREATE TABLE IF NOT EXISTS habits (
    id SERIAL PRIMARY KEY,
    state habit_state NOT NULL DEFAULT 'no_info',
    "date" DATE NOT NULL,
    "type" habit_type NOT NULL,
    origin TEXT NOT NULL DEFAULT 'unknown',
    UNIQUE("date", "type")
);

