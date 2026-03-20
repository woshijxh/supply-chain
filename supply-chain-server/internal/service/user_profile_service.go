package service

import (
	"errors"
	"supply-chain-server/internal/model"
	"supply-chain-server/internal/repository"
)

type UserProfileService struct {
	profileRepo *repository.UserProfileRepository
	userRepo    *repository.UserRepository
}

func NewUserProfileService(profileRepo *repository.UserProfileRepository, userRepo *repository.UserRepository) *UserProfileService {
	return &UserProfileService{
		profileRepo: profileRepo,
		userRepo:    userRepo,
	}
}

func (s *UserProfileService) CreateProfile(profile *model.UserProfile) error {
	return s.profileRepo.Create(profile)
}

func (s *UserProfileService) GetProfile(userID uint) (*model.UserProfile, error) {
	return s.profileRepo.GetByUserID(userID)
}

func (s *UserProfileService) UpdateProfile(profile *model.UserProfile) error {
	return s.profileRepo.Update(profile)
}

// AssignRole 为用户分配角色
func (s *UserProfileService) AssignRole(userID, roleID uint) error {
	// 检查用户是否存在
	_, err := s.userRepo.GetByID(userID)
	if err != nil {
		return errors.New("用户不存在")
	}
	return s.profileRepo.AssignRole(userID, roleID)
}

// RemoveRole 移除用户的角色
func (s *UserProfileService) RemoveRole(userID, roleID uint) error {
	return s.profileRepo.RemoveRole(userID, roleID)
}

// GetUserRoles 获取用户的所有角色
func (s *UserProfileService) GetUserRoles(userID uint) ([]model.Role, error) {
	return s.profileRepo.GetUserRoles(userID)
}

// SetUserRoles 设置用户的所有角色
func (s *UserProfileService) SetUserRoles(userID uint, roleIDs []uint) error {
	// 检查用户是否存在
	_, err := s.userRepo.GetByID(userID)
	if err != nil {
		return errors.New("用户不存在")
	}
	return s.profileRepo.SetUserRoles(userID, roleIDs)
}

// GetUserRoleNames 获取用户的所有角色名称
func (s *UserProfileService) GetUserRoleNames(userID uint) ([]string, error) {
	return s.profileRepo.GetUserRoleNames(userID)
}

// HasRole 检查用户是否拥有指定角色
func (s *UserProfileService) HasRole(userID uint, roleID uint) bool {
	return s.profileRepo.HasRole(userID, roleID)
}