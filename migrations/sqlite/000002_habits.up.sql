CREATE TABLE IF NOT EXISTS habit_categories (
    id                  INTEGER primary key,

    name                TEXT not null,
    type                TEXT not null,
    code                TEXT not null default "default",
    description         TEXT default "",
    color               INTEGER default "#ffffff"
);

INSERT INTO habit_categories (name, type)
VALUES  ("Eat healthy", "food"),
        ("Wake up early", "wake_up"),
        ("Workout", "fitness"),
        ("Deep work", "deep_work");

CREATE TABLE IF NOT EXISTS habits (
    id                  INTEGER primary key,
    day_id              INTEGER not null,
    habit_category_id   INTEGER not null,

    success             BOOLEAN default false,
    origin              TEXT not null default "unknown",
    is_automated        BOOLEAN not null default false,
    note                TEXT default "",

    FOREIGN KEY (habit_category_id) REFERENCES habit_categories (id) ON DELETE CASCADE,
    FOREIGN KEY (day_id) REFERENCES "days" (id) ON DELETE CASCADE,
    UNIQUE (day_id, habit_category_id)
);
