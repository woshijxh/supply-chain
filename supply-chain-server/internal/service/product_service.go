package service

import (
	"supply-chain-server/internal/model"
	"supply-chain-server/internal/repository"
)

type ProductService struct {
	repo *repository.ProductRepository
}

func NewProductService(r *repository.ProductRepository) *ProductService {
	return &ProductService{repo: r}
}

func (s *ProductService) List(page, pageSize int, keyword string) ([]model.Product, int64, error) {
	return s.repo.List(page, pageSize, keyword)
}

func (s *ProductService) GetByID(id uint) (*model.Product, error) {
	return s.repo.GetByID(id)
}

func (s *ProductService) Create(product *model.Product) error {
	return s.repo.Create(product)
}
