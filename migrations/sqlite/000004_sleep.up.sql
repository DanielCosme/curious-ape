CREATE TABLE IF NOT EXISTS sleep_logs (
    id                  INTEGER primary key,
    day_id              INTEGER not null,


    FOREIGN KEY (day_id) REFERENCES "days" (id) ON DELETE CASCADE
);
