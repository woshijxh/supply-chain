package middleware

import (
	"strconv"
	"supply-chain-server/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RowFilter 行级权限过滤中间件
// 用于根据用户角色过滤查询数据
func RowFilter(resource string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取用户信息
		userID, exists := c.Get("userId")
		if !exists {
			response.Unauthorized(c, "请先登录")
			c.Abort()
			return
		}

		role, exists := c.Get("role")
		if !exists {
			role = "operator" // 默认角色
		}
		roleStr := role.(string)

		// 获取用户部门信息
		department := ""
		if dept, exists := c.Get("department"); exists {
			department = dept.(string)
		}

		// 将权限信息存储到 context 中，供后续 handler 使用
		c.Set("rowFilter", &RowFilterInfo{
			UserID:     userID.(uint),
			Role:       roleStr,
			Department: department,
			Resource:   resource,
		})

		c.Next()
	}
}

// RowFilterInfo 行级过滤信息
type RowFilterInfo struct {
	UserID     uint
	Role       string
	Department string
	Resource   string
}

// GetRowFilter 获取行级过滤信息
func GetRowFilter(c *gin.Context) *RowFilterInfo {
	if filter, exists := c.Get("rowFilter"); exists {
		return filter.(*RowFilterInfo)
	}
	return nil
}

// ApplySupplierFilter 应用供应商过滤
func ApplySupplierFilter(db *gorm.DB, c *gin.Context) *gorm.DB {
	filter := GetRowFilter(c)
	if filter == nil {
		return db
	}

	// admin 无过滤
	if filter.Role == "admin" {
		return db
	}

	// manager 按部门过滤
	if filter.Role == "manager" && filter.Department != "" {
		return db.Where("department = ?", filter.Department)
	}

	// operator 只看自己创建的数据
	return db.Where("created_by = ?", filter.UserID)
}

// ApplyProcurementFilter 应用采购订单过滤
func ApplyProcurementFilter(db *gorm.DB, c *gin.Context) *gorm.DB {
	filter := GetRowFilter(c)
	if filter == nil {
		return db
	}

	// admin 无过滤
	if filter.Role == "admin" {
		return db
	}

	// manager 按部门过滤
	if filter.Role == "manager" && filter.Department != "" {
		return db.Where("department = ?", filter.Department)
	}

	// operator 只看自己创建的数据
	return db.Where("created_by = ?", filter.UserID)
}

// ApplySalesFilter 应用销售订单过滤
func ApplySalesFilter(db *gorm.DB, c *gin.Context) *gorm.DB {
	filter := GetRowFilter(c)
	if filter == nil {
		return db
	}

	// admin 无过滤
	if filter.Role == "admin" {
		return db
	}

	// manager 按部门过滤
	if filter.Role == "manager" && filter.Department != "" {
		return db.Where("department = ?", filter.Department)
	}

	// operator 只看自己创建的数据
	return db.Where("created_by = ?", filter.UserID)
}

// ApplyLogisticsFilter 应用物流订单过滤
func ApplyLogisticsFilter(db *gorm.DB, c *gin.Context) *gorm.DB {
	filter := GetRowFilter(c)
	if filter == nil {
		return db
	}

	// admin 无过滤
	if filter.Role == "admin" {
		return db
	}

	// manager 按部门过滤
	if filter.Role == "manager" && filter.Department != "" {
		return db.Where("department = ?", filter.Department)
	}

	// operator 只看自己创建的数据
	return db.Where("created_by = ?", filter.UserID)
}

// SetCreatedByAndDepartment 设置创建人和部门（在创建数据时自动设置）
func SetCreatedByAndDepartment(userID uint, role, department string) (createdBy *uint, dept string) {
	createdBy = &userID
	return createdBy, department
}

// GetUserDepartment 获取用户部门
func GetUserDepartment(c *gin.Context) string {
	if dept, exists := c.Get("department"); exists {
		return dept.(string)
	}
	return ""
}

// GetUserIDFromContext 从 context 获取用户ID
func GetUserIDFromContext(c *gin.Context) uint {
	if userID, exists := c.Get("userId"); exists {
		return userID.(uint)
	}
	return 0
}

// ParseUint 解析 uint
func ParseUint(s string) uint {
	u, _ := strconv.ParseUint(s, 10, 32)
	return uint(u)
}