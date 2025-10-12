UPDATE work_days
SET day_of_week = $1, shift_start = $2, shift_end = $3, updated_at = $4
WHERE id = $5