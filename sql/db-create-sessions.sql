CREATE TABLE snippetbox.sessions
(
    token  TEXT PRIMARY KEY,
    data   BYTEA       NOT NULL,
    expiry TIMESTAMPTZ NOT NULL
);

CREATE INDEX sessions_expiry_idx ON snippetbox.sessions (expiry);