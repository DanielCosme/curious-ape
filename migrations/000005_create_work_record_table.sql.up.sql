CREATE TABLE work_records (
   id SERIAL PRIMARY KEY,
   "date" DATE NOT NULL UNIQUE,
   grand_total INT NOT NULL CHECK(grand_total >= 0),
   raw_json JSONB NOT NULL
);
