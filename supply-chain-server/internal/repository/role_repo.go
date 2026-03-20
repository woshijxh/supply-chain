package repository

import (
	"supply-chain-server/internal/model"

	"gorm.io/gorm"
)

type RoleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) *RoleRepository {
	return &RoleRepository{db: db}
}

func (r *RoleRepository) Create(role *model.Role) error {
	return r.db.Create(role).Error
}

func (r *RoleRepository) GetByID(id uint) (*model.Role, error) {
	var role model.Role
	err := r.db.Preload("Permissions").First(&role, id).Error
	return &role, err
}

func (r *RoleRepository) GetByName(name string) (*model.Role, error) {
	var role model.Role
	err := r.db.Where("name = ?", name).Preload("Permissions").First(&role).Error
	return &role, err
}

func (r *RoleRepository) List(page, pageSize int, keyword string) ([]model.Role, int64, error) {
	var roles []model.Role
	var total int64

	query := r.db.Model(&model.Role{})

	if keyword != "" {
		query = query.Where("name LIKE ? OR description LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err = query.Preload("Permissions").Offset(offset).Limit(pageSize).Find(&roles).Error

	return roles, total, err
}

func (r *RoleRepository) Update(role *model.Role) error {
	return r.db.Save(role).Error
}

func (r *RoleRepository) Delete(id uint) error {
	return r.db.Delete(&model.Role{}, id).Error
}

func (r *RoleRepository) GetAll() ([]model.Role, error) {
	var roles []model.Role
	err := r.db.Preload("Permissions").Find(&roles).Error
	return roles, err
}

func (r *RoleRepository) AssignPermissions(roleID uint, permissionIDs []uint) error {
	var permissions []model.Permission
	if err := r.db.Find(&permissions, permissionIDs).Error; err != nil {
		return err
	}
	return r.db.Model(&model.Role{}).First(&model.Role{}, roleID).Error
}