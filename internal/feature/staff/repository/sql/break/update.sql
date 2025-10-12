UPDATE breaks
SET break_begin = $1, break_end = $2, updated_at = $3
WHERE id = $4;
