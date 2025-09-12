package usecase

import (
	"Test/internal/feature/table/interface_adapters/dto"
	"context"
)

type TableUseCase interface {
	CreateTable(ctx context.Context, req dto.CreateTableRequest) error
	UpdateTableSeats(ctx context.Context, req dto.UpdateSeatsRequest) error
	UpdateTableNumber(ctx context.Context, req dto.UpdateNumberRequest) error
	UpdateTableStatus(ctx context.Context, req dto.UpdateStatusRequest) error
	GetTable(ctx context.Context, req dto.GetTableRequest) (dto.GetTableResponse, error)
	GetOrganizationTables(ctx context.Context, req dto.GetOrganizationTablesRequest) ([]dto.GetTableResponse, error)
	DeleteTable(ctx context.Context, req dto.DeleteRequest) error
}
