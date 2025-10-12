package usecase

import (
	"Test/internal/feature/table/entity"
	"Test/internal/feature/table/interface_adapters/dto"
	"Test/internal/feature/table/repository"
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
)

type TableUseCaseInteractor struct {
	repo repository.TableRepository
}

func NewTableUseCaseInteractor(repo repository.TableRepository) *TableUseCaseInteractor {
	return &TableUseCaseInteractor{repo: repo}
}

func (uc *TableUseCaseInteractor) CreateTable(ctx context.Context, req dto.CreateTableRequest) error {
	table := &entity.Table{
		OrganizationID: req.OrganizationID,
		Number:         req.Number,
		Seats:          req.Seats,
	}

	err := uc.repo.Create(ctx, table)
	if err != nil {
		log.Warn().Err(err).
			Str("usecase", "Table").
			Str("method", "CreateTable").
			Msg("failed to create table")
		return errors.New("failed to create table")
	}
	return nil
}

func (uc *TableUseCaseInteractor) UpdateTableSeats(ctx context.Context, req dto.UpdateSeatsRequest) error {

	err := uc.repo.UpdateSeats(ctx, req.Seats, req.ID, req.OrganizationID)
	if err != nil {
		log.Warn().Err(err).
			Str("usecase", "Table").
			Str("method", "UpdateTableSeats").
			Msg("failed to update table seats")
		return errors.New("failed to update table seats")
	}
	return nil
}

func (uc *TableUseCaseInteractor) UpdateTableNumber(ctx context.Context, req dto.UpdateNumberRequest) error {

	err := uc.repo.UpdateNumber(ctx, req.Number, req.OrganizationID, req.ID)
	if err != nil {
		log.Warn().Err(err).
			Str("usecase", "Table").
			Str("method", "UpdateTableNumber").
			Msg("failed to update table number")
		return errors.New("failed to update table number")
	}
	return nil
}

func (uc *TableUseCaseInteractor) UpdateTableStatus(ctx context.Context, req dto.UpdateStatusRequest) error {

	err := uc.repo.UpdateStatus(ctx, req.Status, req.ID, req.OrganizationID)
	if err != nil {
		log.Warn().Err(err).
			Str("usecase", "Table").
			Str("method", "UpdateTableStatus").
			Msg("failed to update table status")
		return errors.New("failed to update table status")
	}
	return nil
}

func (uc *TableUseCaseInteractor) GetTable(ctx context.Context, req dto.GetTableRequest) (*dto.GetTableResponse, error) {

	resp, err := uc.repo.Get(ctx, req.ID)
	if err != nil {
		if errors.As(err, pgx.ErrNoRows) {
			return nil, errors.New("table with such id does not exist")
		}
		log.Warn().Err(err).
			Str("usecase", "Table").
			Str("method", "GetTable").
			Msg("failed to get table")
		return nil, errors.New("failed to get table")
	}
	return &dto.GetTableResponse{
		ID:             resp.ID,
		OrganizationID: resp.OrganizationID,
		Number:         resp.Number,
		Status:         resp.Status,
		Seats:          resp.Seats,
		CreatedAt:      resp.CreatedAt,
		UpdatedAt:      resp.UpdatedAt,
	}, err
}

func (uc *TableUseCaseInteractor) GetOrganizationTables(ctx context.Context, req dto.GetOrganizationTablesRequest) ([]*dto.GetTableResponse, error) {

	tables, err := uc.repo.GetAll(ctx, req.OrganizationID)
	if err != nil {
		if errors.As(err, pgx.ErrNoRows) {
			return nil, errors.New("table with such organizational identifier does not exist")
		}
		log.Warn().Err(err).
			Str("usecase", "Table").
			Str("method", "GetOrganizationTables").
			Msg("failed to get organization tables")
		return nil, errors.New("failed to get organization tables")
	}

	dtoResp := make([]*dto.GetTableResponse, 0, len(tables))
	for _, i := range tables {
		dtoResp = append(dtoResp, &dto.GetTableResponse{
			ID:             i.ID,
			OrganizationID: i.OrganizationID,
			Number:         i.Number,
			Status:         i.Status,
			Seats:          i.Seats,
			CreatedAt:      i.CreatedAt,
			UpdatedAt:      i.UpdatedAt,
		})
	}

	return dtoResp, nil
}

func (uc *TableUseCaseInteractor) DeleteTable(ctx context.Context, req dto.DeleteRequest) error {

	err := uc.repo.Delete(ctx, req.ID, req.OrganizationID)
	if err != nil {
		log.Warn().Err(err).
			Str("usecase", "Table").
			Str("method", "DeleteTable").
			Msg("failed to delete table")
		return errors.New("failed to delete table")
	}
	return nil
}
