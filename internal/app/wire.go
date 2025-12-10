package app

import (
	http4 "Test/internal/feature/inventory/delivery/http"
	postgres4 "Test/internal/feature/inventory/repository/postgres"
	usecase4 "Test/internal/feature/inventory/usecase"
	http5 "Test/internal/feature/menu/delivery/http"
	postgres5 "Test/internal/feature/menu/repository/postgres"
	usecase5 "Test/internal/feature/menu/usecase"
	http7 "Test/internal/feature/order/delivery/http"
	postgres7 "Test/internal/feature/order/repository/postgres"
	usecase7 "Test/internal/feature/order/usecase"
	http3 "Test/internal/feature/organization/delivery/http"
	usecase8 "Test/internal/feature/permission/usecase"
	"Test/internal/feature/staff/delivery/http"
	repository2 "Test/internal/feature/staff/repository"
	usecase9 "Test/internal/feature/staff/usecase"
	http6 "Test/internal/feature/table/delivery/http"
	taskHTTP "Test/internal/feature/task/delivery/http"
	"Test/internal/feature/task/repository/postgres"
	"Test/internal/feature/task/usecase"
	http2 "Test/internal/feature/user/delivery/http"
	"Test/internal/middleware"
	"Test/pkg/store"
	"time"

	"github.com/DmytroBuzhylov/ratelimiter"
)

//func (a *App) setupDependencies() {
//	limiter := ratelimiter.NewRedisTokenBucket(a.Redis, &ratelimiter.BucketOptions{
//		BucketSize: 60,
//		RefillRate: 0.5,
//		TTL:        500 * time.Second,
//	})
//	authMiddleware := middleware.NewAuthMiddleware(a.Config.JWT, limiter)
//
//	redisClient := store.NewRedisClient(a.Redis)
//
//	dependencyBuilder := initDependencyBuilder(a, redisClient)
//
//	permRepo := repository.NewPermissionRepo(a.DataProvider)
//	roleRepo := repository.NewRoleRepo(a.DataProvider)
//	rolePermRepo := repository.NewRolePermissionRepo(a.DataProvider)
//	staffRepo := repository2.NewStaffRepo(a.DataProvider)
//
//	permUC := usecase8.NewPermissionUCInteractor(permRepo, roleRepo, rolePermRepo, staffRepo, redisClient)
//	permMiddleware := middleware.NewPermissionMiddleware(permUC)
//
//	a.Router.Use(authMiddleware.RateLimiter())
//
//	a.initUserModule(redisClient)
//
//	a.Router.Use(authMiddleware.JWTMiddleware())
//
//	a.initInvitationAndOrganizationModule(permMiddleware)
//	a.initTableModule()
//	a.initTaskModule()
//	a.initOrderModule()
//	a.initIngredientModule()
//	a.initMenuModule()
//	a.InitRolePermissionModule(permUC, permMiddleware)
//}

func (a *App) setupDependencies() {
	limiter := ratelimiter.NewRedisTokenBucket(a.Redis, &ratelimiter.BucketOptions{
		BucketSize: 60,
		RefillRate: 0.5,
		TTL:        500 * time.Second,
	})

	redisClient := store.NewRedisClient(a.Redis)

	dependencyBuilder := initDependencyBuilder(a, redisClient, limiter)

	a.initUserRoutes(dependencyBuilder)

	a.Router.Use(dependencyBuilder.AuthMiddleware.RateLimiter())

	a.initInvitationAndOrganizationRoutes(dependencyBuilder)
	a.initTableRoutes(dependencyBuilder)
}

//	func (a *App) initUserModule(cash store.Cash) {
//		userRepo := postgres2.NewPostgresUserRepository(a.DB)
//		userUC := usecase2.NewUserUseCaseInteractor(userRepo, a.Config.JWT, cash)
//		authHandler := http2.NewUserHandler(userUC, a.Config.JWT)
//
//		authAPI := a.Router.Group("/auth")
//		{
//			authAPI.POST("/create", authHandler.Register)
//			authAPI.POST("/login", authHandler.Login)
//			authAPI.POST("/refresh", authHandler.Refresh)
//		}
//	}
func (a *App) initUserRoutes(dependencyBuilder *DependencyBuilder) {
	authHandler := http2.NewUserHandler(dependencyBuilder.UserUC, a.Config.JWT)

	authAPI := a.Router.Group("/auth")
	{
		authAPI.POST("/create", authHandler.Register)
		authAPI.POST("/login", authHandler.Login)
		authAPI.POST("/refresh", authHandler.Refresh)
	}
}

//	func (a *App) initInvitationAndOrganizationModule(permissionMiddleware *middleware.PermissionMiddleware) {
//		invitationRepo := postgres3.NewPostgresInvitationRepository(a.DB)
//		invitationUC := usecase3.NewInvitationUseCaseInteractor(invitationRepo)
//
//		organizationRepo := postgres3.NewPostgresOrganizationRepository(a.DB)
//		organizationUC := usecase3.NewOrganizationUseCaseInteractor(organizationRepo)
//		organizationHandler := http3.NewOrganizationHandler(organizationUC, invitationUC)
//
//		organizationAPI := a.Router.Group("/organization")
//		{
//			organizationAPI.POST("/create", organizationHandler.CreateOrganization)
//			organizationAPI.POST("/invite", permissionMiddleware.CheckPermission(""), organizationHandler.CreateInvite)
//			organizationAPI.GET("/invite/:code", organizationHandler.AcceptInvite)
//		}
//	}
func (a *App) initInvitationAndOrganizationRoutes(builder *DependencyBuilder) {
	organizationHandler := http3.NewOrganizationHandler(builder.OrganizationUC)

	organizationAPI := a.Router.Group("/organization")
	{
		organizationAPI.POST("/create", organizationHandler.CreateOrganization)
	}
}

//	func (a *App) initTableModule() {
//		tableRepo := postgres6.NewPostgresTableRepository(a.DB)
//		tableUC := usecase6.NewTableUseCaseInteractor(tableRepo)
//		tableHandler := http6.NewTableHandler(tableUC)
//
//		tableAPI := a.Router.Group("/organization/table")
//		{
//			tableAPI.POST("/create", tableHandler.CreateTableHandler)
//		}
//	}
func (a *App) initTableRoutes(builder *DependencyBuilder) {
	tableHandler := http6.NewTableHandler(builder.TableUC)

	tableAPI := a.Router.Group("/organization/table")
	{
		tableAPI.POST("/create", tableHandler.CreateTableHandler)
	}
}

func (a *App) initTaskModule() {
	taskRepo := postgres.NewTaskRepository(a.DB, a.DataProvider)
	taskUC := usecase.NewTaskUseCaseInteractor(taskRepo)
	taskHandler := taskHTTP.NewTaskHandler(taskUC)

	tasksAPI := a.Router.Group("/tasks")
	{
		tasksAPI.POST("/", taskHandler.CreateTask)
		tasksAPI.GET("/:id", taskHandler.GetTask)
		tasksAPI.GET("/", taskHandler.GetTasks)
		tasksAPI.PUT("/:id", taskHandler.UpdateTask)
		tasksAPI.PATCH("/:id/status", taskHandler.UpdateTaskStatus)
		tasksAPI.PATCH("/:id/deadline", taskHandler.UpdateTaskDeadline)
		tasksAPI.DELETE("/:id", taskHandler.DeleteTask)
	}

}

func (a *App) initOrderModule() {
	orderRepo := postgres7.NewPostgresOrderRepository(a.DB)
	orderUC := usecase7.NewOrderUseCaseInteractor(orderRepo)
	orderHandler := http7.NewOrderHandler(orderUC)

	orderAPI := a.Router.Group("/organization/order")
	{
		orderAPI.GET("/:id", orderHandler.GetOrderHandler)
		orderAPI.GET("/", orderHandler.GetOrganizationsOrdersHandler)
		orderAPI.POST("/create", orderHandler.CreateOrderHandler)
	}

}

func (a *App) initIngredientModule() {
	inventoryRepo := postgres4.NewPostgresInventoryRepository(a.DB)
	inventoryUC := usecase4.NewInventoryUseCaseInteractor(inventoryRepo, a.StorageS3)
	inventoryHandler := http4.NewInventoryHandler(inventoryUC)

	ingredientAPI := a.Router.Group("/organization/ingredient")
	{
		ingredientAPI.GET("/", inventoryHandler.GetOrganizationIngredients)
		ingredientAPI.GET("/:id", inventoryHandler.GetIngredientHandler)
		ingredientAPI.POST("/create", inventoryHandler.CreateIngredientHandler)
		ingredientAPI.PUT("/update", inventoryHandler.UpdateIngredientHandler)
		ingredientAPI.DELETE("/:id", inventoryHandler.DeleteIngredientHandler)
		ingredientAPI.DELETE("/delete", inventoryHandler.DeleteManyIngredientsHandler)
	}
}

func (a *App) initMenuModule() {
	menuRepo := postgres5.NewPostgresMenuRepository(a.DB)
	menuUC := usecase5.NewMenuUseCaseInteractor(menuRepo, a.StorageS3)
	menuHandler := http5.NewMenuHandler(menuUC)

	menuAPI := a.Router.Group("/organization/menu")
	{
		menuAPI.GET("/:id")
		menuAPI.PUT("/update", menuHandler.UpdateMenuItemHandler)
		menuAPI.POST("/create", menuHandler.CreateMenuItemHandler)
	}
}

func (a *App) initStaffModule(permissionMiddleware *middleware.PermissionMiddleware) {

	staffRepo := repository2.NewStaffRepo(a.DataProvider)
	inviteRepo := repository2.NewInvitationStaffRepo(a.DataProvider)
	breakRepo := repository2.NewBreakRepo(a.DataProvider)
	workDayRepo := repository2.NewWorkDayRepo(a.DataProvider)

	staffUC := usecase9.NewStaffUCInteractor(staffRepo, workDayRepo, breakRepo, inviteRepo)

	staffHandler := http.NewStaffHandler(staffUC)

	staffAPI := a.Router.Group("/organization/staff")
	{
		staffAPI.GET("/invite/:code")
		staffAPI.POST("/invite", permissionMiddleware.CheckPermission("staff:manage"), staffHandler.CreateInviteHandler)
		staffAPI.DELETE("/:id", permissionMiddleware.CheckPermission("staff:manage"), staffHandler.DismissStaffHandler)
	}

}

func (a *App) initNotificationModule() {

}

func (a *App) InitRolePermissionModule(permUseCase usecase8.PermissionUseCase, permissionMiddleware *middleware.PermissionMiddleware) {

}
