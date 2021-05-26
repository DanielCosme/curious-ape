CREATE TABLE IF NOT EXISTS food_habits (
    id serial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    state boolean NOT NULL,
    "date" date UNIQUE NOT NULL DEFAULT current_date,
    tags text[] NOT NULL
);


ALTER TABLE food_habits ADD CONSTRAINT tags_length_check 
CHECK (array_length(tags, 1) between 1 and 5);
-- TODO add constraint that only valid tags are passed on, enum.
