CREATE TABLE auth_tokens (
    service VARCHAR(50) PRIMARY KEY NOT NULL,
    access_token TEXT UNIQUE NOT NULL,
    refresh_token TEXT UNIQUE NOT NULL 
);

INSERT INTO auth_tokens (service, access_token, refresh_token)
VALUES ('fitbit', '', '')
VALUES ('google', '', '');
