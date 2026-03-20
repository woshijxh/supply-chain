package repository

import (
	"supply-chain-server/internal/model"

	"gorm.io/gorm"
)

type ProcurementRepository struct {
	db *gorm.DB
}

func NewProcurementRepository(db *gorm.DB) *ProcurementRepository {
	return &ProcurementRepository{db: db}
}

func (r *ProcurementRepository) Create(order *model.ProcurementOrder) error {
	return r.db.Create(order).Error
}

func (r *ProcurementRepository) GetByID(id uint) (*model.ProcurementOrder, error) {
	var order model.ProcurementOrder
	err := r.db.Preload("Supplier").Preload("Items").First(&order, id).Error
	return &order, err
}

func (r *ProcurementRepository) Update(order *model.ProcurementOrder) error {
	return r.db.Save(order).Error
}

func (r *ProcurementRepository) Delete(id uint) error {
	return r.db.Delete(&model.ProcurementOrder{}, id).Error
}

func (r *ProcurementRepository) List(page, pageSize int, status string) ([]model.ProcurementOrder, int64, error) {
	var orders []model.ProcurementOrder
	var total int64

	query := r.db.Model(&model.ProcurementOrder{}).Preload("Supplier")
	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)
	err := query.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&orders).Error
	return orders, total, err
}

func (r *ProcurementRepository) CreateItem(item *model.ProcurementItem) error {
	return r.db.Create(item).Error
}

func (r *ProcurementRepository) UpdateItem(item *model.ProcurementItem) error {
	return r.db.Save(item).Error
}

type SalesRepository struct {
	db *gorm.DB
}

func NewSalesRepository(db *gorm.DB) *SalesRepository {
	return &SalesRepository{db: db}
}

func (r *SalesRepository) Create(order *model.SalesOrder) error {
	return r.db.Create(order).Error
}

func (r *SalesRepository) GetByID(id uint) (*model.SalesOrder, error) {
	var order model.SalesOrder
	err := r.db.Preload("Items").First(&order, id).Error
	return &order, err
}

func (r *SalesRepository) Update(order *model.SalesOrder) error {
	return r.db.Save(order).Error
}

func (r *SalesRepository) Delete(id uint) error {
	return r.db.Delete(&model.SalesOrder{}, id).Error
}

func (r *SalesRepository) List(page, pageSize int, status string) ([]model.SalesOrder, int64, error) {
	var orders []model.SalesOrder
	var total int64

	query := r.db.Model(&model.SalesOrder{})
	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)
	err := query.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&orders).Error
	return orders, total, err
}

func (r *SalesRepository) CreateItem(item *model.SalesOrderItem) error {
	return r.db.Create(item).Error
}

type LogisticsRepository struct {
	db *gorm.DB
}

func NewLogisticsRepository(db *gorm.DB) *LogisticsRepository {
	return &LogisticsRepository{db: db}
}

func (r *LogisticsRepository) Create(order *model.LogisticsOrder) error {
	return r.db.Create(order).Error
}

func (r *LogisticsRepository) GetByID(id uint) (*model.LogisticsOrder, error) {
	var order model.LogisticsOrder
	err := r.db.Preload("Timeline").First(&order, id).Error
	return &order, err
}

func (r *LogisticsRepository) Update(order *model.LogisticsOrder) error {
	return r.db.Save(order).Error
}

func (r *LogisticsRepository) Delete(id uint) error {
	return r.db.Delete(&model.LogisticsOrder{}, id).Error
}

func (r *LogisticsRepository) List(page, pageSize int, status string) ([]model.LogisticsOrder, int64, error) {
	var orders []model.LogisticsOrder
	var total int64

	query := r.db.Model(&model.LogisticsOrder{})
	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)
	err := query.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&orders).Error
	return orders, total, err
}

func (r *LogisticsRepository) CreateTimeline(item *model.LogisticsTimeline) error {
	return r.db.Create(item).Error
}
