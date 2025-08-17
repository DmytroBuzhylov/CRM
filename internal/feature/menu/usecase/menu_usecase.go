package usecase

import (
	"Test/internal/feature/menu/interface_adapter/dto"
	"context"
)

type MenuItemUseCase interface {
	CreateMenuItem(ctx context.Context, item *dto.MenuItemWithRecipeDTO) (*dto.CreateMenuItemResponse, error)
	UpdateMenuItem(ctx context.Context, item *dto.MenuItemWithRecipeDTO) (*dto.UpdateMenuItemResponse, error)
	DeleteMenuItem(ctx context.Context, request *dto.DeleteMenuItemRequest) (*dto.DeleteMenuItemResponse, error)
	GetMenuItem(ctx context.Context, request *dto.GetMenuItemRequest) (*dto.MenuItemWithRecipeDTO, error)
	GetAllMenuItems(ctx context.Context, request *dto.GetAllMenuItemsRequest) ([]*dto.MenuItemWithRecipeDTO, error)
}
