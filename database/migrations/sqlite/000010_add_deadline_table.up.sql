CREATE TABLE IF NOT EXISTS deadline (
        id                  INTEGER PRIMARY KEY,

        title               VARCHAR(200) NOT NULL,
        start_time          DATE NOT NULL,
        end_time            DATE NOT NULL,
        recurring						BOOLEAN NOT NULL DEFAULT false
);
