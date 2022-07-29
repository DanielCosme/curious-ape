CREATE TABLE IF NOT EXISTS oauths (
    id                  INTEGER primary key,

    provider            TEXT not null UNIQUE,
    access_token        TEXT not null UNIQUE,
    refresh_token       TEXT,
    type                TEXT,
    expiration          DATE,

    toggl_workspace_id      INTEGER default "",
    toggl_organization_id   INTEGER default "",
    toggl_project_ids       TEXT defatult ""
)
