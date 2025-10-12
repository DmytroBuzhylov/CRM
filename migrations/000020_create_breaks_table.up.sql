CREATE TABLE IF NOT EXISTS breaks
(
    id           uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    work_day_id  uuid NOT NULL REFERENCES work_days (id) ON DELETE CASCADE,
    break_begin  time NOT NULL,
    break_end    time NOT NULL,
    created_at   timestamptz DEFAULT current_timestamp,
    updated_at   timestamptz DEFAULT current_timestamp
);
