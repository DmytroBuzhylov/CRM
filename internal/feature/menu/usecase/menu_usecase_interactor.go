package usecase

import (
	"Test/internal/feature/menu/entity"
	"Test/internal/feature/menu/interface_adapter/dto"
	"Test/internal/feature/menu/repository"
	"Test/pkg/storage"
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
	"github.com/shopspring/decimal"
	"time"
)

type MenuUseCaseInteractor struct {
	menuRepo repository.MenuRepository
	storage  storage.Storage
}

func NewMenuUseCaseInteractor(menuRepo repository.MenuRepository, storage storage.Storage) *MenuUseCaseInteractor {
	return &MenuUseCaseInteractor{menuRepo: menuRepo, storage: storage}
}

func (uc *MenuUseCaseInteractor) CreateMenuItem(ctx context.Context, itemDTO *dto.MenuItemWithRecipeDTO, fileOpts *storage.StorageFileOptions) error {

	price, err := decimal.NewFromString(itemDTO.MenuItem.Price)
	if err != nil {
		return err
	}

	menuItem := entity.MenuItem{
		ID:             uuid.New(),
		OrganizationID: itemDTO.MenuItem.OrganizationID,
		Name:           itemDTO.MenuItem.Name,
		Description:    itemDTO.MenuItem.Description,
		Price: decimal.NullDecimal{
			Decimal: price,
		},
		Category:    itemDTO.MenuItem.Category,
		IsAvailable: true,
		UpdatedAt:   time.Now(),
		CreatedAt:   time.Now(),
	}

	var recipeItems = make([]entity.RecipeItem, 0, len(itemDTO.RecipeItems))

	for _, i := range itemDTO.RecipeItems {
		recipeItems = append(recipeItems, entity.RecipeItem{
			ID:             uuid.New(),
			MenuItemID:     menuItem.ID,
			IngredientID:   i.IngredientID,
			QuantityNeeded: i.QuantityNeeded,
			UnitOfMeasure:  i.UnitOfMeasure,
			UpdatedAt:      time.Now(),
			CreatedAt:      time.Now(),
		})
	}

	txRepo, err := uc.menuRepo.StartTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() {
		if err := txRepo.Rollback(ctx); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			log.Error().Err(err).Send()
		}
	}()

	if err := txRepo.CreateMenuItem(ctx, menuItem); err != nil {
		return err
	}

	if err := txRepo.AddRecipeItems(ctx, recipeItems); err != nil {
		return err
	}

	if err := uc.storage.UploadFile(ctx, menuItem.ID.String(), fileOpts.File, fileOpts.FileSize, fileOpts.ContentType); err != nil {
		return err
	}

	if err := txRepo.Commit(ctx); err != nil {
		log.Error().Err(err).Send()
	}

	return nil
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

	//err = tx.UpdateRecipeItem(ctx, item.RecipeItems[0])
	//if err != nil {
	//	return nil, err
	//}
	//
	//err = tx.UpdateMenuItem(ctx, item.MenuItem)
	//if err != nil {
	//	return nil, err
	//}

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
