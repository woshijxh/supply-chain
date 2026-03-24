package repository

import (
	"supply-chain-server/internal/model"

	"gorm.io/gorm"
)

type TraceRepository struct {
	db *gorm.DB
}

func NewTraceRepository(db *gorm.DB) *TraceRepository {
	return &TraceRepository{db: db}
}

// FindProductByCode 通过编码查找产品(支持产品编码/SKU/批次号)
// 返回产品、编码类型、错误
func (r *TraceRepository) FindProductByCode(code string) (*model.Product, string, error) {
	var product model.Product

	// 1. 尝试作为产品编码查询
	if err := r.db.Where("code = ?", code).First(&product).Error; err == nil {
		return &product, "product_code", nil
	}

	// 2. 尝试作为SKU查询
	if err := r.db.Where("sku = ?", code).First(&product).Error; err == nil {
		return &product, "sku", nil
	}

	// 3. 尝试通过批次号查询
	var inventory model.Inventory
	if err := r.db.Where("batch_no = ?", code).First(&inventory).Error; err == nil {
		if err := r.db.First(&product, inventory.ProductID).Error; err == nil {
			return &product, "batch_no", nil
		}
	}

	return nil, "", gorm.ErrRecordNotFound
}

// GetProcurementTrace 获取采购追溯数据
func (r *TraceRepository) GetProcurementTrace(productID uint) ([]model.ProcurementTrace, error) {
	var traces []model.ProcurementTrace
	err := r.db.Table("procurement_orders po").
		Select(`po.id as order_id, po.order_no, po.supplier_id, s.name as supplier_name,
				pi.quantity, pi.received_qty, pi.unit_price, po.order_date,
				po.actual_date, po.warehouse, po.status`).
		Joins("JOIN procurement_items pi ON po.id = pi.order_id").
		Joins("JOIN suppliers s ON po.supplier_id = s.id").
		Where("pi.product_id = ?", productID).
		Order("po.created_at DESC").
		Scan(&traces).Error
	return traces, err
}

// GetInventoryTrace 获取库存流水追溯数据
func (r *TraceRepository) GetInventoryTrace(productID uint) ([]model.InventoryTrace, error) {
	var traces []model.InventoryTrace
	err := r.db.Table("inventory_logs").
		Where("product_id = ?", productID).
		Order("created_at DESC").
		Find(&traces).Error
	return traces, err
}

// GetSalesTrace 获取销售追溯数据
func (r *TraceRepository) GetSalesTrace(productID uint) ([]model.SalesTrace, error) {
	var traces []model.SalesTrace
	err := r.db.Table("sales_orders so").
		Select(`so.id as order_id, so.order_no, so.customer_name, so.customer_phone,
				si.quantity, si.unit_price, so.order_date, so.delivery_date,
				so.status, so.payment_status, lo.tracking_no, lo.status as logistics_status`).
		Joins("JOIN sales_order_items si ON so.id = si.order_id").
		Joins("LEFT JOIN logistics_orders lo ON so.id = lo.sales_order_id").
		Where("si.product_id = ?", productID).
		Order("so.created_at DESC").
		Scan(&traces).Error
	return traces, err
}

// GetLogisticsTrace 获取物流追溯数据
func (r *TraceRepository) GetLogisticsTrace(productID uint) ([]model.LogisticsTrace, error) {
	var traces []model.LogisticsTrace
	err := r.db.Table("logistics_orders lo").
		Select(`lo.id as logistics_id, lo.tracking_no, lo.carrier, lo.sales_order_no,
				lo.receiver_name, lo.receiver_phone, lo.status, lo.shipping_fee,
				lo.estimated_delivery, lo.actual_delivery`).
		Joins("JOIN sales_orders so ON lo.sales_order_id = so.id").
		Joins("JOIN sales_order_items si ON so.id = si.order_id").
		Where("si.product_id = ?", productID).
		Order("lo.created_at DESC").
		Scan(&traces).Error

	if err != nil {
		return nil, err
	}

	// 查询物流轨迹
	for i := range traces {
		var timeline []model.LogisticsTimelineItem
		r.db.Table("logistics_timelines").
			Select("time, status, location, description").
			Where("logistics_id = ?", traces[i].LogisticsID).
			Order("time DESC").
			Find(&timeline)
		traces[i].Timeline = timeline
	}

	return traces, nil
}

// GetReturnTrace 获取退货追溯数据
func (r *TraceRepository) GetReturnTrace(productID uint) ([]model.ReturnTrace, error) {
	var traces []model.ReturnTrace

	// 销售退货
	var salesReturns []model.ReturnTrace
	r.db.Table("sales_returns sr").
		Select(`sr.id as return_id, sr.return_no, 'sales_return' as type,
				sr.sales_order_no as order_no, sr.customer_name as related_name,
				sri.quantity, sri.unit_price, sr.reason, sr.status, sr.created_at`).
		Joins("JOIN sales_return_items sri ON sr.id = sri.return_id").
		Where("sri.product_id = ?", productID).
		Scan(&salesReturns)
	traces = append(traces, salesReturns...)

	// 采购退货
	var procReturns []model.ReturnTrace
	r.db.Table("procurement_returns pr").
		Select(`pr.id as return_id, pr.return_no, 'procurement_return' as type,
				pr.procurement_order_no as order_no, pr.supplier_name as related_name,
				pri.quantity, pri.unit_price, pr.reason, pr.status, pr.created_at`).
		Joins("JOIN procurement_return_items pri ON pr.id = pri.return_id").
		Where("pri.product_id = ?", productID).
		Scan(&procReturns)
	traces = append(traces, procReturns...)

	return traces, nil
}

// GetTotalStock 获取产品总库存
func (r *TraceRepository) GetTotalStock(productID uint) (int, error) {
	var total int
	err := r.db.Table("inventories").
		Where("product_id = ?", productID).
		Select("COALESCE(SUM(quantity), 0)").
		Scan(&total).Error
	return total, err
}