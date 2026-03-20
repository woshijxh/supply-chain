package repository

import (
	"supply-chain-server/internal/model"

	"gorm.io/gorm"
)

type InventoryRepository struct {
	db *gorm.DB
}

func NewInventoryRepository(db *gorm.DB) *InventoryRepository {
	return &InventoryRepository{db: db}
}

func (r *InventoryRepository) Create(item *model.Inventory) error {
	return r.db.Create(item).Error
}

func (r *InventoryRepository) GetByID(id uint) (*model.Inventory, error) {
	var item model.Inventory
	err := r.db.Preload("Product").First(&item, id).Error
	return &item, err
}

func (r *InventoryRepository) Update(item *model.Inventory) error {
	return r.db.Save(item).Error
}

func (r *InventoryRepository) List(page, pageSize int, status, warehouse string) ([]model.Inventory, int64, error) {
	var items []model.Inventory
	var total int64

	query := r.db.Model(&model.Inventory{}).Preload("Product")
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if warehouse != "" {
		query = query.Where("warehouse = ?", warehouse)
	}

	query.Count(&total)
	err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&items).Error
	return items, total, err
}

func (r *InventoryRepository) GetByProductIDForUpdate(tx *gorm.DB, productID uint, warehouse string) (*model.Inventory, error) {
	var item model.Inventory
	query := tx.Where("product_id = ?", productID)
	if warehouse != "" {
		query = query.Where("warehouse = ?", warehouse)
	}
	err := query.First(&item).Error
	return &item, err
}

func (r *InventoryRepository) CreateInTransaction(tx *gorm.DB, item *model.Inventory) error {
	return tx.Create(item).Error
}

func (r *InventoryRepository) UpdateInTransaction(tx *gorm.DB, item *model.Inventory) error {
	return tx.Save(item).Error
}

func (r *InventoryRepository) CountByStatus(status string) (int64, error) {
	var count int64
	err := r.db.Model(&model.Inventory{}).Where("status = ?", status).Count(&count).Error
	return count, err
}

func (r *InventoryRepository) CountAll() (int64, error) {
	var count int64
	err := r.db.Model(&model.Inventory{}).Count(&count).Error
	return count, err
}
