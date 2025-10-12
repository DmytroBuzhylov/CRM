CREATE TABLE IF NOT EXISTS staff_invitations
(
    id              uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id uuid NOT NULL REFERENCES organizations (id) ON DELETE CASCADE,
    invited_email   varchar(255) NOT NULL,
    inviter_user_id uuid NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    invitation_code varchar(255) UNIQUE NOT NULL,
    position        varchar(100),
    role            varchar(50),
    expires_at      timestamp with time zone NOT NULL,
    status          varchar(50) NOT NULL DEFAULT 'pending' CHECK (status IN('pending', 'accepted', 'rejected', 'expired')),
    created_at      timestamp with time zone NOT NULL DEFAULT now()
);

CREATE INDEX idx_staff_invitations_code ON staff_invitations (invitation_code);
CREATE INDEX idx_staff_invitations_email ON staff_invitations (invited_email);