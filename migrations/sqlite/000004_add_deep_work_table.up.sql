CREATE TABLE IF NOT EXISTS deep_work_logs (
    id                  INTEGER PRIMARY KEY,
    day_id              INTEGER NOT NULL,

    "date"              DATE NOT NULL,
    seconds             INTEGER NOT NULL CHECK (length(seconds) < 60),
    is_automated        BOOLEAN DEFAULT false,
    origin              TEXT NOT NULL CHECK (length(origin) > 1),

    UNIQUE (day_id, origin),
    FOREIGN KEY (day_id) REFERENCES "days" (id) ON DELETE CASCADE
);
