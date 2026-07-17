CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE,
    username VARCHAR(16) UNIQUE,
    password VARCHAR(60),
    creationDate TIMESTAMP NOT NULL,
    lastAccessDate TIMESTAMP
)