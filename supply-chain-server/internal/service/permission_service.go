package service

import (
	"errors"
	"supply-chain-server/internal/model"
	"supply-chain-server/internal/repository"
)

type PermissionService struct {
	repo *repository.PermissionRepository
}

func NewPermissionService(r *repository.PermissionRepository) *PermissionService {
	return &PermissionService{repo: r}
}

func (s *PermissionService) Create(permission *model.Permission) error {
	// 检查权限名是否已存在
	existing, _ := s.repo.GetByName(permission.Name)
	if existing != nil {
		return errors.New("权限名已存在")
	}
	return s.repo.Create(permission)
}

func (s *PermissionService) GetByID(id uint) (*model.Permission, error) {
	return s.repo.GetByID(id)
}

func (s *PermissionService) GetByName(name string) (*model.Permission, error) {
	return s.repo.GetByName(name)
}

func (s *PermissionService) List(page, pageSize int, keyword string) ([]model.Permission, int64, error) {
	return s.repo.List(page, pageSize, keyword)
}

func (s *PermissionService) GetAll() ([]model.Permission, error) {
	return s.repo.GetAll()
}

func (s *PermissionService) GetByRoleID(roleID uint) ([]model.Permission, error) {
	return s.repo.GetByRoleID(roleID)
}

func (s *PermissionService) Delete(id uint) error {
	return s.repo.Delete(id)
}

func (s *PermissionService) Update(permission *model.Permission) error {
	return s.repo.Update(permission)
}

// GetRepo 获取权限仓库实例
func (s *PermissionService) GetRepo() *repository.PermissionRepository {
	return s.repo
}