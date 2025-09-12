package usecase

import (
	"Test/internal/feature/menu/interface_adapter/dto"
	"Test/internal/feature/menu/repository"
	"Test/pkg/storage"
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
)

type MenuUseCaseInteractor struct {
	menuRepo repository.MenuRepository
	storage  storage.Storage
}

func NewMenuUseCaseInteractor(menuRepo repository.MenuRepository, storage storage.Storage) *MenuUseCaseInteractor {
	return &MenuUseCaseInteractor{menuRepo: menuRepo, storage: storage}
}

func (uc *MenuUseCaseInteractor) CreateMenuItem(ctx context.Context, item *dto.MenuItemWithRecipeDTO, fileOpts *storage.StorageFileOptions) (*dto.CreateMenuItemResponse, error) {
	item.MenuItem.ID = uuid.New()
	for _, i := range item.RecipeItems {
		i.ID = uuid.New()
		i.MenuItemID = item.MenuItem.ID
	}

	txRepo, err := uc.menuRepo.StartTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := txRepo.Rollback(ctx); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			log.Error().Err(err).Send()
		}
	}()

	if err := txRepo.CreateMenuItem(ctx, item.MenuItem); err != nil {
		return nil, err
	}

	if err := txRepo.AddRecipeItems(ctx, item.RecipeItems); err != nil {
		return nil, err
	}

	if err := uc.storage.UploadFile(ctx, item.MenuItem.ID.String(), fileOpts.File, fileOpts.FileSize, fileOpts.ContentType); err != nil {
		return nil, err
	}

	if err := txRepo.Commit(ctx); err != nil {
		log.Error().Err(err).Send()
	}

	return &dto.CreateMenuItemResponse{Status: "Ok"}, err
}

func (uc *MenuUseCaseInteractor) UpdateMenuItem(ctx context.Context, item *dto.MenuItemWithRecipeDTO) (*dto.UpdateMenuItemResponse, error) {

	tx, err := uc.menuRepo.StartTx(ctx, pgx.TxOptions{})
	if err != nil {
		log.Warn().Err(err).Msg("MenuUseCase@UpdateMenuItem")
		return nil, errors.New("failed to start transaction")
	}

	defer func() {
		if err = tx.Rollback(ctx); err != nil {
			log.Error().Err(err).Msg("MenuUseCase@UpdateMenuItem")
		}
	}()

	err = tx.UpdateRecipeItem(ctx, item.RecipeItems[0])
	if err != nil {
		return nil, err
	}

	err = tx.UpdateMenuItem(ctx, item.MenuItem)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(ctx); err != nil {
		log.Error().Err(err).Msg("MenuUseCase@UpdateMenuItem")
	}

	return &dto.UpdateMenuItemResponse{Status: "ok"}, nil
}

func (uc *MenuUseCaseInteractor) DeleteMenuItem(ctx context.Context, request *dto.DeleteMenuItemRequest) (*dto.DeleteMenuItemResponse, error) {
	return nil, nil
}

func (uc *MenuUseCaseInteractor) GetMenuItem(ctx context.Context, request *dto.GetMenuItemRequest) (*dto.MenuItemWithRecipeDTO, error) {
	return nil, nil
}

func (uc *MenuUseCaseInteractor) GetAllMenuItems(ctx context.Context, request *dto.GetAllMenuItemsRequest) ([]*dto.MenuItemWithRecipeDTO, error) {
	return nil, nil
}
