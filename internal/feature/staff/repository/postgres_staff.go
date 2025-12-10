package repository

import (
	"Test/internal/feature/staff/entity"
	"Test/internal/pkg/provider"
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

//go:embed sql/staff/get_by_id.sql
var getStaffByIDQuery string

//go:embed sql/staff/get_by_user_id.sql
var getStaffByUserIDQuery string

//go:embed sql/staff/get_by_organization_id.sql
var getStaffByOrganizationIDQuery string

//go:embed sql/staff/create.sql
var createStaffQuery string

//go:embed sql/staff/delete_by_id.sql
var deleteStaffByIDQuery string

type StaffRepo struct {
	connector provider.SQLConnector
	qb        sq.StatementBuilderType
}

func NewStaffRepo(connector provider.SQLConnector) *StaffRepo {
	return &StaffRepo{
		connector: connector,
		qb:        sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *StaffRepo) Create(ctx context.Context, staff *entity.Staff) error {
	_, err := r.connector.Exec(ctx, createStaffQuery,
		staff.ID,
		staff.UserID,
		staff.OrganizationID,
		staff.RoleID,
		staff.Position,
		staff.CreatedAt,
		staff.UpdatedAt,
	)
	return err
}

func (r *StaffRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.connector.Exec(ctx, deleteStaffByIDQuery, id)
	return err
}

func (r *StaffRepo) Get(ctx context.Context, id uuid.UUID) (*entity.Staff, error) {
	var staff entity.Staff

	err := r.connector.QueryRow(ctx, getStaffByIDQuery, id).Scan(
		&staff.ID,
		&staff.UserID,
		&staff.OrganizationID,
		&staff.RoleID,
		&staff.Position,
		&staff.CreatedAt,
		&staff.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &staff, nil
}

func (r *StaffRepo) GetByUserID(ctx context.Context, userID uuid.UUID) (*entity.Staff, error) {
	var staff entity.Staff

	err := r.connector.QueryRow(ctx, getStaffByUserIDQuery, userID).Scan(
		&staff.ID,
		&staff.UserID,
		&staff.OrganizationID,
		&staff.RoleID,
		&staff.Position,
		&staff.CreatedAt,
		&staff.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &staff, nil
}

func (r *StaffRepo) GetMany(ctx context.Context, organizationId uuid.UUID) ([]*entity.Staff, error) {
	rows, err := r.connector.Query(ctx, getStaffByOrganizationIDQuery, organizationId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var staffList []*entity.Staff
	for rows.Next() {
		var s entity.Staff
		err := rows.Scan(
			&s.ID,
			&s.UserID,
			&s.OrganizationID,
			&s.RoleID,
			&s.Position,
			&s.CreatedAt,
			&s.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		staffList = append(staffList, &s)
	}
	return staffList, rows.Err()
}

func (r *StaffRepo) GetRoleIDByUserID(ctx context.Context, userID uuid.UUID) (uuid.UUID, error) {
	query, arg, err := r.qb.Select("role_id").From("staff").Where(sq.Eq{
		"id": userID,
	}).ToSql()
	if err != nil {
		return [16]byte{}, err
	}

	var roleID uuid.UUID
	err = r.connector.QueryRow(ctx, query, arg...).Scan(&roleID)
	if err != nil {
		return [16]byte{}, err
	}
	return roleID, nil
}

//go:embed sql/staff/work_day.sql
var workDaySQL string

//go:embed sql/staff/work_day.sql
var breakSQL string

func (r *StaffRepo) SaveTimetable(ctx context.Context, staffID uuid.UUID, workDays []entity.WorkDay, breaksByWorkDayID map[uuid.UUID][]entity.Break) (err error) {

	workDayIDMap := make(map[uuid.UUID]uuid.UUID)

	for _, wd := range workDays {
		var realWorkDayID uuid.UUID

		err = r.connector.QueryRow(ctx, workDaySQL,
			wd.ID,
			staffID,
			int(wd.DayOfWeak),
			wd.ShiftStart.Format("15:04:05"),
			wd.ShiftEnd.Format("15:04:05"),
		).Scan(&realWorkDayID)

		if err != nil {
			return fmt.Errorf("insert/update workday failed: %w", err)
		}

		workDayIDMap[wd.ID] = realWorkDayID
	}

	var realWorkDayIDs []uuid.UUID
	for _, realID := range workDayIDMap {
		realWorkDayIDs = append(realWorkDayIDs, realID)
	}

	if len(realWorkDayIDs) > 0 {
		delQuery, delArgs, err := r.qb.Delete("breaks").Where(sq.Eq{"work_day_id": realWorkDayIDs}).ToSql()
		if err != nil {
			return fmt.Errorf("delete breaks query build failed: %w", err)
		}

		if _, err = r.connector.Exec(ctx, delQuery, delArgs...); err != nil {
			return fmt.Errorf("deleting old breaks failed: %w", err)
		}
	}

	for _, breaks := range breaksByWorkDayID {
		for _, b := range breaks {
			realWorkDayID := workDayIDMap[b.WorkDayID]

			if _, err = r.connector.Exec(ctx, breakSQL,
				uuid.New(),
				realWorkDayID,
				b.BreakBegin.Format("15:04:05"),
				b.BreakEnd.Format("15:04:05"),
			); err != nil {
				return fmt.Errorf("inserting break failed: %w", err)
			}
		}
	}

	return nil
}

func (r *StaffRepo) DeleteTimetable(ctx context.Context, staffID uuid.UUID, days []time.Weekday) error {

	var dayInts []int
	for _, d := range days {
		dayInts = append(dayInts, int(d))
	}

	selectQuery, selectArgs, _ := r.qb.Select("id").From("work_days").Where(sq.Eq{
		"staff_id":    staffID,
		"day_of_week": dayInts,
	}).ToSql()

	rows, err := r.connector.Query(ctx, selectQuery, selectArgs...)
	if err != nil {
		return fmt.Errorf("failed to query work day IDs: %w", err)
	}
	defer rows.Close()

	var workDayIDsToDelete []uuid.UUID
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return fmt.Errorf("failed to scan work day ID: %w", err)
		}
		workDayIDsToDelete = append(workDayIDsToDelete, id)
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("row iteration error: %w", err)
	}

	if len(workDayIDsToDelete) > 0 {
		deleteWDQuery, deleteWDArgs, _ := r.qb.Delete("work_days").Where(sq.Eq{"id": workDayIDsToDelete}).ToSql()
		if _, err = r.connector.Exec(ctx, deleteWDQuery, deleteWDArgs...); err != nil {
			return fmt.Errorf("failed to delete work days: %w", err)
		}
	}

	return nil
}

//go:embed sql/staff/get_staff_with_schedule.sql
var getStaffWithScheduleQuery string

func (r *StaffRepo) FindAvailableStaffByPermission(ctx context.Context, orgID uuid.UUID, permissionName string) ([]*StaffWithSchedule, error) {
	rows, err := r.connector.Query(ctx, getStaffWithScheduleQuery, orgID, permissionName)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var results []*StaffWithSchedule

	for rows.Next() {
		var s StaffWithSchedule

		err := rows.Scan(
			&s.ID,
			&s.UserID,
			&s.Position,
			&s.RoleName,
			&s.WorkSchedule,
		)
		if err != nil {
			return nil, fmt.Errorf("row scan error: %w", err)
		}

		var schedule []WorkDaySchedule
		if len(s.WorkSchedule) > 0 && string(s.WorkSchedule) != "[]" {
			if err := json.Unmarshal(s.WorkSchedule, &schedule); err != nil {
				return nil, fmt.Errorf("json unmarshal schedule error: %w", err)
			}
		}
		
		results = append(results, &s)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return results, nil
}
