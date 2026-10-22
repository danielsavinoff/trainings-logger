CREATE TABLE identities (
    id TEXT PRIMARY KEY NOT NULL,
    provider TEXT NOT NULL CHECK (provider IN ('Dropbox'))
);

CREATE TABLE identity_credentials (
    identity_id TEXT PRIMARY KEY NOT NULL REFERENCES identities(id) ON DELETE CASCADE,
    access_token_encrypted TEXT NOT NULL,
    refresh_token_encrypted TEXT NOT NULL,
    access_token_expires_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);
