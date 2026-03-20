package handler

import (
	"strconv"
	"supply-chain-server/internal/model"
	"supply-chain-server/internal/service"
	"supply-chain-server/pkg/response"

	"github.com/gin-gonic/gin"
)

type PermissionHandler struct {
	service *service.PermissionService
}

func NewPermissionHandler(service *service.PermissionService) *PermissionHandler {
	return &PermissionHandler{service: service}
}

func (h *PermissionHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	keyword := c.Query("keyword")

	permissions, total, err := h.service.List(page, pageSize, keyword)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.PageSuccess(c, permissions, total, page, pageSize)
}

func (h *PermissionHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的权限ID")
		return
	}

	permission, err := h.service.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "权限不存在")
		return
	}

	response.Success(c, permission)
}

func (h *PermissionHandler) Create(c *gin.Context) {
	var permission struct {
		Name        string `json:"name" binding:"required"`
		Code        string `json:"code"`
		Type        string `json:"type"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&permission); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	newPermission := &model.Permission{
		Name:        permission.Name,
		Code:        permission.Code,
		Type:        permission.Type,
		Description: permission.Description,
	}

	if err := h.service.Create(newPermission); err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, newPermission)
}

func (h *PermissionHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的权限ID")
		return
	}

	var req struct {
		Name        string `json:"name"`
		Code        string `json:"code"`
		Type        string `json:"type"`
		Description string `json:"description"`
		Status      int8   `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	permission, err := h.service.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "权限不存在")
		return
	}

	if req.Name != "" {
		permission.Name = req.Name
	}
	if req.Code != "" {
		permission.Code = req.Code
	}
	if req.Type != "" {
		permission.Type = req.Type
	}
	if req.Description != "" {
		permission.Description = req.Description
	}
	// 允许更新 status，包括设为 0（禁用）
	permission.Status = req.Status

	if err := h.service.Update(permission); err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, permission)
}

func (h *PermissionHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的权限ID")
		return
	}

	if err := h.service.Delete(uint(id)); err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, nil)
}

// GetAll 获取所有权限（不需要分页）
func (h *PermissionHandler) GetAll(c *gin.Context) {
	permissions, err := h.service.GetAll()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, permissions)
}