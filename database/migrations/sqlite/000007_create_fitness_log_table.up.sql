CREATE TABLE IF NOT EXISTS fitness_log (
  id                  INTEGER PRIMARY KEY,
  day_id              INTEGER NOT NULL,

  title               TEXT NOT NULL,
  start_time          DATE NOT NULL,
  end_time            DATE NOT NULL,
  note                TEXT NOT NULL DEFAULT "",
  
  type                TEXT NOT NULL,
  
  origin              TEXT NOT NULL CHECK (length(origin) > 1),
  raw                 TEXT,

  FOREIGN KEY (day_id) REFERENCES "day" (id) ON DELETE CASCADE,
  UNIQUE (day_id, start_time),
  CHECK (type IN ('strength', 'cardio'))
);

