CREATE TABLE IF NOT EXISTS users
(
    id              uuid PRIMARY KEY         DEFAULT gen_random_uuid(),
    full_name       VARCHAR(255) NOT NULL,
    username        VARCHAR(255) NOT NULL UNIQUE,
    hashed_password TEXT         NOT NULL,
    email           VARCHAR(255) UNIQUE,
    phone           VARCHAR(50),
    role            VARCHAR(50),
    updated_at      TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_at      TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_users_email ON users (email);