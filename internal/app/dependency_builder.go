package app

import (
	inventory_repo "Test/internal/feature/inventory/repository/postgres"
	inventory_usecase "Test/internal/feature/inventory/usecase"
	menu_repo "Test/internal/feature/menu/repository/postgres"
	menu_usecase "Test/internal/feature/menu/usecase"
	order_repo "Test/internal/feature/order/repository/postgres"
	order_usecase "Test/internal/feature/order/usecase"
	organization_repo "Test/internal/feature/organization/repository/postgres"
	organization_usecase "Test/internal/feature/organization/usecase"
	perm_repo "Test/internal/feature/permission/repository"
	perm_usecase "Test/internal/feature/permission/usecase"
	staff_repo "Test/internal/feature/staff/repository"
	staff_usecase "Test/internal/feature/staff/usecase"
	table_repo "Test/internal/feature/table/repository/postgres"
	table_usecase "Test/internal/feature/table/usecase"
	task_repo "Test/internal/feature/task/repository/postgres"
	task_usecase "Test/internal/feature/task/usecase"
	user_repo "Test/internal/feature/user/repository/postgres"
	user_usecase "Test/internal/feature/user/usecase"
	"Test/internal/middleware"
	"Test/pkg/store"

	"github.com/DmytroBuzhylov/ratelimiter"
)

type DependencyBuilder struct {
	RolePermissionRepo *perm_repo.RolePermissionRepo
	RoleRepo           *perm_repo.RoleRepo
	PermRepo           *perm_repo.PermissionRepo
	StaffRepo          *staff_repo.StaffRepo
	TaskRepo           *task_repo.TaskRepo
	UserRepo           *user_repo.UserRepo
	OrganizationRepo   *organization_repo.PostgresOrganizationRepository
	TableRepo          *table_repo.PostgresTableRepository
	OrderRepo          *order_repo.PostgresOrderRepository
	InventoryRepo      *inventory_repo.PostgresInventoryRepository
	MenuRepo           *menu_repo.PostgresMenuRepository
	InviteRepo         *staff_repo.InvitationStaffRepo
	BreakRepo          *staff_repo.BreakRepo
	WorkDayRepo        *staff_repo.WorkDayRepo

	//invitation feature is not used, it is currently InivitationStaff used
	//InvitationUC   *organization_usecase.InvitationUseCaseInteractor
	PermUC         *perm_usecase.PermissionUCInteractor
	UserUC         *user_usecase.UserUseCaseInteractor
	OrganizationUC *organization_usecase.OrganizationUseCaseInteractor
	TableUC        *table_usecase.TableUseCaseInteractor
	TaskUC         *task_usecase.TaskUseCaseInteractor
	OrderUC        *order_usecase.OrderUseCaseInteractor
	InventoryUC    *inventory_usecase.InventoryUseCaseInteractor
	MenuUC         *menu_usecase.MenuUseCaseInteractor
	StaffUC        *staff_usecase.StaffUCInteractor

	PermMiddleware *middleware.PermissionMiddleware
	AuthMiddleware *middleware.AuthMiddleware
}

func initDependencyBuilder(a *App, redisClient *store.RedisClient, limiter *ratelimiter.RedisTokenBucket) *DependencyBuilder {
	db := &DependencyBuilder{}
	db.PermRepo = perm_repo.NewPermissionRepo(a.DataProvider)
	db.RoleRepo = perm_repo.NewRoleRepo(a.DataProvider)
	db.RolePermissionRepo = perm_repo.NewRolePermissionRepo(a.DataProvider)
	db.StaffRepo = staff_repo.NewStaffRepo(a.DataProvider)
	db.UserRepo = user_repo.NewPostgresUserRepository(a.DB)
	db.OrganizationRepo = organization_repo.NewPostgresOrganizationRepository(a.DB)
	db.TableRepo = table_repo.NewPostgresTableRepository(a.DB)
	db.TaskRepo = task_repo.NewTaskRepository(a.DB, a.DataProvider)
	db.OrderRepo = order_repo.NewPostgresOrderRepository(a.DB)
	db.InventoryRepo = inventory_repo.NewPostgresInventoryRepository(a.DB)
	db.MenuRepo = menu_repo.NewPostgresMenuRepository(a.DB)
	db.StaffRepo = staff_repo.NewStaffRepo(a.DataProvider)
	db.InviteRepo = staff_repo.NewInvitationStaffRepo(a.DataProvider)
	db.BreakRepo = staff_repo.NewBreakRepo(a.DataProvider)
	db.WorkDayRepo = staff_repo.NewWorkDayRepo(a.DataProvider)

	db.StaffUC = staff_usecase.NewStaffUCInteractor(db.StaffRepo, db.WorkDayRepo, db.BreakRepo, db.InviteRepo)
	db.MenuUC = menu_usecase.NewMenuUseCaseInteractor(db.MenuRepo, a.StorageS3)
	db.InventoryUC = inventory_usecase.NewInventoryUseCaseInteractor(db.InventoryRepo, a.StorageS3)
	db.OrderUC = order_usecase.NewOrderUseCaseInteractor(db.OrderRepo)
	db.TaskUC = task_usecase.NewTaskUseCaseInteractor(db.TaskRepo)
	db.TableUC = table_usecase.NewTableUseCaseInteractor(db.TableRepo)
	db.OrganizationUC = organization_usecase.NewOrganizationUseCaseInteractor(db.OrganizationRepo)
	db.UserUC = user_usecase.NewUserUseCaseInteractor(db.UserRepo, a.Config.JWT, redisClient)
	db.PermUC = perm_usecase.NewPermissionUCInteractor(db.PermRepo, db.RoleRepo, db.RolePermissionRepo, db.StaffRepo, redisClient)

	db.PermMiddleware = middleware.NewPermissionMiddleware(db.PermUC)
	db.AuthMiddleware = middleware.NewAuthMiddleware(a.Config.JWT, limiter)

	return db
}
