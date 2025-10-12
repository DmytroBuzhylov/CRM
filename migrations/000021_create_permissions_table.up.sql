CREATE TABLE IF NOT EXISTS permissions
(
    id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name        varchar(100) UNIQUE NOT NULL,
    description text,
    created_at  timestamptz DEFAULT current_timestamp
);