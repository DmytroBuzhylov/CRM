CREATE TABLE IF NOT EXISTS work_days
(
    id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    staff_id    uuid NOT NULL REFERENCES staff (id) ON DELETE CASCADE,
    day_of_week smallint NOT NULL CHECK (day_of_week BETWEEN 0 AND 6), -- 0=Sunday, 6=Saturday
    shift_start time NOT NULL,
    shift_end   time NOT NULL,
    created_at  timestamptz DEFAULT current_timestamp,
    updated_at  timestamptz DEFAULT current_timestamp,

    UNIQUE (staff_id, day_of_week)
);
