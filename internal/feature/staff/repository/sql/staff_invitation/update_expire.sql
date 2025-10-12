UPDATE staff_invitations
SET status = 'expired'
WHERE id = $1 AND status = 'pending'