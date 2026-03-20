package repository

import (
	"supply-chain-server/internal/model"

	"gorm.io/gorm"
)

type UserProfileRepository struct {
	db *gorm.DB
}

func NewUserProfileRepository(db *gorm.DB) *UserProfileRepository {
	return &UserProfileRepository{db: db}
}

func (r *UserProfileRepository) Create(profile *model.UserProfile) error {
	return r.db.Create(profile).Error
}

func (r *UserProfileRepository) GetByUserID(userID uint) (*model.UserProfile, error) {
	var profile model.UserProfile
	err := r.db.First(&profile, userID).Error
	return &profile, err
}

func (r *UserProfileRepository) Update(profile *model.UserProfile) error {
	return r.db.Save(profile).Error
}

func (r *UserProfileRepository) Delete(userID uint) error {
	return r.db.Delete(&model.UserProfile{}, userID).Error
}

// ========== 用户角色关联 ==========

// AssignRole 为用户分配角色
func (r *UserProfileRepository) AssignRole(userID uint, roleID uint) error {
	userRole := model.UserRole{
		UserID: userID,
		RoleID: roleID,
	}
	return r.db.Create(&userRole).Error
}

// RemoveRole 移除用户的角色
func (r *UserProfileRepository) RemoveRole(userID uint, roleID uint) error {
	return r.db.Where("user_id = ? AND role_id = ?", userID, roleID).Delete(&model.UserRole{}).Error
}

// GetUserRoles 获取用户的所有角色
func (r *UserProfileRepository) GetUserRoles(userID uint) ([]model.Role, error) {
	var roles []model.Role
	err := r.db.Model(&model.Role{}).
		Joins("JOIN user_roles ON roles.id = user_roles.role_id").
		Where("user_roles.user_id = ?", userID).
		Find(&roles).Error
	return roles, err
}

// SetUserRoles 设置用户的所有角色（先删除后添加）
func (r *UserProfileRepository) SetUserRoles(userID uint, roleIDs []uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 删除用户已有的所有角色
		if err := tx.Where("user_id = ?", userID).Delete(&model.UserRole{}).Error; err != nil {
			return err
		}

		// 添加新角色
		for _, roleID := range roleIDs {
			userRole := model.UserRole{
				UserID: userID,
				RoleID: roleID,
			}
			if err := tx.Create(&userRole).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// HasRole 检查用户是否拥有指定角色
func (r *UserProfileRepository) HasRole(userID uint, roleID uint) bool {
	var count int64
	r.db.Model(&model.UserRole{}).Where("user_id = ? AND role_id = ?", userID, roleID).Count(&count)
	return count > 0
}

// GetUserRoleNames 获取用户的所有角色名称
func (r *UserProfileRepository) GetUserRoleNames(userID uint) ([]string, error) {
	var roleNames []string
	err := r.db.Table("user_roles").
		Joins("JOIN roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ?", userID).
		Pluck("roles.name", &roleNames).Error
	return roleNames, err
}