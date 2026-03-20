package handler

import (
	"strconv"
	"supply-chain-server/internal/model"
	"supply-chain-server/internal/service"
	"supply-chain-server/pkg/response"
	"supply-chain-server/pkg/utils"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService    *service.UserService
	profileService *service.UserProfileService
}

func NewUserHandler(userService *service.UserService, profileService *service.UserProfileService) *UserHandler {
	return &UserHandler{
		userService:    userService,
		profileService: profileService,
	}
}

func (h *UserHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	keyword := c.Query("keyword")

	users, total, err := h.userService.List(page, pageSize, keyword)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.PageSuccess(c, users, total, page, pageSize)
}

func (h *UserHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的用户ID")
		return
	}

	user, err := h.userService.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "用户不存在")
		return
	}

	// 获取用户角色
	roles, _ := h.profileService.GetUserRoles(uint(id))
	user.Roles = roles

	response.Success(c, user)
}

func (h *UserHandler) Create(c *gin.Context) {
	var req struct {
		Username    string `json:"username" binding:"required"`
		Password    string `json:"password" binding:"required"`
		Email       string `json:"email"`
		Phone       string `json:"phone"`
		Role        string `json:"role"`
		Department  string `json:"department"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	// 检查用户名是否已存在
	if h.userService.ExistsByUsername(req.Username) {
		response.BadRequest(c, "用户名已存在")
		return
	}

	// 验证角色
	role := req.Role
	if role == "" {
		role = "operator"
	}
	if role != "admin" && role != "manager" && role != "operator" {
		response.BadRequest(c, "无效的角色")
		return
	}

	// 加密密码
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		response.ServerError(c, "密码加密失败")
		return
	}

	user := &model.User{
		Username:   req.Username,
		Password:   hashedPassword,
		Email:      req.Email,
		Phone:      req.Phone,
		Role:       role,
		Department: req.Department,
		Status:     1,
	}

	if err := h.userService.Create(user); err != nil {
		response.ServerError(c, err.Error())
		return
	}

	// 创建用户扩展信息
	if req.Department != "" {
		currentUserID := c.GetUint("userId")
		profile := &model.UserProfile{
			UserID:     user.ID,
			Department: req.Department,
			CreatedBy:  &currentUserID,
		}
		h.profileService.CreateProfile(profile)
	}

	response.Success(c, user)
}

func (h *UserHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的用户ID")
		return
	}

	var req struct {
		Email       string `json:"email"`
		Phone       string `json:"phone"`
		Role        string `json:"role"`
		Department  string `json:"department"`
		Status      int8   `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	user, err := h.userService.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "用户不存在")
		return
	}

	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.Role != "" {
		// 验证角色
		if req.Role != "admin" && req.Role != "manager" && req.Role != "operator" {
			response.BadRequest(c, "无效的角色")
			return
		}
		user.Role = req.Role
	}
	if req.Department != "" {
		user.Department = req.Department

		// 更新用户扩展信息
		profile, err := h.profileService.GetProfile(uint(id))
		if err != nil {
			// 如果不存在则创建
			currentUserID := c.GetUint("userId")
			newProfile := &model.UserProfile{
				UserID:     user.ID,
				Department: req.Department,
				CreatedBy:  &currentUserID,
			}
			h.profileService.CreateProfile(newProfile)
		} else {
			// 如果存在则更新
			profile.Department = req.Department
			h.profileService.UpdateProfile(profile)
		}
	}
	if req.Status != 0 {
		user.Status = req.Status
	}
	// 允许将用户设为禁用状态
	if req.Status == 0 {
		user.Status = 0
	}

	if err := h.userService.Update(user); err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, user)
}

func (h *UserHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的用户ID")
		return
	}

	// 禁止删除 admin 用户
	user, err := h.userService.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "用户不存在")
		return
	}

	if user.Username == "admin" {
		response.BadRequest(c, "不能删除 admin 用户")
		return
	}

	if err := h.userService.Delete(uint(id)); err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, nil)
}

type ResetPasswordRequest struct {
	Password string `json:"password" binding:"required"`
}

func (h *UserHandler) ResetPassword(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的用户ID")
		return
	}

	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请输入新密码")
		return
	}

	if err := h.userService.ResetPassword(uint(id), req.Password); err != nil {
		response.Error(c, 400, err.Error())
		return
	}

	response.SuccessWithMessage(c, "密码重置成功", nil)
}

// SetRoles 设置用户角色
func (h *UserHandler) SetRoles(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的用户ID")
		return
	}

	var req struct {
		RoleIDs []uint `json:"roleIds" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	if err := h.profileService.SetUserRoles(uint(id), req.RoleIDs); err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, nil)
}

// GetRoles 获取用户角色
func (h *UserHandler) GetRoles(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的用户ID")
		return
	}

	roles, err := h.profileService.GetUserRoles(uint(id))
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, roles)
}

// GetPermissions 获取用户权限
func (h *UserHandler) GetPermissions(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的用户ID")
		return
	}

	permissions, err := h.userService.GetUserPermissions(uint(id))
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, permissions)
}

// SetPermissions 设置用户权限
func (h *UserHandler) SetPermissions(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的用户ID")
		return
	}

	var req struct {
		PermissionIDs []uint `json:"permissionIds" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	if err := h.userService.SetUserPermissions(uint(id), req.PermissionIDs); err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, nil)
}