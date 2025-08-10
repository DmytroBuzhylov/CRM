package usecase

import (
	"Test/internal/feature/inventory/interface_adapter/dto"
	"Test/internal/pkg/storage"
	"context"
)

type InventoryUseCase interface {
	CreateIngredient(ctx context.Context, req dto.CreateIngredientRequest, fileOptions storage.StorageFileOptions) (dto.CreateIngredientResponse, error)
	GetIngredient(ctx context.Context, req dto.GetIngredientRequest) (dto.GetIngredientResponse, error)
	GetAllIngredients(ctx context.Context, req dto.GetAllIngredientsRequest) ([]dto.GetIngredientResponse, error)
	DeleteIngredient(ctx context.Context, req dto.DeleteIngredientRequest) (dto.DeleteIngredientResponse, error)
	DeleteManyIngredients(ctx context.Context, req dto.DeleteManyIngredientsRequest) (dto.DeleteIngredientResponse, error)
	UpdateIngredient(ctx context.Context, req dto.UpdateIngredientRequest) (dto.UpdateIngredientResponse, error)
}
