package service

import (
	"errors"
	"supply-chain-server/internal/model"
	"supply-chain-server/internal/repository"
)

type RoleService struct {
	repo *repository.RoleRepository
}

func NewRoleService(r *repository.RoleRepository) *RoleService {
	return &RoleService{repo: r}
}

func (s *RoleService) Create(role *model.Role) error {
	// 检查角色名是否已存在
	existing, _ := s.repo.GetByName(role.Name)
	if existing != nil {
		return errors.New("角色名已存在")
	}
	return s.repo.Create(role)
}

func (s *RoleService) GetByID(id uint) (*model.Role, error) {
	return s.repo.GetByID(id)
}

func (s *RoleService) GetByName(name string) (*model.Role, error) {
	return s.repo.GetByName(name)
}

func (s *RoleService) List(page, pageSize int, keyword string) ([]model.Role, int64, error) {
	return s.repo.List(page, pageSize, keyword)
}

func (s *RoleService) Update(role *model.Role) error {
	return s.repo.Update(role)
}

func (s *RoleService) Delete(id uint) error {
	return s.repo.Delete(id)
}

func (s *RoleService) GetAll() ([]model.Role, error) {
	return s.repo.GetAll()
}

func (s *RoleService) SetPermissions(roleID uint, permissionIDs []uint, permRepo *repository.PermissionRepository) error {
	return permRepo.SetRolePermissions(roleID, permissionIDs)
}