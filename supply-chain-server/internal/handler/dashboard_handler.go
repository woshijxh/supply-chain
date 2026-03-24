package handler

import (
	"math"
	"supply-chain-server/internal/service"
	"supply-chain-server/pkg/response"

	"github.com/gin-gonic/gin"
)

type DashboardHandler struct {
	supplierService    *service.SupplierService
	inventoryService   *service.InventoryService
	salesService       *service.SalesService
	procurementService *service.ProcurementService
	productService     *service.ProductService
	customerService    *service.CustomerService
	logisticsService   *service.LogisticsService
	salesReturnService *service.SalesReturnService
}

func NewDashboardHandler(
	ss *service.SupplierService,
	is *service.InventoryService,
	sales *service.SalesService,
	ps *service.ProcurementService,
	productSvc *service.ProductService,
	customerSvc *service.CustomerService,
	logisticsSvc *service.LogisticsService,
	salesReturnSvc *service.SalesReturnService,
) *DashboardHandler {
	return &DashboardHandler{
		supplierService:    ss,
		inventoryService:   is,
		salesService:       sales,
		procurementService: ps,
		productService:     productSvc,
		customerService:    customerSvc,
		logisticsService:   logisticsSvc,
		salesReturnService: salesReturnSvc,
	}
}

func (h *DashboardHandler) Stats(c *gin.Context) {
	// 获取库存统计
	inventoryStatsRaw, err := h.inventoryService.GetStats()
	var inventoryStats map[string]interface{}
	if err != nil {
		inventoryStats = map[string]interface{}{"total": 0, "normal": 0, "low": 0, "totalSku": 0}
	} else {
		// 转换为 map[string]interface{}
		inventoryStats = make(map[string]interface{})
		for k, v := range inventoryStatsRaw {
			inventoryStats[k] = v
		}
	}

	// 获取活跃供应商数量
	suppliers, _ := h.supplierService.GetActiveSuppliers()
	activeSuppliers := len(suppliers)

	// 获取今日销售统计
	todaySales, todayOrders, _ := h.salesService.GetTodayStats()

	// 获取昨日销售统计（用于计算增长）
	yesterdaySales, yesterdayOrders, _ := h.salesService.GetYesterdayStats()

	// 计算增长率
	var salesGrowth, ordersGrowth float64
	if yesterdaySales > 0 {
		salesGrowth = math.Round((todaySales-yesterdaySales)/yesterdaySales*100*10) / 10
	} else if todaySales > 0 {
		salesGrowth = 100
	}
	if yesterdayOrders > 0 {
		ordersGrowth = math.Round(float64(todayOrders-yesterdayOrders)/float64(yesterdayOrders)*100*10) / 10
	} else if todayOrders > 0 {
		ordersGrowth = 100
	}

	// 获取待处理采购订单数量
	pendingProcurement, _ := h.procurementService.CountPending()

	// 获取客户数量
	_, customerCount, _ := h.customerService.List(1, 1000, "")

	// 获取产品数量
	products, _, _ := h.productService.List(1, 1000, "")
	productCount := products

	// 获取待确认销售订单
	pendingSales, _ := h.salesService.CountByStatus("pending")

	// 获取配送中物流
	shippingLogistics, _ := h.logisticsService.CountByStatus("shipping")

	// 获取待处理退货
	pendingReturns, _ := h.salesReturnService.CountByStatus("pending")

	stats := map[string]interface{}{
		"todaySales":         todaySales,
		"todayOrders":        todayOrders,
		"inventoryAlert":     inventoryStats["low"],
		"pendingProcurement": pendingProcurement,
		"salesGrowth":        salesGrowth,
		"ordersGrowth":       ordersGrowth,
		"inventoryStats":     inventoryStats,
		"activeSuppliers":    activeSuppliers,
		"customerCount":      customerCount,
		"productCount":       productCount,
		"pendingSales":       pendingSales,
		"shippingLogistics":  shippingLogistics,
		"pendingReturns":     pendingReturns,
	}

	response.Success(c, stats)
}
