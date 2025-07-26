CREATE TABLE IF NOT EXISTS invitations (
    id SERIAL PRIMARY KEY,
    organization_id BIGINT NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    invited_email VARCHAR(255) NOT NULL,
    invitation_code VARCHAR(255) UNIQUE NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_invitations_code ON invitations (invitation_code);
CREATE INDEX idx_invitations_email ON invitations (invited_email);

