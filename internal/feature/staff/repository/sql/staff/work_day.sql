INSERT INTO work_days (id, staff_id, day_of_week, shift_start, shift_end)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (staff_id, day_of_week)
    DO UPDATE SET shift_start = EXCLUDED.shift_start, shift_end = EXCLUDED.shift_end, updated_at = NOW()
RETURNING id