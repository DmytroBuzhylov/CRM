package usecase

import (
	"Test/internal/feature/permission/entity"
	"Test/internal/feature/permission/repository"
	repository2 "Test/internal/feature/staff/repository"
	"Test/pkg/store"
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type PermissionUCInteractor struct {
	permRepo     repository.PermissionRepository
	roleRepo     repository.RoleRepository
	rolePermRepo repository.RolePermissionRepository
	staffRepo    repository2.StaffRepository
	store        store.Cash
}

func NewPermissionUCInteractor(
	permRepo repository.PermissionRepository,
	roleRepo repository.RoleRepository,
	rolePermRepo repository.RolePermissionRepository,
	staffRepo repository2.StaffRepository,
	store store.Cash,
) *PermissionUCInteractor {
	return &PermissionUCInteractor{
		permRepo:     permRepo,
		roleRepo:     roleRepo,
		rolePermRepo: rolePermRepo,
		staffRepo:    staffRepo,
		store:        store,
	}
}

func (i *PermissionUCInteractor) RequestUserPermissions(ctx context.Context, userID uuid.UUID) ([]*entity.Permission, error) {

	roleID, err := i.staffRepo.GetRoleIDByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	permissions, err := i.rolePermRepo.GetPermissionsByRole(ctx, roleID)
	if err != nil {
		return nil, err
	}

	return permissions, nil
}

func (i *PermissionUCInteractor) HasPermission(ctx context.Context, userID uuid.UUID, organizationID uuid.UUID, permissionName string) bool {

	key := userID.String() + permissionName

	_, err := i.store.Get(ctx, key)
	if err == nil {
		return true
	}

	exist, err := i.permRepo.PermissionExists(ctx, userID, organizationID, permissionName)
	if err != nil {
		log.Warn().
			Str("usecase", "permissionUCInteractor").
			Str("method", "HasPermission").
			Err(err).
			Send()
		return false
	}
	if exist {
		err = i.store.Set(ctx, key, true, 1*time.Hour)
		return true
	}

	return exist
}

func (i *PermissionUCInteractor) AddPermissions(ctx context.Context, roleID uuid.UUID, permissionIDs []uuid.UUID) error {
	return i.rolePermRepo.AddPermissionToRole(ctx, roleID, permissionIDs)
}

func (i *PermissionUCInteractor) DeletePermissions(ctx context.Context, roleID uuid.UUID, permissions []uuid.UUID) error {
	return i.rolePermRepo.RemovePermissionFromRole(ctx, roleID, permissions)
}
