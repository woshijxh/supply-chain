package handler

import (
	"strconv"
	"supply-chain-server/internal/model"
	"supply-chain-server/internal/service"
	"supply-chain-server/pkg/response"

	"github.com/gin-gonic/gin"
)

type SupplierHandler struct {
	service *service.SupplierService
}

func NewSupplierHandler(s *service.SupplierService) *SupplierHandler {
	return &SupplierHandler{service: s}
}

func (h *SupplierHandler) List(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}
	pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	if err != nil || pageSize < 1 {
		pageSize = 10
	}
	keyword := c.Query("keyword")

	suppliers, total, err := h.service.List(page, pageSize, keyword)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.PageSuccess(c, suppliers, total, page, pageSize)
}

func (h *SupplierHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}
	supplier, err := h.service.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "供应商不存在")
		return
	}
	response.Success(c, supplier)
}

func (h *SupplierHandler) Create(c *gin.Context) {
	var supplier model.Supplier
	if err := c.ShouldBindJSON(&supplier); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	// 确保前端传来的 code 被清空，由后端自动生成
	supplier.Code = ""

	if err := h.service.Create(&supplier); err != nil {
		response.ServerError(c, err.Error())
		return
	}

	// 重新查询获取完整数据（包含生成的code）
	created, _ := h.service.GetByID(supplier.ID)
	if created != nil {
		response.SuccessWithMessage(c, "创建成功", created)
	} else {
		response.SuccessWithMessage(c, "创建成功", supplier)
	}
}

func (h *SupplierHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}
	var supplier model.Supplier
	if err := c.ShouldBindJSON(&supplier); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	// 编码不允许修改，清空以使用原有值
	supplier.Code = ""

	supplier.ID = uint(id)
	if err := h.service.Update(&supplier); err != nil {
		response.ServerError(c, err.Error())
		return
	}

	// 重新查询获取完整数据
	updated, _ := h.service.GetByID(uint(id))
	if updated != nil {
		response.SuccessWithMessage(c, "更新成功", updated)
	} else {
		response.SuccessWithMessage(c, "更新成功", supplier)
	}
}

func (h *SupplierHandler) Delete(c *gin.Context) {
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
