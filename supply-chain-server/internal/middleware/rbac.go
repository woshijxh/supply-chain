package middleware

import (
	"strconv"
	"supply-chain-server/pkg/casbin"
	"supply-chain-server/pkg/response"

	"github.com/gin-gonic/gin"
)

// RBACAuth RBAC 权限检查中间件
func RBACAuth(obj, act string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取用户角色
		role, exists := c.Get("role")
		if !exists {
			response.Forbidden(c, "未获取到用户角色")
			c.Abort()
			return
		}

		roleStr, ok := role.(string)
		if !ok {
			response.Forbidden(c, "用户角色格式错误")
			c.Abort()
			return
		}

		// 管理员拥有所有权限
		if roleStr == "admin" {
			c.Next()
			return
		}

		// 获取用户名
		username, exists := c.Get("username")
		if !exists {
			response.Forbidden(c, "未获取到用户名")
			c.Abort()
			return
		}

		usernameStr, ok := username.(string)
		if !ok {
			response.Forbidden(c, "用户名格式错误")
			c.Abort()
			return
		}

		// 获取用户ID
		userID, exists := c.Get("userId")
		if !exists {
			response.Forbidden(c, "未获取到用户ID")
			c.Abort()
			return
		}

		userIDStr := strconv.FormatUint(uint64(userID.(uint)), 10)

		// 检查用户是否有权限（先检查角色，再检查用户）
		hasPermission, err := casbin.CheckPermission(roleStr, obj, act)
		if err == nil && hasPermission {
			c.Next()
			return
		}

		// 检查用户直接权限
		hasPermission, err = casbin.CheckPermission(usernameStr, obj, act)
		if err == nil && hasPermission {
			c.Next()
			return
		}

		// 检查用户ID权限
		hasPermission, err = casbin.CheckPermission(userIDStr, obj, act)
		if err == nil && hasPermission {
			c.Next()
			return
		}

		response.Forbidden(c, "权限不足")
		c.Abort()
	}
}

// RequirePermission 要求特定权限
func RequirePermission(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		username, exists := c.Get("username")
		if !exists {
			response.Forbidden(c, "未获取到用户名")
			c.Abort()
			return
		}

		role, exists := c.Get("role")
		if !exists {
			response.Forbidden(c, "未获取到用户角色")
			c.Abort()
			return
		}

		// 管理员拥有所有权限
		if role == "admin" {
			c.Next()
			return
		}

		// 检查角色是否有权限
		hasPermission, _ := casbin.CheckPermission(role.(string), permission, "access")
		if hasPermission {
			c.Next()
			return
		}

		// 检查用户是否有权限
		hasPermission, _ = casbin.CheckPermission(username.(string), permission, "access")
		if hasPermission {
			c.Next()
			return
		}

		response.Forbidden(c, "权限不足")
		c.Abort()
	}
}