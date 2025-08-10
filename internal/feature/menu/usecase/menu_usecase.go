package usecase

import "context"

type MenuItemUseCase interface {
	GetMenuItem(ctx context.Context)
	GetOrganizationMenuItems(ctx context.Context)
	CreateMenuItem(ctx context.Context)
	UpdateMenuItem(ctx context.Context)
	DeleteMenuItem(ctx context.Context)
}
