package casbin

import (
	"strings"

	"github.com/casbin/casbin/v2"
	"gorm.io/gorm"
)

var enforcer *casbin.Enforcer

// InitCasbin 初始化 Casbin
func InitCasbin(db *gorm.DB) error {
	// 使用内存适配器
	enforcer, _ = casbin.NewEnforcer("pkg/casbin/rbac_model.conf")

	// 从数据库加载权限策略
	var permissions []struct {
		Sub string
		Obj string
		Act string
	}

	// 查询 role_permissions
	db.Table("role_permissions rp").
		Select("r.name as sub, p.name as obj, 'read' as act").
		Joins("JOIN roles r ON r.id = rp.role_id").
		Joins("JOIN permissions p ON p.id = rp.permission_id").
		Where("p.name LIKE '%:read'").Scan(&permissions)

	for _, p := range permissions {
		enforcer.AddPolicy(p.Sub, strings.TrimSuffix(p.Obj, ":read"), "read")
	}

	db.Table("role_permissions rp").
		Select("r.name as sub, p.name as obj, 'write' as act").
		Joins("JOIN roles r ON r.id = rp.role_id").
		Joins("JOIN permissions p ON p.id = rp.permission_id").
		Where("p.name LIKE '%:write'").Scan(&permissions)

	for _, p := range permissions {
		enforcer.AddPolicy(p.Sub, strings.TrimSuffix(p.Obj, ":write"), "write")
	}

	// 加载策略
	_ = enforcer.LoadPolicy()

	return nil
}

// GetEnforcer 获取 Casbin 执行器
func GetEnforcer() *casbin.Enforcer {
	return enforcer
}

// CheckPermission 检查用户是否有权限
func CheckPermission(sub, obj, act string) (bool, error) {
	if enforcer == nil {
		return false, nil
	}
	return enforcer.Enforce(sub, obj, act)
}

// AddRoleForUser 为用户添加角色
func AddRoleForUser(user, role string) error {
	if enforcer == nil {
		return nil
	}
	_, err := enforcer.AddRoleForUser(user, role)
	return err
}

// RemoveRoleForUser 移除用户的角色
func RemoveRoleForUser(user, role string) error {
	if enforcer == nil {
		return nil
	}
	_, err := enforcer.DeleteRoleForUser(user, role)
	return err
}

// GetRolesForUser 获取用户的所有角色
func GetRolesForUser(user string) ([]string, error) {
	if enforcer == nil {
		return nil, nil
	}
	return enforcer.GetRolesForUser(user)
}

// AddPolicy 添加策略
func AddPolicy(sub, obj, act string) error {
	if enforcer == nil {
		return nil
	}
	_, err := enforcer.AddPolicy(sub, obj, act)
	return err
}

// RemovePolicy 移除策略
func RemovePolicy(sub, obj, act string) error {
	if enforcer == nil {
		return nil
	}
	_, err := enforcer.RemovePolicy(sub, obj, act)
	return err
}

// GetPermissionsForUser 获取用户的所有权限
func GetPermissionsForUser(user string) [][]string {
	if enforcer == nil {
		return nil
	}
	permissions, _ := enforcer.GetPermissionsForUser(user)
	return permissions
}
