CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE links (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    shortId VARCHAR(8) UNIQUE,
    origin VARCHAR(4096),
    creationDate TIMESTAMP NOT NULL,
    lastAccessDate TIMESTAMP
);

