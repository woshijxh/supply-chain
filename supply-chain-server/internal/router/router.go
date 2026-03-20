package router

import (
	"supply-chain-server/internal/handler"
	"supply-chain-server/internal/middleware"
	"supply-chain-server/internal/repository"
	"supply-chain-server/internal/service"
	"supply-chain-server/pkg/casbin"
	"supply-chain-server/pkg/database"

	"github.com/gin-gonic/gin"
)

func Setup() *gin.Engine {
	r := gin.Default()
	r.Use(middleware.CORS())

	// 初始化 Casbin
	if err := casbin.InitCasbin(database.DB); err != nil {
		panic("Casbin 初始化失败: " + err.Error())
	}

	// 初始化 Repository
	userRepo := repository.NewUserRepository(database.DB)
	roleRepo := repository.NewRoleRepository(database.DB)
	permRepo := repository.NewPermissionRepository(database.DB)
	userProfileRepo := repository.NewUserProfileRepository(database.DB)
	supplierRepo := repository.NewSupplierRepository(database.DB)
	inventoryRepo := repository.NewInventoryRepository(database.DB)
	procurementRepo := repository.NewProcurementRepository(database.DB)
	salesRepo := repository.NewSalesRepository(database.DB)
	logisticsRepo := repository.NewLogisticsRepository(database.DB)
	productRepo := repository.NewProductRepository(database.DB)

	// 初始化 Service
	userService := service.NewUserService(userRepo)
	roleService := service.NewRoleService(roleRepo)
	permService := service.NewPermissionService(permRepo)
	userProfileService := service.NewUserProfileService(userProfileRepo, userRepo)
	supplierService := service.NewSupplierService(supplierRepo)
	inventoryService := service.NewInventoryService(inventoryRepo, procurementRepo)
	procurementService := service.NewProcurementService(procurementRepo, inventoryService)
	salesService := service.NewSalesService(salesRepo)
	logisticsService := service.NewLogisticsService(logisticsRepo)
	productService := service.NewProductService(productRepo)

	// 初始化 Handler
	authHandler := handler.NewAuthHandler(userService)
	userHandler := handler.NewUserHandler(userService, userProfileService)
	roleHandler := handler.NewRoleHandler(roleService, permService)
	permHandler := handler.NewPermissionHandler(permService)
	supplierHandler := handler.NewSupplierHandler(supplierService)
	inventoryHandler := handler.NewInventoryHandler(inventoryService)
	dashboardHandler := handler.NewDashboardHandler(supplierService, inventoryService)
	procurementHandler := handler.NewProcurementHandler(procurementService, productRepo)
	salesHandler := handler.NewSalesHandler(salesService)
	logisticsHandler := handler.NewLogisticsHandler(logisticsService)
	productHandler := handler.NewProductHandler(productService)

	api := r.Group("/api")
	{
		api.POST("/auth/login", authHandler.Login)
		api.POST("/auth/register", authHandler.Register)
	}

	protected := api.Group("")
	protected.Use(middleware.JWTAuth())
	{
		// 用户相关 - 需要 user:read 或 user:write
		protected.GET("/auth/profile", authHandler.Profile)
		protected.PUT("/auth/password", authHandler.ChangePassword)
		protected.PUT("/auth/avatar", authHandler.UpdateAvatar)
		protected.GET("/users", middleware.RBACAuth("user", "read"), userHandler.List)
		protected.GET("/users/:id", middleware.RBACAuth("user", "read"), userHandler.Get)
		protected.POST("/users", middleware.RBACAuth("user", "write"), userHandler.Create)
		protected.PUT("/users/:id", middleware.RBACAuth("user", "write"), userHandler.Update)
		protected.DELETE("/users/:id", middleware.RBACAuth("user", "write"), userHandler.Delete)
		protected.PUT("/users/:id/password", middleware.RBACAuth("user", "write"), userHandler.ResetPassword)
		protected.GET("/users/:id/roles", middleware.RBACAuth("user", "read"), userHandler.GetRoles)
		protected.PUT("/users/:id/roles", middleware.RBACAuth("user", "write"), userHandler.SetRoles)
		protected.GET("/users/:id/permissions", middleware.RBACAuth("user", "read"), userHandler.GetPermissions)
		protected.PUT("/users/:id/permissions", middleware.RBACAuth("user", "write"), userHandler.SetPermissions)

		// 角色管理 - 仅 admin 可访问
		protected.GET("/roles", middleware.RBACAuth("role", "read"), roleHandler.List)
		protected.GET("/roles/all", middleware.RBACAuth("role", "read"), roleHandler.All)
		protected.GET("/roles/:id", middleware.RBACAuth("role", "read"), roleHandler.Get)
		protected.POST("/roles", middleware.RBACAuth("role", "write"), roleHandler.Create)
		protected.PUT("/roles/:id", middleware.RBACAuth("role", "write"), roleHandler.Update)
		protected.DELETE("/roles/:id", middleware.RBACAuth("role", "write"), roleHandler.Delete)
		protected.GET("/roles/:id/permissions", middleware.RBACAuth("role", "read"), roleHandler.GetPermissions)
		protected.POST("/roles/:id/permissions", middleware.RBACAuth("role", "write"), roleHandler.SetPermissions)

		// 权限管理 - 仅 admin 可访问
		protected.GET("/permissions", middleware.RBACAuth("permission", "read"), permHandler.List)
		protected.GET("/permissions/all", middleware.RBACAuth("permission", "read"), permHandler.GetAll)
		protected.GET("/permissions/:id", middleware.RBACAuth("permission", "read"), permHandler.Get)
		protected.POST("/permissions", middleware.RBACAuth("permission", "write"), permHandler.Create)
		protected.PUT("/permissions/:id", middleware.RBACAuth("permission", "write"), permHandler.Update)
		protected.DELETE("/permissions/:id", middleware.RBACAuth("permission", "write"), permHandler.Delete)

		// 供应商 - 需要 supplier:read 或 supplier:write
		protected.GET("/suppliers", middleware.RBACAuth("supplier", "read"), supplierHandler.List)
		protected.GET("/suppliers/:id", middleware.RBACAuth("supplier", "read"), supplierHandler.Get)
		protected.POST("/suppliers", middleware.RBACAuth("supplier", "write"), supplierHandler.Create)
		protected.PUT("/suppliers/:id", middleware.RBACAuth("supplier", "write"), supplierHandler.Update)
		protected.DELETE("/suppliers/:id", middleware.RBACAuth("supplier", "write"), supplierHandler.Delete)

		// 库存 - 需要 inventory:read 或 inventory:write
		protected.GET("/inventory", middleware.RBACAuth("inventory", "read"), inventoryHandler.List)
		protected.GET("/inventory/:id", middleware.RBACAuth("inventory", "read"), inventoryHandler.Get)
		protected.POST("/inventory/stock-in", middleware.RBACAuth("inventory", "write"), inventoryHandler.StockIn)
		protected.POST("/inventory/stock-out", middleware.RBACAuth("inventory", "write"), inventoryHandler.StockOut)
		protected.GET("/inventory/stats", middleware.RBACAuth("inventory", "read"), inventoryHandler.Stats)

		// 产品 - 需要 product:read 或 product:write
		protected.GET("/products", middleware.RBACAuth("product", "read"), productHandler.List)
		protected.GET("/products/:id", middleware.RBACAuth("product", "read"), productHandler.Get)
		protected.POST("/products", middleware.RBACAuth("product", "write"), productHandler.Create)

		// 采购 - 需要 procurement:read 或 procurement:write
		protected.GET("/procurement", middleware.RBACAuth("procurement", "read"), procurementHandler.List)
		protected.GET("/procurement/:id", middleware.RBACAuth("procurement", "read"), procurementHandler.Get)
		protected.POST("/procurement", middleware.RBACAuth("procurement", "write"), procurementHandler.Create)
		protected.PUT("/procurement/:id/status", middleware.RBACAuth("procurement", "write"), procurementHandler.UpdateStatus)
		protected.DELETE("/procurement/:id", middleware.RBACAuth("procurement", "write"), procurementHandler.Delete)

		// 销售 - 需要 sales:read 或 sales:write
		protected.GET("/sales", middleware.RBACAuth("sales", "read"), salesHandler.List)
		protected.GET("/sales/:id", middleware.RBACAuth("sales", "read"), salesHandler.Get)
		protected.POST("/sales", middleware.RBACAuth("sales", "write"), salesHandler.Create)
		protected.PUT("/sales/:id/status", middleware.RBACAuth("sales", "write"), salesHandler.UpdateStatus)
		protected.DELETE("/sales/:id", middleware.RBACAuth("sales", "write"), salesHandler.Delete)

		// 物流 - 需要 logistics:read 或 logistics:write
		protected.GET("/logistics", middleware.RBACAuth("logistics", "read"), logisticsHandler.List)
		protected.GET("/logistics/:id", middleware.RBACAuth("logistics", "read"), logisticsHandler.Get)
		protected.POST("/logistics", middleware.RBACAuth("logistics", "write"), logisticsHandler.Create)
		protected.PUT("/logistics/:id/status", middleware.RBACAuth("logistics", "write"), logisticsHandler.UpdateStatus)
		protected.DELETE("/logistics/:id", middleware.RBACAuth("logistics", "write"), logisticsHandler.Delete)

		// 仪表盘 - 需要 dashboard:read
		protected.GET("/dashboard/stats", middleware.RBACAuth("dashboard", "read"), dashboardHandler.Stats)
	}

	return r
}