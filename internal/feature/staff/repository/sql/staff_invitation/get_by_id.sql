SELECT (organization_id, invited_email, inviter_user_id, invitation_code,
        position, role, expires_at, status, created_at)
FROM staff_invitations
WHERE id = $1