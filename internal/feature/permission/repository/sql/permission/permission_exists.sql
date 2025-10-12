SELECT EXISTS (
    SELECT 1
    FROM staff s
             JOIN roles r ON s.role_id = r.id
             JOIN permissions p ON p.name = $3
             JOIN role_permissions rp ON rp.role_id = s.role_id AND rp.permission_id = p.id

    WHERE s.user_id = $1
      AND s.organization_id = $2
      AND r.organization_id = $2
      AND s.role_id IS NOT NULL
);