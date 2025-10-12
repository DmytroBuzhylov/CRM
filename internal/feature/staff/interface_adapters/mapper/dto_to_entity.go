package mapper

import (
	"Test/internal/feature/staff/entity"
	"Test/internal/feature/staff/interface_adapters/dto"
	"github.com/google/uuid"
)

func MapAddTimetableReqToEntities(req *dto.AddTimetableStaffReq) ([]entity.WorkDay, map[uuid.UUID][]entity.Break) {
	workDays := make([]entity.WorkDay, len(req.WorkDays))
	for idx, wdDto := range req.WorkDays {
		workDays[idx] = entity.WorkDay{
			ID:         wdDto.StaffID,
			DayOfWeak:  wdDto.DayOfWeak,
			ShiftStart: wdDto.ShiftStart,
			ShiftEnd:   wdDto.ShiftEnd,
		}
		workDays[idx].ID = uuid.New()
	}

	breaksByWorkDayID := make(map[uuid.UUID][]entity.Break)
	for _, bDto := range req.Breaks {
		workDayID := bDto.WorkDayID

		breaksByWorkDayID[workDayID] = append(breaksByWorkDayID[workDayID], entity.Break{
			WorkDayID:  workDayID,
			BreakBegin: bDto.BreakBegin,
			BreakEnd:   bDto.BreakEnd,
		})
	}

	return workDays, breaksByWorkDayID
}
