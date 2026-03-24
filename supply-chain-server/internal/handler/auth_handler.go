package handler

import (
	"supply-chain-server/internal/model"
	"supply-chain-server/internal/service"
	"supply-chain-server/pkg/captcha"
	"supply-chain-server/pkg/response"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service *service.UserService
}

func NewAuthHandler(s *service.UserService) *AuthHandler {
	return &AuthHandler{service: s}
}

type LoginRequest struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	CaptchaID string `json:"captchaId"`
	Captcha   string `json:"captcha"`
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"currentPassword" binding:"required"`
	NewPassword    string `json:"newPassword" binding:"required"`
}

type UpdateAvatarRequest struct {
	Avatar string `json:"avatar" binding:"required"`
}

type UserResponse struct {
	ID          uint            `json:"id"`
	Username    string          `json:"username"`
	Email       string          `json:"email"`
	Phone       string          `json:"phone"`
	Role        string          `json:"role"`
	Department  string          `json:"department"`
	Position    string          `json:"position"`
	Avatar      string          `json:"avatar"`
	Status      int8            `json:"status"`
	Roles       []model.Role    `json:"roles"`
	Permissions []model.Permission `json:"permissions"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	// Verify captcha - always required
	if req.CaptchaID == "" || req.Captcha == "" {
		response.Error(c, 400, "请输入验证码")
		return
	}

	if !captcha.Verify(req.CaptchaID, req.Captcha) {
		response.Error(c, 400, "验证码错误")
		return
	}

	token, user, err := h.service.LoginWithRoles(req.Username, req.Password)
	if err != nil {
		response.Error(c, 400, err.Error())
		return
	}

	// 构建用户响应，包含角色和权限
	userResp := UserResponse{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		Phone:       user.Phone,
		Role:        user.Role,
		Department:  user.Department,
		Position:    user.Position,
		Avatar:      user.Avatar,
		Status:      user.Status,
		Roles:       user.Roles,
		Permissions: user.Permissions,
	}

	response.Success(c, gin.H{
		"token": token,
		"user":  userResp,
	})
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	user := &model.User{
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
		Phone:    req.Phone,
		Role:     "operator",
		Status:   1,
	}

	if err := h.service.Create(user); err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.SuccessWithMessage(c, "注册成功", nil)
}

func (h *AuthHandler) Profile(c *gin.Context) {
	userID := c.GetUint("userId")
	user, err := h.service.GetByIDWithRoles(userID)
	if err != nil {
		response.NotFound(c, "用户不存在")
		return
	}

	// 构建用户响应，包含角色和权限
	userResp := UserResponse{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		Phone:       user.Phone,
		Role:        user.Role,
		Department:  user.Department,
		Position:    user.Position,
		Avatar:      user.Avatar,
		Status:      user.Status,
		Roles:       user.Roles,
		Permissions: user.Permissions,
	}

	response.Success(c, userResp)
}

func (h *AuthHandler) ChangePassword(c *gin.Context) {
	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	userID := c.GetUint("userId")
	err := h.service.ChangePassword(userID, req.CurrentPassword, req.NewPassword)
	if err != nil {
		response.Error(c, 400, err.Error())
		return
	}

	response.SuccessWithMessage(c, "密码修改成功", nil)
}

func (h *AuthHandler) UpdateAvatar(c *gin.Context) {
	var req UpdateAvatarRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	userID := c.GetUint("userId")
	err := h.service.UpdateAvatar(userID, req.Avatar)
	if err != nil {
		response.Error(c, 400, err.Error())
		return
	}

	response.SuccessWithMessage(c, "头像更新成功", nil)
}
