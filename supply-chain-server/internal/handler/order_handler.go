package handler

import (
	"strconv"
	"supply-chain-server/internal/model"
	"supply-chain-server/internal/repository"
	"supply-chain-server/internal/service"
	"supply-chain-server/pkg/response"
	"time"

	"github.com/gin-gonic/gin"
)

// ProcurementOrderRequest 用于接收前端传来的采购订单数据
type ProcurementOrderRequest struct {
	SupplierID  uint                   `json:"supplierId" binding:"required"`
	SupplierName string                `json:"supplierName"`
	Warehouse   string                 `json:"warehouse"`
	ExpectedDate string                `json:"expectedDate"`
	Remark      string                 `json:"remark"`
	Items       []ProcurementItemRequest `json:"items"`
	TotalAmount float64                `json:"totalAmount"`
	Status      string                 `json:"status"`
}

type ProcurementItemRequest struct {
	ProductName string  `json:"productName"`
	Quantity    int     `json:"quantity"`
	Unit        string  `json:"unit"`
	UnitPrice   float64 `json:"unitPrice"`
	Amount      float64 `json:"amount"`
}

type SalesOrderRequest struct {
	CustomerName    string              `json:"customerName" binding:"required"`
	CustomerPhone   string              `json:"customerPhone"`
	CustomerAddress string              `json:"customerAddress"`
	PaymentMethod   string              `json:"paymentMethod"`
	DeliveryDate    string              `json:"deliveryDate"`
	Remark          string              `json:"remark"`
	Items           []SalesItemRequest  `json:"items"`
	TotalAmount     float64             `json:"totalAmount"`
	Status          string              `json:"status"`
	PaymentStatus   string              `json:"paymentStatus"`
}

type SalesItemRequest struct {
	ProductID   uint    `json:"productId"`
	ProductName string  `json:"productName"`
	Quantity    int     `json:"quantity"`
	UnitPrice   float64 `json:"unitPrice"`
	Amount      float64 `json:"amount"`
	Discount    float64 `json:"discount"`
}

type ProcurementHandler struct {
	service         *service.ProcurementService
	productRepo     *repository.ProductRepository
}

func NewProcurementHandler(s *service.ProcurementService, p *repository.ProductRepository) *ProcurementHandler {
	return &ProcurementHandler{
		service:     s,
		productRepo: p,
	}
}

func (h *ProcurementHandler) List(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}
	pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	if err != nil || pageSize < 1 {
		pageSize = 10
	}
	status := c.Query("status")

	orders, total, err := h.service.List(page, pageSize, status)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.PageSuccess(c, orders, total, page, pageSize)
}

func (h *ProcurementHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}
	order, err := h.service.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "采购订单不存在")
		return
	}
	response.Success(c, order)
}

func (h *ProcurementHandler) Create(c *gin.Context) {
	var req ProcurementOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 转换日期格式
	var expectedDate *model.Date
	if req.ExpectedDate != "" {
		parsed, err := time.Parse("2006-01-02", req.ExpectedDate)
		if err != nil {
			response.BadRequest(c, "日期格式错误: "+err.Error())
			return
		}
		expectedDate = &model.Date{Time: parsed}
	}

	// 转换 items - look up product ID by name
	items := make([]model.ProcurementItem, len(req.Items))
	for i, item := range req.Items {
		// Look up product by name
		product, err := h.productRepo.GetByName(item.ProductName)
		if err != nil {
			// If product not found, we can still save without ProductID
			// or return an error depending on business logic
			// For now, continue without ProductID (ProductID will be 0)
			items[i] = model.ProcurementItem{
				ProductName: item.ProductName,
				Quantity:    item.Quantity,
				Unit:        item.Unit,
				UnitPrice:   item.UnitPrice,
				Amount:      item.Amount,
			}
		} else {
			items[i] = model.ProcurementItem{
				ProductID:   product.ID,
				ProductName: item.ProductName,
				Quantity:    item.Quantity,
				Unit:        item.Unit,
				UnitPrice:   item.UnitPrice,
				Amount:      item.Amount,
			}
		}
	}

	order := model.ProcurementOrder{
		SupplierID:    req.SupplierID,
		Warehouse:     req.Warehouse,
		ExpectedDate:  expectedDate,
		Remark:        req.Remark,
		Items:         items,
		TotalAmount:   req.TotalAmount,
		Status:        req.Status,
		AttachmentURL: "",
	}

	if err := h.service.Create(&order); err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, "创建成功", order)
}

func (h *ProcurementHandler) UpdateStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}
	var req struct {
		Status        string `json:"status"`
		AttachmentURL string `json:"attachmentUrl"`
		Remark        string `json:"remark"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	if err := h.service.UpdateStatus(uint(id), req.Status, req.AttachmentURL, req.Remark); err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, "状态更新成功", nil)
}

func (h *ProcurementHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}
	if err := h.service.Delete(uint(id)); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessWithMessage(c, "删除成功", nil)
}

type SalesHandler struct {
	service *service.SalesService
}

func NewSalesHandler(s *service.SalesService) *SalesHandler {
	return &SalesHandler{service: s}
}

func (h *SalesHandler) List(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}
	pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	if err != nil || pageSize < 1 {
		pageSize = 10
	}
	status := c.Query("status")

	orders, total, err := h.service.List(page, pageSize, status)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.PageSuccess(c, orders, total, page, pageSize)
}

func (h *SalesHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}
	order, err := h.service.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "销售订单不存在")
		return
	}
	response.Success(c, order)
}

func (h *SalesHandler) Create(c *gin.Context) {
	var req SalesOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 转换日期格式
	var deliveryDate *time.Time
	if req.DeliveryDate != "" {
		parsed, err := time.Parse("2006-01-02", req.DeliveryDate)
		if err != nil {
			response.BadRequest(c, "日期格式错误: "+err.Error())
			return
		}
		deliveryDate = &parsed
	}

	// 转换 items
	items := make([]model.SalesOrderItem, len(req.Items))
	for i, item := range req.Items {
		items[i] = model.SalesOrderItem{
			ProductID:   item.ProductID,
			ProductName: item.ProductName,
			Quantity:    item.Quantity,
			UnitPrice:   item.UnitPrice,
			Amount:      item.Amount,
			Discount:    item.Discount,
		}
	}

	order := model.SalesOrder{
		CustomerName:    req.CustomerName,
		CustomerPhone:   req.CustomerPhone,
		CustomerAddress: req.CustomerAddress,
		PaymentMethod:   req.PaymentMethod,
		DeliveryDate:    deliveryDate,
		Remark:          req.Remark,
		Items:           items,
		TotalAmount:     req.TotalAmount,
		Status:          req.Status,
		PaymentStatus:   req.PaymentStatus,
	}

	if err := h.service.Create(&order); err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, "创建成功", order)
}

func (h *SalesHandler) UpdateStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}
	var req struct {
		Status string `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	if err := h.service.UpdateStatus(uint(id), req.Status); err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, "状态更新成功", nil)
}

func (h *SalesHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}
	if err := h.service.Delete(uint(id)); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessWithMessage(c, "删除成功", nil)
}

type LogisticsHandler struct {
	service *service.LogisticsService
}

func NewLogisticsHandler(s *service.LogisticsService) *LogisticsHandler {
	return &LogisticsHandler{service: s}
}

func (h *LogisticsHandler) List(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}
	pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	if err != nil || pageSize < 1 {
		pageSize = 10
	}
	status := c.Query("status")

	orders, total, err := h.service.List(page, pageSize, status)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.PageSuccess(c, orders, total, page, pageSize)
}

func (h *LogisticsHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}
	order, err := h.service.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "物流订单不存在")
		return
	}
	response.Success(c, order)
}

func (h *LogisticsHandler) Create(c *gin.Context) {
	var order model.LogisticsOrder
	if err := c.ShouldBindJSON(&order); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	if err := h.service.Create(&order); err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, "创建成功", order)
}

func (h *LogisticsHandler) UpdateStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}
	var req struct {
		Status   string `json:"status"`
		Location string `json:"location"`
		Desc     string `json:"description"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	if err := h.service.UpdateStatus(uint(id), req.Status); err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, "状态更新成功", nil)
}

func (h *LogisticsHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}
	if err := h.service.Delete(uint(id)); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessWithMessage(c, "删除成功", nil)
}
