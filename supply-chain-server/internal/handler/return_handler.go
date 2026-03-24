package handler

import (
	"strconv"
	"supply-chain-server/internal/model"
	"supply-chain-server/internal/service"
	"supply-chain-server/pkg/response"

	"github.com/gin-gonic/gin"
)

type SalesReturnHandler struct {
	service *service.SalesReturnService
}

func NewSalesReturnHandler(s *service.SalesReturnService) *SalesReturnHandler {
	return &SalesReturnHandler{service: s}
}

func (h *SalesReturnHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	status := c.Query("status")

	returns, total, err := h.service.List(page, pageSize, status)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.PageSuccess(c, returns, total, page, pageSize)
}

func (h *SalesReturnHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	ret, err := h.service.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "退货单不存在")
		return
	}

	response.Success(c, ret)
}

func (h *SalesReturnHandler) Create(c *gin.Context) {
	var req struct {
		SalesOrderID uint    `json:"salesOrderId"`
		SalesOrderNo string  `json:"salesOrderNo"`
		CustomerName string  `json:"customerName"`
		Reason       string  `json:"reason"`
		Remark       string  `json:"remark"`
		Items        []struct {
			ProductID   uint    `json:"productId"`
			ProductName string  `json:"productName"`
			Quantity    int     `json:"quantity"`
			UnitPrice   float64 `json:"unitPrice"`
			Reason      string  `json:"reason"`
		} `json:"items"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	ret := &model.SalesReturn{
		SalesOrderID: req.SalesOrderID,
		SalesOrderNo: req.SalesOrderNo,
		CustomerName: req.CustomerName,
		Reason:       req.Reason,
		Remark:       req.Remark,
	}

	for _, item := range req.Items {
		ret.Items = append(ret.Items, model.SalesReturnItem{
			ProductID:   item.ProductID,
			ProductName: item.ProductName,
			Quantity:    item.Quantity,
			UnitPrice:   item.UnitPrice,
			Reason:      item.Reason,
		})
	}

	if err := h.service.Create(ret); err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, "创建成功", ret)
}

func (h *SalesReturnHandler) Approve(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	if err := h.service.Approve(uint(id)); err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, "批准成功", nil)
}

func (h *SalesReturnHandler) Reject(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	var req struct {
		Reason string `json:"reason"`
	}
	c.ShouldBindJSON(&req)

	if err := h.service.Reject(uint(id), req.Reason); err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, "拒绝成功", nil)
}

func (h *SalesReturnHandler) Complete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	var req struct {
		RefundAmount float64 `json:"refundAmount"`
	}
	c.ShouldBindJSON(&req)

	if err := h.service.Complete(uint(id), req.RefundAmount); err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, "退货完成", nil)
}

func (h *SalesReturnHandler) Delete(c *gin.Context) {
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

type ProcurementReturnHandler struct {
	service *service.ProcurementReturnService
}

func NewProcurementReturnHandler(s *service.ProcurementReturnService) *ProcurementReturnHandler {
	return &ProcurementReturnHandler{service: s}
}

func (h *ProcurementReturnHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	status := c.Query("status")

	returns, total, err := h.service.List(page, pageSize, status)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.PageSuccess(c, returns, total, page, pageSize)
}

func (h *ProcurementReturnHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	ret, err := h.service.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "退货单不存在")
		return
	}

	response.Success(c, ret)
}

func (h *ProcurementReturnHandler) Create(c *gin.Context) {
	var req struct {
		ProcurementOrderID uint    `json:"procurementOrderId"`
		ProcurementOrderNo string  `json:"procurementOrderNo"`
		SupplierName       string  `json:"supplierName"`
		Reason             string  `json:"reason"`
		Remark             string  `json:"remark"`
		Items              []struct {
			ProductID   uint    `json:"productId"`
			ProductName string  `json:"productName"`
			Quantity    int     `json:"quantity"`
			UnitPrice   float64 `json:"unitPrice"`
			Reason      string  `json:"reason"`
		} `json:"items"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	ret := &model.ProcurementReturn{
		ProcurementOrderID: req.ProcurementOrderID,
		ProcurementOrderNo: req.ProcurementOrderNo,
		SupplierName:       req.SupplierName,
		Reason:             req.Reason,
		Remark:             req.Remark,
	}

	for _, item := range req.Items {
		ret.Items = append(ret.Items, model.ProcurementReturnItem{
			ProductID:   item.ProductID,
			ProductName: item.ProductName,
			Quantity:    item.Quantity,
			UnitPrice:   item.UnitPrice,
			Reason:      item.Reason,
		})
	}

	if err := h.service.Create(ret); err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, "创建成功", ret)
}

func (h *ProcurementReturnHandler) Approve(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	if err := h.service.Approve(uint(id)); err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, "批准成功", nil)
}

func (h *ProcurementReturnHandler) Reject(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	var req struct {
		Reason string `json:"reason"`
	}
	c.ShouldBindJSON(&req)

	if err := h.service.Reject(uint(id), req.Reason); err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, "拒绝成功", nil)
}

func (h *ProcurementReturnHandler) Complete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	var req struct {
		RefundAmount float64 `json:"refundAmount"`
	}
	c.ShouldBindJSON(&req)

	if err := h.service.Complete(uint(id), req.RefundAmount); err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, "退货完成", nil)
}

func (h *ProcurementReturnHandler) Delete(c *gin.Context) {
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