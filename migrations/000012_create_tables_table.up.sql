CREATE TABLE IF NOT EXISTS tables
(
    id              uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id uuid not null REFERENCES organizations (id),
    number INTEGER NOT NULL,
    status varchar(50) NOT NULL DEFAULT 'available',
    seats INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_tables_organization_number ON tables (organization_id, number);