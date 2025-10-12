CREATE TABLE IF NOT EXISTS roles
(
    id              uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id uuid NOT NULL REFERENCES organizations (id) ON DELETE CASCADE,
    name            varchar(100) NOT NULL,
    description     text,
    created_at      timestamptz DEFAULT current_timestamp,
    updated_at      timestamptz DEFAULT current_timestamp,

    UNIQUE (organization_id, name)
);