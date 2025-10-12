package repository

import (
	"Test/internal/feature/staff/entity"
	"Test/internal/pkg/provider"
	"context"
	_ "embed"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

//go:embed sql/break/create.sql
var createBreakQuery string

//go:embed sql/break/get_by_id.sql
var getBreakByIDQuery string

//go:embed sql/break/get_many.sql
var getManyBreaksQuery string

//go:embed sql/break/update.sql
var updateBreakQuery string

//go:embed sql/break/delete.sql
var deleteBreakQuery string

type BreakRepo struct {
	provider provider.DataProvider
	qb       sq.StatementBuilderType
}

func NewBreakRepo(provider provider.DataProvider) *BreakRepo {
	return &BreakRepo{
		provider: provider,
		qb:       sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *BreakRepo) Create(ctx context.Context, br *entity.Break) error {
	_, err := r.provider.Exec(ctx, createBreakQuery,
		br.ID,
		br.WorkDayID,
		br.BreakBegin,
		br.BreakEnd,
		br.CreatedAt,
		br.UpdatedAt,
	)
	return err
}

func (r *BreakRepo) Get(ctx context.Context, id uuid.UUID) (*entity.Break, error) {
	var br entity.Break
	err := r.provider.QueryRow(ctx, getBreakByIDQuery, id).Scan(
		&br.ID,
		&br.WorkDayID,
		&br.BreakBegin,
		&br.BreakEnd,
		&br.CreatedAt,
		&br.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &br, nil
}

func (r *BreakRepo) GetMany(ctx context.Context, workDayId uuid.UUID) ([]*entity.Break, error) {
	rows, err := r.provider.Query(ctx, getManyBreaksQuery, workDayId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var breaks []*entity.Break
	for rows.Next() {
		var b entity.Break
		err := rows.Scan(
			&b.ID,
			&b.WorkDayID,
			&b.BreakBegin,
			&b.BreakEnd,
			&b.CreatedAt,
			&b.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		breaks = append(breaks, &b)
	}
	return breaks, rows.Err()
}

func (r *BreakRepo) Update(ctx context.Context, br *entity.Break) error {
	_, err := r.provider.Exec(ctx, updateBreakQuery,
		br.BreakBegin,
		br.BreakEnd,
		br.UpdatedAt,
		br.ID,
	)
	return err
}

func (r *BreakRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.provider.Exec(ctx, deleteBreakQuery, id)
	return err
}
