package handler

import (
	"strconv"
	"supply-chain-server/internal/model"
	"supply-chain-server/internal/service"
	"supply-chain-server/pkg/response"

	"github.com/gin-gonic/gin"
)

type CustomerHandler struct {
	service *service.CustomerService
}

func NewCustomerHandler(s *service.CustomerService) *CustomerHandler {
	return &CustomerHandler{service: s}
}

func (h *CustomerHandler) List(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}
	pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	if err != nil || pageSize < 1 {
		pageSize = 10
	}
	keyword := c.Query("keyword")

	customers, total, err := h.service.List(page, pageSize, keyword)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.PageSuccess(c, customers, total, page, pageSize)
}

func (h *CustomerHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	customer, err := h.service.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "客户不存在")
		return
	}

	response.Success(c, customer)
}

func (h *CustomerHandler) Create(c *gin.Context) {
	var req struct {
		Name    string `json:"name" binding:"required"`
		Contact string `json:"contact"`
		Phone   string `json:"phone"`
		Email   string `json:"email"`
		Address string `json:"address"`
		Level   string `json:"level"`
		Source  string `json:"source"`
		Status  int8   `json:"status"`
		Remark  string `json:"remark"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	customer := &model.Customer{
		Name:    req.Name,
		Contact: req.Contact,
		Phone:   req.Phone,
		Email:   req.Email,
		Address: req.Address,
		Level:   req.Level,
		Source:  req.Source,
		Status:  req.Status,
		Remark:  req.Remark,
	}

	if customer.Level == "" {
		customer.Level = "C"
	}
	if customer.Status == 0 {
		customer.Status = 1
	}

	if err := h.service.Create(customer); err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, "创建成功", customer)
}

func (h *CustomerHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	var req struct {
		Name    string `json:"name"`
		Contact string `json:"contact"`
		Phone   string `json:"phone"`
		Email   string `json:"email"`
		Address string `json:"address"`
		Level   string `json:"level"`
		Source  string `json:"source"`
		Status  int8   `json:"status"`
		Remark  string `json:"remark"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	customer, err := h.service.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "客户不存在")
		return
	}

	if req.Name != "" {
		customer.Name = req.Name
	}
	customer.Contact = req.Contact
	customer.Phone = req.Phone
	customer.Email = req.Email
	customer.Address = req.Address
	customer.Level = req.Level
	customer.Source = req.Source
	customer.Status = req.Status
	customer.Remark = req.Remark

	if err := h.service.Update(customer); err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, "更新成功", customer)
}

func (h *CustomerHandler) Delete(c *gin.Context) {
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

func (h *CustomerHandler) All(c *gin.Context) {
	customers, err := h.service.GetActiveCustomers()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, customers)
}