SELECT id, work_day_id, break_begin, break_end, created_at, updated_at
FROM breaks
WHERE work_day_id = $1;
