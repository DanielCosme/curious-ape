-- schema
CREATE TABLE habits (
    id              INTEGER primary key,
    uuid            TEXT not null,
    state           TEXT not null,
    time            DATETIME not null,
    creation_time   DATETIME not null,
    update_time     DATETIME not null
);
