package handler

import (
	"strconv"
	"supply-chain-server/internal/model"
	"supply-chain-server/internal/service"
	"supply-chain-server/pkg/response"

	"github.com/gin-gonic/gin"
)

type RoleHandler struct {
	service     *service.RoleService
	permService *service.PermissionService
}

func NewRoleHandler(roleService *service.RoleService, permService *service.PermissionService) *RoleHandler {
	return &RoleHandler{
		service:     roleService,
		permService: permService,
	}
}

func (h *RoleHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	keyword := c.Query("keyword")

	roles, total, err := h.service.List(page, pageSize, keyword)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.PageSuccess(c, roles, total, page, pageSize)
}

func (h *RoleHandler) All(c *gin.Context) {
	roles, err := h.service.GetAll()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, roles)
}

func (h *RoleHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的角色ID")
		return
	}

	role, err := h.service.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "角色不存在")
		return
	}

	response.Success(c, role)
}

func (h *RoleHandler) Create(c *gin.Context) {
	var role struct {
		Name        string `json:"name" binding:"required"`
		Code        string `json:"code" binding:"required"`
		Description string `json:"description"`
		Status      int8   `json:"status"`
	}

	if err := c.ShouldBindJSON(&role); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	newRole := &model.Role{
		Name:        role.Name,
		Code:        role.Code,
		Description: role.Description,
		Status:      role.Status,
	}
	if newRole.Status == 0 {
		newRole.Status = 1
	}

	if err := h.service.Create(newRole); err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, newRole)
}

func (h *RoleHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的角色ID")
		return
	}

	var role struct {
		Name        string `json:"name"`
		Code        string `json:"code"`
		Description string `json:"description"`
		Status      int8   `json:"status"`
	}

	if err := c.ShouldBindJSON(&role); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	existingRole, err := h.service.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "角色不存在")
		return
	}

	if role.Name != "" {
		existingRole.Name = role.Name
	}
	if role.Code != "" {
		existingRole.Code = role.Code
	}
	if role.Description != "" {
		existingRole.Description = role.Description
	}
	existingRole.Status = role.Status

	if err := h.service.Update(existingRole); err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, existingRole)
}

func (h *RoleHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的角色ID")
		return
	}

	// 禁止删除 admin 角色
	role, err := h.service.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "角色不存在")
		return
	}

	if role.Name == "admin" {
		response.BadRequest(c, "不能删除 admin 角色")
		return
	}

	if err := h.service.Delete(uint(id)); err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, nil)
}

// SetPermissions 设置角色权限
func (h *RoleHandler) SetPermissions(c *gin.Context) {
	roleID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的角色ID")
		return
	}

	var req struct {
		PermissionIDs []uint `json:"permissionIds" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	// 调用 service 保存权限关联
	if err := h.service.SetPermissions(uint(roleID), req.PermissionIDs, h.permService.GetRepo()); err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, nil)
}

// GetPermissions 获取角色权限
func (h *RoleHandler) GetPermissions(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的角色ID")
		return
	}

	permissions, err := h.permService.GetByRoleID(uint(id))
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, permissions)
}