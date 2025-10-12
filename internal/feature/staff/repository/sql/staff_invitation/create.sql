INSERT INTO staff_invitations (
    id, organization_id, invited_email, inviter_user_id, invitation_code,
    position, role, expires_at, status, created_at
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)