package usecase

import (
	"Test/internal/feature/menu/interface_adapter/dto"
	"Test/internal/feature/menu/repository"
	"Test/pkg/transactor"
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
)

type MenuUseCaseInteractor struct {
	menuRepo   repository.MenuRepository
	menuTxRepo repository.MenuTransactionalRepository
	transactor transactor.Transactor
}

func NewMenuUseCaseInteractor(menuRepo repository.MenuRepository, menuTxRepo repository.MenuTransactionalRepository, transactor transactor.Transactor) *MenuUseCaseInteractor {
	return &MenuUseCaseInteractor{menuRepo: menuRepo, menuTxRepo: menuTxRepo, transactor: transactor}
}

func (uc *MenuUseCaseInteractor) CreateMenuItem(ctx context.Context, item *dto.MenuItemWithRecipeDTO) (*dto.CreateMenuItemResponse, error) {
	item.MenuItem.ID = uuid.New()
	for _, i := range item.RecipeItems {
		i.ID = uuid.New()
		i.MenuItemID = item.MenuItem.ID
	}

	txRepo, err := uc.menuTxRepo.WithTx(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := txRepo.Tx.Rollback(ctx); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			log.Error().Err(err).Send()
		}
	}()

	if err := txRepo.CreateMenuItem(ctx, item.MenuItem); err != nil {
		return nil, err
	}

	if err := txRepo.AddRecipeItems(ctx, item.RecipeItems); err != nil {
		return nil, err
	}

	if err := txRepo.Tx.Commit(ctx); err != nil {
		log.Error().Err(err).Send()
	}

	return &dto.CreateMenuItemResponse{Status: "Ok"}, err
}

func (uc *MenuUseCaseInteractor) UpdateMenuItem(ctx context.Context, item *dto.MenuItemWithRecipeDTO) (*dto.UpdateMenuItemResponse, error) {
	return nil, nil
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
