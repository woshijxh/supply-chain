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
	customerRepo := repository.NewCustomerRepository(database.DB)
	inventoryLogRepo := repository.NewInventoryLogRepository(database.DB)
	salesReturnRepo := repository.NewSalesReturnRepository(database.DB)
	procurementReturnRepo := repository.NewProcurementReturnRepository(database.DB)
	traceRepo := repository.NewTraceRepository(database.DB)

	// 初始化 Service
	userService := service.NewUserService(userRepo)
	roleService := service.NewRoleService(roleRepo)
	permService := service.NewPermissionService(permRepo)
	userProfileService := service.NewUserProfileService(userProfileRepo, userRepo)
	supplierService := service.NewSupplierService(supplierRepo)
	inventoryService := service.NewInventoryService(inventoryRepo, procurementRepo, inventoryLogRepo)
	inventoryService.SetProductRepo(productRepo) // 设置产品仓库用于库存流水记录
	procurementService := service.NewProcurementService(procurementRepo, inventoryService)
	logisticsService := service.NewLogisticsService(logisticsRepo)
	salesService := service.NewSalesService(salesRepo, inventoryService, logisticsService)
	productService := service.NewProductService(productRepo)
	customerService := service.NewCustomerService(customerRepo)
	salesReturnService := service.NewSalesReturnService(salesReturnRepo, salesRepo, inventoryService)
	procurementReturnService := service.NewProcurementReturnService(procurementReturnRepo, procurementRepo, inventoryService)
	traceService := service.NewTraceService(traceRepo, productRepo)

	// 初始化 Handler
	captchaHandler := handler.NewCaptchaHandler()
	authHandler := handler.NewAuthHandler(userService)
	userHandler := handler.NewUserHandler(userService, userProfileService)
	roleHandler := handler.NewRoleHandler(roleService, permService)
	permHandler := handler.NewPermissionHandler(permService)
	supplierHandler := handler.NewSupplierHandler(supplierService)
	inventoryHandler := handler.NewInventoryHandler(inventoryService)
	procurementHandler := handler.NewProcurementHandler(procurementService, productRepo)
	salesHandler := handler.NewSalesHandler(salesService)
	logisticsHandler := handler.NewLogisticsHandler(logisticsService)
	productHandler := handler.NewProductHandler(productService)
	customerHandler := handler.NewCustomerHandler(customerService)
	salesReturnHandler := handler.NewSalesReturnHandler(salesReturnService)
	procurementReturnHandler := handler.NewProcurementReturnHandler(procurementReturnService)
	traceHandler := handler.NewTraceHandler(traceService)
	dashboardHandler := handler.NewDashboardHandler(supplierService, inventoryService, salesService, procurementService, productService, customerService, logisticsService, salesReturnService)

	api := r.Group("/api")
	{
		api.POST("/auth/login", authHandler.Login)
		api.POST("/auth/register", authHandler.Register)
		api.GET("/auth/captcha", captchaHandler.Get)
	}

	protected := api.Group("")
	protected.Use(middleware.JWTAuth())
	{
		// 用户相关
		protected.GET("/auth/profile", authHandler.Profile)
		protected.PUT("/auth/password", authHandler.ChangePassword)
		protected.PUT("/auth/avatar", authHandler.UpdateAvatar)

		// 用户管理
		protected.GET("/users", middleware.RBACAuth("user", "read"), userHandler.List)
		protected.GET("/users/:id", middleware.RBACAuth("user", "read"), userHandler.Get)
		protected.POST("/users", middleware.RBACAuth("user", "create"), userHandler.Create)
		protected.PUT("/users/:id", middleware.RBACAuth("user", "update"), userHandler.Update)
		protected.DELETE("/users/:id", middleware.RBACAuth("user", "delete"), userHandler.Delete)
		protected.PUT("/users/:id/password", middleware.RBACAuth("user", "update"), userHandler.ResetPassword)
		protected.GET("/users/:id/roles", middleware.RBACAuth("user", "read"), userHandler.GetRoles)
		protected.PUT("/users/:id/roles", middleware.RBACAuth("user", "update"), userHandler.SetRoles)
		protected.GET("/users/:id/permissions", middleware.RBACAuth("user", "read"), userHandler.GetPermissions)
		protected.PUT("/users/:id/permissions", middleware.RBACAuth("user", "update"), userHandler.SetPermissions)

		// 角色管理
		protected.GET("/roles", middleware.RBACAuth("role", "read"), roleHandler.List)
		protected.GET("/roles/all", middleware.RBACAuth("role", "read"), roleHandler.All)
		protected.GET("/roles/:id", middleware.RBACAuth("role", "read"), roleHandler.Get)
		protected.POST("/roles", middleware.RBACAuth("role", "create"), roleHandler.Create)
		protected.PUT("/roles/:id", middleware.RBACAuth("role", "update"), roleHandler.Update)
		protected.DELETE("/roles/:id", middleware.RBACAuth("role", "delete"), roleHandler.Delete)
		protected.GET("/roles/:id/permissions", middleware.RBACAuth("role", "read"), roleHandler.GetPermissions)
		protected.POST("/roles/:id/permissions", middleware.RBACAuth("role", "update"), roleHandler.SetPermissions)

		// 权限管理
		protected.GET("/permissions", middleware.RBACAuth("permission", "read"), permHandler.List)
		protected.GET("/permissions/all", middleware.RBACAuth("permission", "read"), permHandler.GetAll)
		protected.GET("/permissions/:id", middleware.RBACAuth("permission", "read"), permHandler.Get)
		protected.POST("/permissions", middleware.RBACAuth("permission", "create"), permHandler.Create)
		protected.PUT("/permissions/:id", middleware.RBACAuth("permission", "update"), permHandler.Update)
		protected.DELETE("/permissions/:id", middleware.RBACAuth("permission", "delete"), permHandler.Delete)

		// 供应商
		protected.GET("/suppliers", middleware.RBACAuth("supplier", "read"), supplierHandler.List)
		protected.GET("/suppliers/:id", middleware.RBACAuth("supplier", "read"), supplierHandler.Get)
		protected.POST("/suppliers", middleware.RBACAuth("supplier", "create"), supplierHandler.Create)
		protected.PUT("/suppliers/:id", middleware.RBACAuth("supplier", "update"), supplierHandler.Update)
		protected.DELETE("/suppliers/:id", middleware.RBACAuth("supplier", "delete"), supplierHandler.Delete)

		// 库存
		protected.GET("/inventory", middleware.RBACAuth("inventory", "read"), inventoryHandler.List)
		protected.GET("/inventory/:id", middleware.RBACAuth("inventory", "read"), inventoryHandler.Get)
		protected.POST("/inventory/stock-in", middleware.RBACAuth("inventory", "create"), inventoryHandler.StockIn)
		protected.POST("/inventory/stock-out", middleware.RBACAuth("inventory", "update"), inventoryHandler.StockOut)
		protected.GET("/inventory/stats", middleware.RBACAuth("inventory", "read"), inventoryHandler.Stats)
		protected.GET("/inventory/logs", middleware.RBACAuth("inventory", "read"), inventoryHandler.Logs)

		// 客户
		protected.GET("/customers", middleware.RBACAuth("customer", "read"), customerHandler.List)
		protected.GET("/customers/all", middleware.RBACAuth("customer", "read"), customerHandler.All)
		protected.GET("/customers/:id", middleware.RBACAuth("customer", "read"), customerHandler.Get)
		protected.POST("/customers", middleware.RBACAuth("customer", "create"), customerHandler.Create)
		protected.PUT("/customers/:id", middleware.RBACAuth("customer", "update"), customerHandler.Update)
		protected.DELETE("/customers/:id", middleware.RBACAuth("customer", "delete"), customerHandler.Delete)

		// 产品
		protected.GET("/products", middleware.RBACAuth("product", "read"), productHandler.List)
		protected.GET("/products/:id", middleware.RBACAuth("product", "read"), productHandler.Get)
		protected.POST("/products", middleware.RBACAuth("product", "create"), productHandler.Create)
		protected.PUT("/products/:id", middleware.RBACAuth("product", "update"), productHandler.Update)
		protected.DELETE("/products/:id", middleware.RBACAuth("product", "delete"), productHandler.Delete)

		// 采购
		protected.GET("/procurement", middleware.RBACAuth("procurement", "read"), procurementHandler.List)
		protected.GET("/procurement/:id", middleware.RBACAuth("procurement", "read"), procurementHandler.Get)
		protected.POST("/procurement", middleware.RBACAuth("procurement", "create"), procurementHandler.Create)
		protected.PUT("/procurement/:id/status", middleware.RBACAuth("procurement", "update"), procurementHandler.UpdateStatus)
		protected.DELETE("/procurement/:id", middleware.RBACAuth("procurement", "delete"), procurementHandler.Delete)

		// 销售
		protected.GET("/sales", middleware.RBACAuth("sales", "read"), salesHandler.List)
		protected.GET("/sales/:id", middleware.RBACAuth("sales", "read"), salesHandler.Get)
		protected.POST("/sales", middleware.RBACAuth("sales", "create"), salesHandler.Create)
		protected.PUT("/sales/:id/status", middleware.RBACAuth("sales", "update"), salesHandler.UpdateStatus)
		protected.DELETE("/sales/:id", middleware.RBACAuth("sales", "delete"), salesHandler.Delete)

		// 物流
		protected.GET("/logistics", middleware.RBACAuth("logistics", "read"), logisticsHandler.List)
		protected.GET("/logistics/:id", middleware.RBACAuth("logistics", "read"), logisticsHandler.Get)
		protected.POST("/logistics", middleware.RBACAuth("logistics", "create"), logisticsHandler.Create)
		protected.PUT("/logistics/:id/status", middleware.RBACAuth("logistics", "update"), logisticsHandler.UpdateStatus)
		protected.DELETE("/logistics/:id", middleware.RBACAuth("logistics", "delete"), logisticsHandler.Delete)

		// 仪表盘
		protected.GET("/dashboard/stats", middleware.RBACAuth("dashboard", "read"), dashboardHandler.Stats)

		// 商品追溯
		protected.GET("/trace", middleware.RBACAuth("trace", "read"), traceHandler.Trace)

		// 销售退货
		protected.GET("/sales-returns", middleware.RBACAuth("sales", "read"), salesReturnHandler.List)
		protected.GET("/sales-returns/:id", middleware.RBACAuth("sales", "read"), salesReturnHandler.Get)
		protected.POST("/sales-returns", middleware.RBACAuth("sales", "create"), salesReturnHandler.Create)
		protected.PUT("/sales-returns/:id/approve", middleware.RBACAuth("sales", "update"), salesReturnHandler.Approve)
		protected.PUT("/sales-returns/:id/reject", middleware.RBACAuth("sales", "update"), salesReturnHandler.Reject)
		protected.PUT("/sales-returns/:id/complete", middleware.RBACAuth("sales", "update"), salesReturnHandler.Complete)
		protected.DELETE("/sales-returns/:id", middleware.RBACAuth("sales", "delete"), salesReturnHandler.Delete)

		// 采购退货
		protected.GET("/procurement-returns", middleware.RBACAuth("procurement", "read"), procurementReturnHandler.List)
		protected.GET("/procurement-returns/:id", middleware.RBACAuth("procurement", "read"), procurementReturnHandler.Get)
		protected.POST("/procurement-returns", middleware.RBACAuth("procurement", "create"), procurementReturnHandler.Create)
		protected.PUT("/procurement-returns/:id/approve", middleware.RBACAuth("procurement", "update"), procurementReturnHandler.Approve)
		protected.PUT("/procurement-returns/:id/reject", middleware.RBACAuth("procurement", "update"), procurementReturnHandler.Reject)
		protected.PUT("/procurement-returns/:id/complete", middleware.RBACAuth("procurement", "update"), procurementReturnHandler.Complete)
		protected.DELETE("/procurement-returns/:id", middleware.RBACAuth("procurement", "delete"), procurementReturnHandler.Delete)
	}

	return r
}