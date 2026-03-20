package service

import (
	"fmt"
	"math/rand"
	"supply-chain-server/internal/model"
	"supply-chain-server/internal/repository"
	"time"

	"github.com/go-playground/validator/v10"
)

type SupplierService struct {
	repo      *repository.SupplierRepository
	validator *validator.Validate
}

func NewSupplierService(r *repository.SupplierRepository) *SupplierService {
	return &SupplierService{
		repo:      r,
		validator: validator.New(),
	}
}

func (s *SupplierService) Create(supplier *model.Supplier) error {
	// 使用验证器验证字段
	if err := s.validator.Struct(supplier); err != nil {
		// 提取验证错误信息
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			for _, e := range validationErrors {
				switch e.Field() {
				case "Name":
					return fmt.Errorf("供应商名称不能为空")
				case "Contact":
					return fmt.Errorf("联系人不能为空")
				case "Phone":
					return fmt.Errorf("联系电话不能为空")
				case "Email":
					return fmt.Errorf("邮箱格式不正确")
				case "Level":
					return fmt.Errorf("供应商等级只能是 A、B 或 C")
				}
			}
		}
		return err
	}

	// 生成唯一编码（使用时间戳+随机数，确保并发安全）
	code, err := s.generateCode()
	if err != nil {
		return fmt.Errorf("生成编码失败: %w", err)
	}

	supplier.Code = code

	return s.repo.Create(supplier)
}

func (s *SupplierService) generateCode() (string, error) {
	// 使用纳秒时间戳+随机数生成唯一编码，确保并发安全
	date := time.Now().Format("20060102")
	rand.Seed(time.Now().UnixNano())
	// 使用纳秒时间戳的最后6位 + 3位随机数，几乎不会冲突
	nano := time.Now().Nanosecond() % 1000000
	random := rand.Intn(1000)
	return fmt.Sprintf("SUP%s%06d%03d", date, nano, random), nil
}

func (s *SupplierService) GetByID(id uint) (*model.Supplier, error) {
	return s.repo.GetByID(id)
}

func (s *SupplierService) Update(supplier *model.Supplier) error {
	// 使用验证器验证字段
	if err := s.validator.Struct(supplier); err != nil {
		// 提取验证错误信息
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			for _, e := range validationErrors {
				switch e.Field() {
				case "Name":
					return fmt.Errorf("供应商名称不能为空")
				case "Contact":
					return fmt.Errorf("联系人不能为空")
				case "Phone":
					return fmt.Errorf("联系电话不能为空")
				case "Email":
					return fmt.Errorf("邮箱格式不正确")
				case "Level":
					return fmt.Errorf("供应商等级只能是 A、B 或 C")
				}
			}
		}
		return err
	}

	return s.repo.Update(supplier)
}

func (s *SupplierService) Delete(id uint) error {
	return s.repo.Delete(id)
}

func (s *SupplierService) List(page, pageSize int, keyword string) ([]model.Supplier, int64, error) {
	return s.repo.List(page, pageSize, keyword)
}

func (s *SupplierService) GetActiveSuppliers() ([]model.Supplier, error) {
	return s.repo.GetActiveSuppliers()
}
