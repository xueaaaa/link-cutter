CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    username VARCHAR(16) UNIQUE NOT NULL,
    password VARCHAR(60) NOT NULL,
    creationDate TIMESTAMP NOT NULL,
    lastAccessDate TIMESTAMP
);