CREATE TABLE IF NOT EXISTS
    users (
        "guid" BYTEA,
        -- "email" TEXT NOT NULL UNIQUE,
        "refresh_token" BYTEA
    );