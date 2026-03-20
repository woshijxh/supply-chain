package repository

import (
	"supply-chain-server/internal/model"

	"gorm.io/gorm"
)

type PermissionRepository struct {
	db *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) *PermissionRepository {
	return &PermissionRepository{db: db}
}

func (r *PermissionRepository) Create(permission *model.Permission) error {
	return r.db.Create(permission).Error
}

func (r *PermissionRepository) GetByID(id uint) (*model.Permission, error) {
	var permission model.Permission
	err := r.db.First(&permission, id).Error
	return &permission, err
}

func (r *PermissionRepository) GetByName(name string) (*model.Permission, error) {
	var permission model.Permission
	err := r.db.Where("name = ?", name).First(&permission).Error
	return &permission, err
}

func (r *PermissionRepository) List(page, pageSize int, keyword string) ([]model.Permission, int64, error) {
	var permissions []model.Permission
	var total int64

	query := r.db.Model(&model.Permission{})

	if keyword != "" {
		query = query.Where("name LIKE ? OR description LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err = query.Offset(offset).Limit(pageSize).Find(&permissions).Error

	return permissions, total, err
}

func (r *PermissionRepository) GetAll() ([]model.Permission, error) {
	var permissions []model.Permission
	err := r.db.Find(&permissions).Error
	return permissions, err
}

func (r *PermissionRepository) Delete(id uint) error {
	return r.db.Delete(&model.Permission{}, id).Error
}

func (r *PermissionRepository) Update(permission *model.Permission) error {
	return r.db.Save(permission).Error
}

func (r *PermissionRepository) GetByRoleID(roleID uint) ([]model.Permission, error) {
	var permissions []model.Permission
	err := r.db.Model(&model.Permission{}).
		Joins("JOIN role_permissions ON permissions.id = role_permissions.permission_id").
		Where("role_permissions.role_id = ?", roleID).
		Find(&permissions).Error
	return permissions, err
}

// SetRolePermissions 设置角色的权限
func (r *PermissionRepository) SetRolePermissions(roleID uint, permissionIDs []uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 删除角色已有的所有权限关联
		if err := tx.Exec("DELETE FROM role_permissions WHERE role_id = ?", roleID).Error; err != nil {
			return err
		}

		// 添加新的权限关联 - 使用 raw SQL 因为 GORM 的 many2many 操作比较复杂
		if len(permissionIDs) > 0 {
			for _, permID := range permissionIDs {
				if err := tx.Exec("INSERT INTO role_permissions (role_id, permission_id) VALUES (?, ?) ON DUPLICATE KEY UPDATE permission_id = permission_id", roleID, permID).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})
}