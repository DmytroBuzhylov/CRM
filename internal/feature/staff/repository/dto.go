package repository

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type BreakSchedule struct {
	BreakBegin time.Time `json:"break_begin"`
	BreakEnd   time.Time `json:"break_end"`
}

type WorkDaySchedule struct {
	WorkDayID  uuid.UUID       `json:"work_day_id"`
	DayOfWeak  int             `json:"day_of_week"`
	ShiftStart time.Time       `json:"shift_start"`
	ShiftEnd   time.Time       `json:"shift_end"`
	Breaks     []BreakSchedule `json:"breaks"`
}

type StaffWithSchedule struct {
	ID           uuid.UUID      `db:"staff_id"`
	UserID       uuid.UUID      `db:"user_id"`
	Position     sql.NullString `db:"position"`
	RoleName     string         `db:"role_name"`
	WorkSchedule []byte         `db:"work_schedule"`
}
