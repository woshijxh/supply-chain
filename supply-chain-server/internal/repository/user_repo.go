package repository

import (
	"supply-chain-server/internal/model"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) GetByID(id uint) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, id).Error
	return &user, err
}

// GetByIDWithRoles 获取用户及角色权限
func (r *UserRepository) GetByIDWithRoles(id uint) (*model.User, error) {
	var user model.User
	err := r.db.Preload("Roles.Permissions").First(&user, id).Error
	return &user, err
}

func (r *UserRepository) GetByUsername(username string) (*model.User, error) {
	var user model.User
	err := r.db.Preload("Roles.Permissions").Where("username = ?", username).First(&user).Error
	return &user, err
}

func (r *UserRepository) ExistsByUsername(username string) bool {
	var count int64
	r.db.Model(&model.User{}).Where("username = ?", username).Count(&count)
	return count > 0
}

func (r *UserRepository) List(page, pageSize int, keyword string) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	query := r.db.Model(&model.User{}).Preload("Roles").Preload("Roles.Permissions")

	if keyword != "" {
		query = query.Where("username LIKE ? OR email LIKE ? OR phone LIKE ?", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err = query.Offset(offset).Limit(pageSize).Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	// 加载每个用户的直接权限
	for i := range users {
		var directPerms []model.Permission
		r.db.Model(&model.Permission{}).
			Joins("JOIN user_permissions ON permissions.id = user_permissions.permission_id").
			Where("user_permissions.user_id = ?", users[i].ID).
			Find(&directPerms)
		// 将直接权限添加到用户模型中（临时使用 Roles 字段存储）
		// 这里我们不做处理，让 service 层处理
	}

	return users, total, err
}

func (r *UserRepository) Update(user *model.User) error {
	return r.db.Save(user).Error
}

func (r *UserRepository) Delete(id uint) error {
	return r.db.Delete(&model.User{}, id).Error
}

// GetUserPermissions 获取用户的直接权限
func (r *UserRepository) GetUserPermissions(userID uint) ([]model.Permission, error) {
	var permissions []model.Permission
	err := r.db.Model(&model.Permission{}).
		Joins("JOIN user_permissions ON permissions.id = user_permissions.permission_id").
		Where("user_permissions.user_id = ?", userID).
		Find(&permissions).Error
	return permissions, err
}

// SetUserPermissions 设置用户的直接权限
func (r *UserRepository) SetUserPermissions(userID uint, permissionIDs []uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 删除用户已有的所有直接权限
		if err := tx.Where("user_id = ?", userID).Delete(&model.UserPermission{}).Error; err != nil {
			return err
		}

		// 添加新权限
		for _, permID := range permissionIDs {
			userPerm := model.UserPermission{
				UserID:       userID,
				PermissionID: permID,
			}
			if err := tx.Create(&userPerm).Error; err != nil {
				return err
			}
		}

		return nil
	})
}
