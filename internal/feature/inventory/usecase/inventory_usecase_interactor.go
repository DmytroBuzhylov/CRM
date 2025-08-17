package usecase

import (
	"Test/internal/feature/inventory/entity"
	"Test/internal/feature/inventory/interface_adapter/dto"
	"Test/internal/feature/inventory/repository"
	"Test/pkg/storage"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type InventoryUseCaseInteractor struct {
	repo    repository.InventoryRepository
	storage storage.Storage
}

func NewInventoryUseCaseInteractor(repo repository.InventoryRepository, storage storage.Storage) *InventoryUseCaseInteractor {
	return &InventoryUseCaseInteractor{repo: repo, storage: storage}
}

func (uc *InventoryUseCaseInteractor) CreateIngredient(ctx context.Context, req dto.CreateIngredientRequest, fileOptions storage.StorageFileOptions) (dto.CreateIngredientResponse, error) {
	newID := uuid.New()
	newIngredient := entity.Ingredient{
		ID:              newID,
		OrganizationID:  req.OrganizationID,
		Name:            req.Name,
		Unit:            req.Unit,
		Quantity:        req.Quantity,
		MinimumQuantity: req.MinimumQuantity,
	}
	err := uc.storage.UploadFile(ctx, newID.String(), fileOptions.File, fileOptions.FileSize, fileOptions.ContentType)
	if err != nil {
		return dto.CreateIngredientResponse{}, errors.New("failed to save ingredient image")
	}

	err = uc.repo.Create(ctx, newIngredient)
	if err != nil {
		uc.storage.DeleteFile(ctx, newID.String())
		return dto.CreateIngredientResponse{}, errors.New("failed to save ingredient")
	}

	return dto.CreateIngredientResponse{Status: "ok"}, nil
}

func (uc *InventoryUseCaseInteractor) GetIngredient(ctx context.Context, req dto.GetIngredientRequest) (dto.GetIngredientResponse, error) {
	ingredient, err := uc.repo.Get(ctx, req.ID, req.OrganizationID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return dto.GetIngredientResponse{}, errors.New("ingredient with the specified id does not exist")
		} else {
			return dto.GetIngredientResponse{}, errors.New("unable to find ingredient with specified ID")
		}
	}

	return dto.GetIngredientResponse{
		ID:              ingredient.ID,
		OrganizationID:  ingredient.OrganizationID,
		Name:            ingredient.Name,
		Unit:            ingredient.Unit,
		Quantity:        ingredient.Quantity,
		MinimumQuantity: ingredient.MinimumQuantity,
		CreatedAt:       ingredient.CreatedAt,
		UpdatedAt:       ingredient.UpdatedAt,
	}, nil
}

func (uc *InventoryUseCaseInteractor) GetAllIngredients(ctx context.Context, req dto.GetAllIngredientsRequest) ([]dto.GetIngredientResponse, error) {
	ingredients, err := uc.repo.GetAll(ctx, req.OrganizationID)
	if err != nil {
		return []dto.GetIngredientResponse{}, nil
	}

	log.Info().Msg(fmt.Sprintf("%v", ingredients))

	return []dto.GetIngredientResponse{}, nil
}

func (uc *InventoryUseCaseInteractor) DeleteIngredient(ctx context.Context, req dto.DeleteIngredientRequest) (dto.DeleteIngredientResponse, error) {
	err := uc.repo.Delete(ctx, req.ID, req.OrganizationID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return dto.DeleteIngredientResponse{}, errors.New("this ingredient does not exist")
		}
		return dto.DeleteIngredientResponse{}, errors.New("error deleting ingredient")
	}

	err = uc.storage.DeleteFile(ctx, req.ID.String())
	if err != nil {
		return dto.DeleteIngredientResponse{}, errors.New("error deleting ingredient image")
	}

	return dto.DeleteIngredientResponse{Status: "ok"}, nil
}

func (uc *InventoryUseCaseInteractor) DeleteManyIngredients(ctx context.Context, req dto.DeleteManyIngredientsRequest) (dto.DeleteIngredientResponse, error) {
	err := uc.repo.DeleteMany(ctx, req.IDs, req.OrganizationID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return dto.DeleteIngredientResponse{}, errors.New("this ingredient does not exist")
		}
		return dto.DeleteIngredientResponse{}, errors.New("error deleting ingredient")
	}

	for _, i := range req.IDs {
		err = uc.storage.DeleteFile(ctx, i.String())
		if err != nil {
			log.Warn().Err(err).Send()
		}
	}

	return dto.DeleteIngredientResponse{Status: "ok"}, nil
}

func (uc *InventoryUseCaseInteractor) UpdateIngredient(ctx context.Context, req dto.UpdateIngredientRequest) (dto.UpdateIngredientResponse, error) {
	ingredient := entity.Ingredient{
		ID:              req.ID,
		OrganizationID:  req.OrganizationID,
		Name:            req.Name,
		Unit:            req.Unit,
		Quantity:        req.Quantity,
		MinimumQuantity: req.MinimumQuantity,
	}

	err := uc.repo.Update(ctx, ingredient)
	if err != nil {
		return dto.UpdateIngredientResponse{}, errors.New("")
	}

	return dto.UpdateIngredientResponse{Status: "ok"}, nil
}
