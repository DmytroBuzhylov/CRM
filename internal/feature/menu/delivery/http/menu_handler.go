package http

import "Test/internal/feature/menu/usecase"

type MenuHandler struct {
	uc *usecase.MenuItemUseCase
}

func NewMenuHandler(uc *usecase.MenuItemUseCase) *MenuHandler {
	return &MenuHandler{uc: uc}
}
