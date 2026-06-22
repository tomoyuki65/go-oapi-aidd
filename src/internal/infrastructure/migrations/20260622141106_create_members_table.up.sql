SET statement_timeout = 0;

--bun:split

SELECT 1

--bun:split

SELECT 2
CREATE TABLE members (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    rank TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
