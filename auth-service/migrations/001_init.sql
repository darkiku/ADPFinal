CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS users (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email         VARCHAR(255) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    name          VARCHAR(255) NOT NULL,
    role          VARCHAR(50) DEFAULT 'student',
    verified      BOOLEAN DEFAULT false,
    verify_token  TEXT DEFAULT '',
    created_at    TIMESTAMP DEFAULT NOW()
);
