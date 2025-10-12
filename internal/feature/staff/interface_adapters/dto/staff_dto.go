package dto

import (
	"github.com/google/uuid"
	"time"
)

type InviteStaffReq struct {
	OrganizationID uuid.UUID `json:"-"`
	InvitedEmail   string    `json:"invited_email,omitempty"`
	InviterUserId  uuid.UUID `json:"inviter_user_id,omitempty"`
	InvitationCode string    `json:"invitation_code,omitempty"`
	Position       string    `json:"position,omitempty"`
	RoleID         uuid.UUID `json:"role,omitempty"`
}

type AcceptInviteReq struct {
	UserID uuid.UUID `json:"-"`
	Code   string    `json:"-"`
}

type UpdateStaffReq struct {
	ID     uuid.UUID `json:"id"`
	RoleID uuid.UUID `json:"role_id"`
}

type AddTimetableStaffReq struct {
	StaffID  uuid.UUID    `json:"staff_id,omitempty"`
	WorkDays []WorkDayDto `json:"work_days,omitempty"`
	Breaks   []BreakDto   `json:"breaks,omitempty"`
}

type UpdateTimetableStaffReq struct {
	StaffID  uuid.UUID    `json:"staff_id,omitempty"`
	WorkDays []WorkDayDto `json:"work_days,omitempty"`
	Breaks   []BreakDto   `json:"breaks,omitempty"`
}

type WorkDayDto struct {
	StaffID    uuid.UUID    `json:"staff_id,omitempty"`
	DayOfWeak  time.Weekday `json:"day_of_weak,omitempty"`
	ShiftStart time.Time    `json:"shift_start"`
	ShiftEnd   time.Time    `json:"shift_end"`
}

type BreakDto struct {
	WorkDayID  uuid.UUID    `json:"work_day_id,omitempty"`
	BreakBegin time.Time    `json:"break_begin"`
	BreakEnd   time.Time    `json:"break_end"`
	DayOfWeak  time.Weekday `json:"day_of_weak,omitempty"`
}

type DeleteTimetableStaffReq struct {
	ID   uuid.UUID      `json:"id,omitempty"`
	Days []time.Weekday `json:"days,omitempty"`
}
