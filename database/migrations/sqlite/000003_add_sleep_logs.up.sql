CREATE TABLE IF NOT EXISTS sleep_log (
  id                  INTEGER PRIMARY KEY,
  day_id              INTEGER NOT NULL,

  start_time          DATE NOT NULL,
  end_time            DATE NOT NULL,
  is_main_sleep       BOOLEAN DEFAULT TRUE,
  total_time_in_bed   INTEGER DEFAULT 0,
  minutes_asleep      INTEGER DEFAULT 0,

  origin              TEXT NOT NULL CHECK (length(origin) > 1),
  raw                 TEXT,

  FOREIGN KEY (day_id) REFERENCES "day" (id) ON DELETE CASCADE,
  UNIQUE (day_id, is_main_sleep)
);

