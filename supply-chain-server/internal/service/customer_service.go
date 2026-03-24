package service

import (
	"fmt"
	"supply-chain-server/internal/model"
	"supply-chain-server/internal/repository"
	"time"
)

type CustomerService struct {
	repo *repository.CustomerRepository
}

func NewCustomerService(r *repository.CustomerRepository) *CustomerService {
	return &CustomerService{repo: r}
}

func (s *CustomerService) Create(customer *model.Customer) error {
	if s.repo.ExistsByName(customer.Name) {
		return fmt.Errorf("客户名称已存在")
	}
	// 自动生成客户编码
	customer.Code = fmt.Sprintf("CU%s%04d", time.Now().Format("20060102"), time.Now().Nanosecond()%10000)
	return s.repo.Create(customer)
}

func (s *CustomerService) GetByID(id uint) (*model.Customer, error) {
	return s.repo.GetByID(id)
}

func (s *CustomerService) Update(customer *model.Customer) error {
	return s.repo.Update(customer)
}

func (s *CustomerService) Delete(id uint) error {
	return s.repo.Delete(id)
}

func (s *CustomerService) List(page, pageSize int, keyword string) ([]model.Customer, int64, error) {
	return s.repo.List(page, pageSize, keyword)
}

func (s *CustomerService) GetActiveCustomers() ([]model.Customer, error) {
	return s.repo.GetActiveCustomers()
}