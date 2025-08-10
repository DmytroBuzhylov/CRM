CREATE TABLE IF NOT EXISTS menu_items
(
    id              uuid PRIMARY KEY                  default gen_random_uuid(),
    organization_id uuid,
    name            varchar(255)             not null,
    description     text,
    price           bigint,
    category        varchar(255),
    is_available    bool,
    updated_at      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (organization_id) REFERENCES organizations (id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_id ON menu_items (id);
CREATE INDEX IF NOT EXISTS idx_name ON menu_items (name);