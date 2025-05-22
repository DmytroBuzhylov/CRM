CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    full_name VARCHAR(255) NOT NULL,
    username VARCHAR(255) NOT NULL,
    hashed_password TEXT NOT NULL,
    email VARCHAR(255),
    phone VARCHAR(10),
    role VARCHAR(10),
    organization_id INTEGER NULL default null,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_users_email ON users (email);