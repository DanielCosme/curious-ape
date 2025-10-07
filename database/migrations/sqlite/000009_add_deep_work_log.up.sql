CREATE TABLE IF NOT EXISTS deep_work_log (
  id                  INTEGER PRIMARY KEY,
  day_id              INTEGER NOT NULL,

  title               TEXT NOT NULL,
  start_time          DATE NOT NULL,
  end_time            DATE NOT NULL,
  note                TEXT NOT NULL DEFAULT "",

  origin              TEXT NOT NULL CHECK (length(origin) > 1),
  raw                 TEXT,

  FOREIGN KEY (day_id) REFERENCES "day" (id) ON DELETE CASCADE,
  UNIQUE (day_id, start_time)
);
