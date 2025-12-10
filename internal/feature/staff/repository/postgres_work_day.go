package repository

import (
	"Test/internal/feature/staff/entity"
	"Test/internal/pkg/provider"
	"context"
	_ "embed"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"time"
)

//go:embed sql/work_day/create.sql
var createWorkDayQuery string

//go:embed sql/work_day/get_by_id.sql
var getByIDQuery string

//go:embed sql/work_day/get_by_staff_id.sql
var getByStaffIDQuery string

//go:embed sql/work_day/update.sql
var updateWorkDayQuery string

//go:embed sql/work_day/delete.sql
var deleteWorkDayQuery string

type WorkDayRepo struct {
	connector provider.SQLConnector
	qb        sq.StatementBuilderType
}

func NewWorkDayRepo(connector provider.SQLConnector) *WorkDayRepo {
	return &WorkDayRepo{
		connector: connector,
		qb:        sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *WorkDayRepo) Create(ctx context.Context, wD *entity.WorkDay) error {
	_, err := r.connector.Exec(ctx, createWorkDayQuery,
		wD.ID,
		wD.StaffID,
		int(wD.DayOfWeak),
		wD.ShiftStart,
		wD.ShiftEnd,
		wD.CreatedAt,
		wD.UpdatedAt,
	)
	return err
}

func (r *WorkDayRepo) Get(ctx context.Context, id uuid.UUID) (*entity.WorkDay, error) {
	var wD = new(entity.WorkDay)
	var dayOfWeek int

	err := r.connector.QueryRow(ctx, getByIDQuery).Scan(
		&wD.ID,
		&wD.StaffID,
		&dayOfWeek,
		&wD.ShiftStart,
		&wD.ShiftEnd,
		&wD.CreatedAt,
		&wD.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	wD.DayOfWeak = time.Weekday(dayOfWeek)

	return wD, nil
}

func (r *WorkDayRepo) GetMany(ctx context.Context, staffId uuid.UUID) ([]*entity.WorkDay, error) {

	rows, err := r.connector.Query(ctx, getByStaffIDQuery, staffId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var workDays []*entity.WorkDay
	for rows.Next() {
		var wd entity.WorkDay
		var dayOfWeek int
		err := rows.Scan(
			&wd.ID,
			&wd.StaffID,
			&dayOfWeek,
			&wd.ShiftStart,
			&wd.ShiftEnd,
			&wd.CreatedAt,
			&wd.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		wd.DayOfWeak = time.Weekday(dayOfWeek)
		workDays = append(workDays, &wd)
	}
	return workDays, rows.Err()
}

func (r *WorkDayRepo) Update(ctx context.Context, wD *entity.WorkDay) error {
	_, err := r.connector.Exec(ctx, updateWorkDayQuery,
		int(wD.DayOfWeak),
		wD.ShiftStart,
		wD.ShiftEnd,
		wD.UpdatedAt,
		wD.ID,
	)
	return err
}

func (r *WorkDayRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.connector.Exec(ctx, deleteWorkDayQuery, id)
	return err
}
