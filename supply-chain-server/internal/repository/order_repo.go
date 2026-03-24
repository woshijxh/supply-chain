package repository

import (
	"supply-chain-server/internal/model"
	"time"

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

	query := r.db.Model(&model.ProcurementOrder{}).Preload("Supplier").Preload("Items")
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

// CountPending 统计待处理采购订单数量
func (r *ProcurementRepository) CountPending() (int64, error) {
	var count int64
	err := r.db.Model(&model.ProcurementOrder{}).
		Where("status IN ?", []string{"pending", "approved", "purchasing"}).
		Count(&count).Error
	return count, err
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

// GetTodayStats 获取今日销售统计
func (r *SalesRepository) GetTodayStats() (float64, int64, error) {
	var totalAmount float64
	var orderCount int64

	today := time.Now().Format("2006-01-02")
	err := r.db.Model(&model.SalesOrder{}).
		Where("DATE(order_date) = ?", today).
		Where("status != ?", "cancelled").
		Select("COALESCE(SUM(total_amount), 0)").
		Scan(&totalAmount).Error
	if err != nil {
		return 0, 0, err
	}

	err = r.db.Model(&model.SalesOrder{}).
		Where("DATE(order_date) = ?", today).
		Where("status != ?", "cancelled").
		Count(&orderCount).Error
	if err != nil {
		return 0, 0, err
	}

	return totalAmount, orderCount, nil
}

// GetYesterdayStats 获取昨日销售统计（用于计算增长）
func (r *SalesRepository) GetYesterdayStats() (float64, int64, error) {
	var totalAmount float64
	var orderCount int64

	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	err := r.db.Model(&model.SalesOrder{}).
		Where("DATE(order_date) = ?", yesterday).
		Where("status != ?", "cancelled").
		Select("COALESCE(SUM(total_amount), 0)").
		Scan(&totalAmount).Error
	if err != nil {
		return 0, 0, err
	}

	err = r.db.Model(&model.SalesOrder{}).
		Where("DATE(order_date) = ?", yesterday).
		Where("status != ?", "cancelled").
		Count(&orderCount).Error
	if err != nil {
		return 0, 0, err
	}

	return totalAmount, orderCount, nil
}

// CountByStatus 按状态统计销售订单数量
func (r *SalesRepository) CountByStatus(status string) (int64, error) {
	var count int64
	err := r.db.Model(&model.SalesOrder{}).Where("status = ?", status).Count(&count).Error
	return count, err
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

// UpdateStatusBySalesOrderID 根据销售订单ID更新物流状态
func (r *LogisticsRepository) UpdateStatusBySalesOrderID(salesOrderID uint, status string) error {
	return r.db.Model(&model.LogisticsOrder{}).
		Where("sales_order_id = ?", salesOrderID).
		Update("status", status).Error
}

// CountByStatus 按状态统计物流订单数量
func (r *LogisticsRepository) CountByStatus(status string) (int64, error) {
	var count int64
	err := r.db.Model(&model.LogisticsOrder{}).Where("status = ?", status).Count(&count).Error
	return count, err
}
