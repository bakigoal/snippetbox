CREATE TABLE snippetbox.users
(
    id              SERIAL PRIMARY KEY,
    name            VARCHAR(255) NOT NULL,
    email           VARCHAR(255) NOT NULL,
    hashed_password CHAR(60)     NOT NULL,
    created         timestamp    NOT NULL
);

ALTER TABLE snippetbox.users
    ADD CONSTRAINT users_uc_email UNIQUE (email);