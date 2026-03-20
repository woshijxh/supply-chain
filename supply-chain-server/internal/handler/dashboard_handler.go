package handler

import (
	"supply-chain-server/internal/service"
	"supply-chain-server/pkg/response"

	"github.com/gin-gonic/gin"
)

type DashboardHandler struct {
	supplierService  *service.SupplierService
	inventoryService *service.InventoryService
}

func NewDashboardHandler(ss *service.SupplierService, is *service.InventoryService) *DashboardHandler {
	return &DashboardHandler{
		supplierService:  ss,
		inventoryService: is,
	}
}

func (h *DashboardHandler) Stats(c *gin.Context) {
	inventoryStats, err := h.inventoryService.GetStats()
	if err != nil {
		response.ServerError(c, "获取库存统计失败")
		return
	}

	suppliers, err := h.supplierService.GetActiveSuppliers()
	if err != nil {
		response.ServerError(c, "获取供应商列表失败")
		return
	}

	stats := map[string]interface{}{
		"todaySales":         156800.00,
		"todayOrders":        28,
		"inventoryAlert":     inventoryStats["low"],
		"pendingProcurement": 2,
		"salesGrowth":        12.5,
		"ordersGrowth":       8.3,
		"inventoryStats":     inventoryStats,
		"activeSuppliers":    len(suppliers),
	}

	response.Success(c, stats)
}
