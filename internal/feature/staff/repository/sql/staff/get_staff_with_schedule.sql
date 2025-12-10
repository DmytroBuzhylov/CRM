SELECT s.id AS staff_id, s.user_id, s.position, r.name AS role_name,
       COALESCE(
                       jsonb_agg(
                       DISTINCT jsonb_build_object(
                               'work_day_id', wd.id, 'day_of_week', wd.day_of_week,
                               'shift_start', wd.shift_start, 'shift_end', wd.shift_end,
                               'breaks', (SELECT COALESCE(jsonb_agg(jsonb_build_object('break_begin', b.break_begin, 'break_end', b.break_end)), '[]'::jsonb) FROM breaks b WHERE b.work_day_id = wd.id)
                                )
                                ) FILTER (WHERE wd.id IS NOT NULL), '[]'::jsonb
       ) AS work_schedule
FROM staff s
         JOIN roles r ON s.role_id = r.id
         JOIN role_permissions rp ON rp.role_id = r.id
         JOIN permissions p ON p.id = rp.permission_id
         LEFT JOIN work_days wd ON wd.staff_id = s.id
WHERE s.organization_id = $1 AND p.name = $2
GROUP BY s.id, s.user_id, s.position, r.name
ORDER BY s.id;