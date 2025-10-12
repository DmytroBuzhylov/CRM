SELECT id, user_id, organization_id, role_id, position, created_at, updated_at
FROM staff
WHERE user_id = $1