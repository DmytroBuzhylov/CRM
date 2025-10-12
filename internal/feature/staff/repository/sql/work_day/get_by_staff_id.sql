SELECT id, staff_id, day_of_week, shift_start, shift_end, created_at, updated_at
FROM work_days
WHERE staff_id = $1