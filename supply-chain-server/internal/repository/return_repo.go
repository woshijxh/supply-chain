package repository

import (
	"supply-chain-server/internal/model"

	"gorm.io/gorm"
)

type SalesReturnRepository struct {
	db *gorm.DB
}

func NewSalesReturnRepository(db *gorm.DB) *SalesReturnRepository {
	return &SalesReturnRepository{db: db}
}

func (r *SalesReturnRepository) Create(ret *model.SalesReturn) error {
	return r.db.Create(ret).Error
}

func (r *SalesReturnRepository) GetByID(id uint) (*model.SalesReturn, error) {
	var ret model.SalesReturn
	err := r.db.Preload("Items").First(&ret, id).Error
	return &ret, err
}

func (r *SalesReturnRepository) Update(ret *model.SalesReturn) error {
	return r.db.Save(ret).Error
}

func (r *SalesReturnRepository) Delete(id uint) error {
	return r.db.Delete(&model.SalesReturn{}, id).Error
}

func (r *SalesReturnRepository) List(page, pageSize int, status string) ([]model.SalesReturn, int64, error) {
	var returns []model.SalesReturn
	var total int64

	query := r.db.Model(&model.SalesReturn{})
	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)
	err := query.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&returns).Error
	return returns, total, err
}

// CountByStatus 按状态统计销售退货数量
func (r *SalesReturnRepository) CountByStatus(status string) (int64, error) {
	var count int64
	err := r.db.Model(&model.SalesReturn{}).Where("status = ?", status).Count(&count).Error
	return count, err
}

type ProcurementReturnRepository struct {
	db *gorm.DB
}

func NewProcurementReturnRepository(db *gorm.DB) *ProcurementReturnRepository {
	return &ProcurementReturnRepository{db: db}
}

func (r *ProcurementReturnRepository) Create(ret *model.ProcurementReturn) error {
	return r.db.Create(ret).Error
}

func (r *ProcurementReturnRepository) GetByID(id uint) (*model.ProcurementReturn, error) {
	var ret model.ProcurementReturn
	err := r.db.Preload("Items").First(&ret, id).Error
	return &ret, err
}

func (r *ProcurementReturnRepository) Update(ret *model.ProcurementReturn) error {
	return r.db.Save(ret).Error
}

func (r *ProcurementReturnRepository) Delete(id uint) error {
	return r.db.Delete(&model.ProcurementReturn{}, id).Error
}

func (r *ProcurementReturnRepository) List(page, pageSize int, status string) ([]model.ProcurementReturn, int64, error) {
	var returns []model.ProcurementReturn
	var total int64

	query := r.db.Model(&model.ProcurementReturn{})
	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)
	err := query.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&returns).Error
	return returns, total, err
}