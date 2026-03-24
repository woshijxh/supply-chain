package handler

import (
	"strconv"
	"supply-chain-server/internal/service"
	"supply-chain-server/pkg/response"

	"github.com/gin-gonic/gin"
)

type InventoryHandler struct {
	service *service.InventoryService
}

func NewInventoryHandler(s *service.InventoryService) *InventoryHandler {
	return &InventoryHandler{service: s}
}

func (h *InventoryHandler) List(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}
	pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	if err != nil || pageSize < 1 {
		pageSize = 10
	}
	status := c.Query("status")
	warehouse := c.Query("warehouse")

	items, total, err := h.service.List(page, pageSize, status, warehouse)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.PageSuccess(c, items, total, page, pageSize)
}

func (h *InventoryHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}
	item, err := h.service.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "库存不存在")
		return
	}
	response.Success(c, item)
}

type StockOperationRequest struct {
	ProductID         uint   `json:"productId" binding:"required"`
	Quantity          int    `json:"quantity" binding:"required"`
	Warehouse         string `json:"warehouse"`
	Remark            string `json:"remark"`
	ProcurementID     uint   `json:"procurementId,omitempty"`
	ProcurementItemID uint   `json:"procurementItemId,omitempty"`
}

func (h *InventoryHandler) StockIn(c *gin.Context) {
	var req StockOperationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	if err := h.service.StockInWithProcurement(req.ProductID, uint(req.Quantity), req.Warehouse, req.ProcurementID, req.ProcurementItemID); err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, "入库成功", nil)
}

func (h *InventoryHandler) StockOut(c *gin.Context) {
	var req StockOperationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	if err := h.service.StockOut(req.ProductID, uint(req.Quantity), req.Warehouse); err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, "出库成功", nil)
}

func (h *InventoryHandler) Stats(c *gin.Context) {
	stats, err := h.service.GetStats()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, stats)
}

func (h *InventoryHandler) Logs(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}
	pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	if err != nil || pageSize < 1 {
		pageSize = 20
	}

	var productID uint
	if pid := c.Query("productId"); pid != "" {
		id, err := strconv.ParseUint(pid, 10, 32)
		if err == nil {
			productID = uint(id)
		}
	}
	logType := c.Query("type")

	logs, total, err := h.service.GetLogs(page, pageSize, productID, logType)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.PageSuccess(c, logs, total, page, pageSize)
}
