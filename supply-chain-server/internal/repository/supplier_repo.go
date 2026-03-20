package repository

import (
	"supply-chain-server/internal/model"

	"gorm.io/gorm"
)

type SupplierRepository struct {
	db *gorm.DB
}

func NewSupplierRepository(db *gorm.DB) *SupplierRepository {
	return &SupplierRepository{db: db}
}

func (r *SupplierRepository) Create(supplier *model.Supplier) error {
	// 编码由Service层生成，这里不再重复生成
	return r.db.Create(supplier).Error
}

func (r *SupplierRepository) GetByID(id uint) (*model.Supplier, error) {
	var supplier model.Supplier
	err := r.db.First(&supplier, id).Error
	return &supplier, err
}

func (r *SupplierRepository) Update(supplier *model.Supplier) error {
	// 更新时忽略 code 字段，编码不允许修改
	return r.db.Model(&model.Supplier{}).Where("id = ?", supplier.ID).Omit("code").Updates(supplier).Error
}

func (r *SupplierRepository) Delete(id uint) error {
	return r.db.Delete(&model.Supplier{}, id).Error
}

func (r *SupplierRepository) List(page, pageSize int, keyword string) ([]model.Supplier, int64, error) {
	var suppliers []model.Supplier
	var total int64

	query := r.db.Model(&model.Supplier{})
	if keyword != "" {
		query = query.Where("name LIKE ? OR code LIKE ? OR contact LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	query.Count(&total)
	err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&suppliers).Error
	return suppliers, total, err
}

func (r *SupplierRepository) GetActiveSuppliers() ([]model.Supplier, error) {
	var suppliers []model.Supplier
	err := r.db.Where("status = ?", 1).Find(&suppliers).Error
	return suppliers, err
}
