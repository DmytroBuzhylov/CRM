CREATE TABLE IF NOT EXISTS role_permissions
(
    role_id       uuid NOT NULL REFERENCES roles (id) ON DELETE CASCADE,
    permission_id uuid NOT NULL REFERENCES permissions (id) ON DELETE CASCADE,
    created_at    timestamptz DEFAULT current_timestamp,

    PRIMARY KEY (role_id, permission_id)
);