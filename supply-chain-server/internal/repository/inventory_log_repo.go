package repository

import (
	"supply-chain-server/internal/model"

	"gorm.io/gorm"
)

type InventoryLogRepository struct {
	db *gorm.DB
}

func NewInventoryLogRepository(db *gorm.DB) *InventoryLogRepository {
	return &InventoryLogRepository{db: db}
}

func (r *InventoryLogRepository) Create(log *model.InventoryLog) error {
	return r.db.Create(log).Error
}

func (r *InventoryLogRepository) List(page, pageSize int, productID uint, logType string) ([]model.InventoryLog, int64, error) {
	var logs []model.InventoryLog
	var total int64

	query := r.db.Model(&model.InventoryLog{})
	if productID > 0 {
		query = query.Where("product_id = ?", productID)
	}
	if logType != "" {
		query = query.Where("type = ?", logType)
	}

	query.Count(&total)
	err := query.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&logs).Error
	return logs, total, err
}

func (r *InventoryLogRepository) GetByProductID(productID uint, limit int) ([]model.InventoryLog, error) {
	var logs []model.InventoryLog
	err := r.db.Where("product_id = ?", productID).Order("created_at DESC").Limit(limit).Find(&logs).Error
	return logs, err
}