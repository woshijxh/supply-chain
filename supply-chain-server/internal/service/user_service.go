package service

import (
	"errors"
	"fmt"
	"supply-chain-server/internal/model"
	"supply-chain-server/internal/repository"
	"supply-chain-server/pkg/utils"
)

type UserWithRoles struct {
	*model.User
	Roles       []model.Role       `json:"roles"`
	Permissions []model.Permission `json:"permissions"`
}

func (u *UserWithRoles) GetRoleNames() []string {
	names := make([]string, 0)
	for _, role := range u.Roles {
		names = append(names, role.Name)
	}
	return names
}

func (u *UserWithRoles) GetPermissionNames() []string {
	names := make([]string, 0)
	// 收集角色权限
	for _, role := range u.Roles {
		for _, perm := range role.Permissions {
			names = append(names, perm.Name)
		}
	}
	// 去重
	seen := make(map[string]bool)
	result := make([]string, 0)
	for _, n := range names {
		if !seen[n] {
			seen[n] = true
			result = append(result, n)
		}
	}
	return result
}

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(r *repository.UserRepository) *UserService {
	return &UserService{repo: r}
}

func (s *UserService) Create(user *model.User) error {
	if s.repo.ExistsByUsername(user.Username) {
		return errors.New("用户名已存在")
	}
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword
	return s.repo.Create(user)
}

func (s *UserService) GetByID(id uint) (*model.User, error) {
	return s.repo.GetByID(id)
}

// GetByIDWithRoles 获取用户信息及角色权限
func (s *UserService) GetByIDWithRoles(id uint) (*UserWithRoles, error) {
	user, err := s.repo.GetByIDWithRoles(id)
	if err != nil {
		return nil, err
	}
	return s.convertToUserWithRoles(user), nil
}

// LoginWithRoles 登录并返回用户角色权限信息
func (s *UserService) LoginWithRoles(username, password string) (string, *UserWithRoles, error) {
	user, err := s.repo.GetByUsername(username)
	if err != nil {
		return "", nil, errors.New("用户不存在")
	}

	if !utils.CheckPassword(password, user.Password) {
		return "", nil, errors.New("密码错误")
	}

	token, err := utils.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		return "", nil, errors.New("生成token失败")
	}

	userWithRoles := s.convertToUserWithRoles(user)
	return token, userWithRoles, nil
}

func (s *UserService) ExistsByUsername(username string) bool {
	return s.repo.ExistsByUsername(username)
}

func (s *UserService) List(page, pageSize int, keyword string) ([]UserWithRoles, int64, error) {
	users, total, err := s.repo.List(page, pageSize, keyword)
	if err != nil {
		return nil, 0, err
	}

	result := make([]UserWithRoles, 0, len(users))
	for _, user := range users {
		result = append(result, UserWithRoles{
			User:        &user,
			Roles:       user.Roles,
			Permissions: nil, // 权限从角色中获取
		})
	}

	return result, total, nil
}

func (s *UserService) Update(user *model.User) error {
	return s.repo.Update(user)
}

func (s *UserService) Delete(id uint) error {
	return s.repo.Delete(id)
}

// ResetPassword 重置用户密码（管理员用）
func (s *UserService) ResetPassword(userID uint, newPassword string) error {
	user, err := s.repo.GetByID(userID)
	if err != nil {
		return fmt.Errorf("用户不存在")
	}

	// 加密新密码
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("密码加密失败")
	}

	user.Password = hashedPassword
	return s.repo.Update(user)
}

// GetUserPermissions 获取用户的直接权限
func (s *UserService) GetUserPermissions(userID uint) ([]model.Permission, error) {
	return s.repo.GetUserPermissions(userID)
}

// SetUserPermissions 设置用户的直接权限
func (s *UserService) SetUserPermissions(userID uint, permissionIDs []uint) error {
	return s.repo.SetUserPermissions(userID, permissionIDs)
}

// ChangePassword 修改用户密码
func (s *UserService) ChangePassword(userID uint, currentPassword, newPassword string) error {
	user, err := s.repo.GetByID(userID)
	if err != nil {
		return fmt.Errorf("用户不存在")
	}

	// 验证当前密码
	if !utils.CheckPassword(currentPassword, user.Password) {
		return fmt.Errorf("当前密码错误")
	}

	// 加密新密码
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("密码加密失败")
	}

	user.Password = hashedPassword
	return s.repo.Update(user)
}

// UpdateAvatar 更新用户头像
func (s *UserService) UpdateAvatar(userID uint, avatar string) error {
	user, err := s.repo.GetByID(userID)
	if err != nil {
		return fmt.Errorf("用户不存在")
	}

	user.Avatar = avatar
	return s.repo.Update(user)
}

// convertToUserWithRoles 将User模型转换为带角色权限的用户信息
func (s *UserService) convertToUserWithRoles(user *model.User) *UserWithRoles {
	var roles []model.Role
	var perms []model.Permission

	// 预加载角色和权限
	s.repo.GetByIDWithRoles(user.ID)

	roles = user.Roles
	// 从角色中获取权限
	for _, role := range roles {
		perms = append(perms, role.Permissions...)
	}

	return &UserWithRoles{
		User:        user,
		Roles:       roles,
		Permissions: perms,
	}
}
