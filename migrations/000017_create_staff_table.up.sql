CREATE TABLE IF NOT EXISTS staff
(
    id              uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         uuid NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    organization_id uuid NOT NULL REFERENCES organizations (id) ON DELETE CASCADE,
    role_id         uuid REFERENCES roles (id) ON DELETE SET NULL,
    position        varchar(100),
    created_at      timestamptz DEFAULT current_timestamp,
    updated_at      timestamptz DEFAULT current_timestamp,

    UNIQUE (user_id, organization_id)
);