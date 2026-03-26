package service

import (
	"errors"
	"fmt"
	"supply-chain-server/internal/model"
	"supply-chain-server/internal/repository"
	"supply-chain-server/pkg/database"
	"time"

	"gorm.io/gorm"
)

type InventoryService struct {
	repo            *repository.InventoryRepository
	procurementRepo *repository.ProcurementRepository
	logRepo         *repository.InventoryLogRepository
	productRepo     *repository.ProductRepository
}

func NewInventoryService(r *repository.InventoryRepository, p *repository.ProcurementRepository, lr *repository.InventoryLogRepository, pr *repository.ProductRepository) *InventoryService {
	return &InventoryService{repo: r, procurementRepo: p, logRepo: lr, productRepo: pr}
}

// SetProductRepository 设置产品仓库（用于获取产品名称）- 保留用于向后兼容
func (s *InventoryService) SetProductRepo(pr *repository.ProductRepository) {
	s.productRepo = pr
}

func (s *InventoryService) Create(item *model.Inventory) error {
	return s.repo.Create(item)
}

func (s *InventoryService) GetByID(id uint) (*model.Inventory, error) {
	return s.repo.GetByID(id)
}

func (s *InventoryService) Update(item *model.Inventory) error {
	return s.repo.Update(item)
}

func (s *InventoryService) List(page, pageSize int, status, warehouse string) ([]model.Inventory, int64, error) {
	return s.repo.List(page, pageSize, status, warehouse)
}

func (s *InventoryService) StockIn(productID, quantity uint, warehouse string) error {
	if warehouse == "" {
		warehouse = "默认仓库"
	}

	return database.DB.Transaction(func(tx *gorm.DB) error {
		// 获取产品信息和名称
		var productName string
		var minStock, maxStock int
		if s.productRepo != nil {
			if product, err := s.productRepo.GetByID(productID); err == nil {
				productName = product.Name
				minStock = product.MinStock
				maxStock = product.MaxStock
			}
		}

		item, err := s.repo.GetByProductIDForUpdate(tx, productID, warehouse)
		var beforeQty int
		if err != nil {
			beforeQty = 0
			item = &model.Inventory{
				ProductID:    productID,
				Warehouse:    warehouse,
				Quantity:     int(quantity),
				AvailableQty: int(quantity),
				LockedQty:    0,
				Status:       "normal",
			}
			// 初始化 Product 用于状态计算
			item.Product = model.Product{MinStock: minStock, MaxStock: maxStock}
			s.updateStatus(item)
			if err := s.repo.CreateInTransaction(tx, item); err != nil {
				return err
			}
		} else {
			beforeQty = item.Quantity
			item.Quantity += int(quantity)
			item.AvailableQty += int(quantity)
			s.updateStatus(item)

			if err := s.repo.UpdateInTransaction(tx, item); err != nil {
				return err
			}
		}

		// 记录库存流水
		s.createLog(productID, productName, "in", int(quantity), beforeQty, item.Quantity, warehouse, "", "", "", "手动入库")

		return nil
	})
}

func (s *InventoryService) StockInWithProcurement(productID, quantity uint, warehouse string, procurementID, procurementItemID uint) error {
	// 如果有关联采购订单，验证并获取信息
	if procurementID > 0 {
		order, err := s.procurementRepo.GetByID(procurementID)
		if err != nil {
			return errors.New("采购订单不存在")
		}
		// 只允许 pending/approved/purchasing 状态的采购订单可以入库
		if order.Status == "received" || order.Status == "cancelled" {
			return errors.New("采购订单已收货或已取消")
		}
		// 使用采购订单的仓库
		if warehouse == "" {
			warehouse = order.Warehouse
		}
		// 如果指定了采购明细，验证产品匹配
		if procurementItemID > 0 {
			var found bool
			for _, item := range order.Items {
				if item.ID == procurementItemID {
					found = true
					if item.ProductID != productID {
						return errors.New("产品与采购明细不匹配")
					}
					// 检查入库数量是否超过待收货数量
					pendingQty := item.Quantity - item.ReceivedQty
					if int(quantity) > pendingQty {
						return errors.New("入库数量超过待收货数量")
					}
					break
				}
			}
			if !found {
				return errors.New("采购明细不存在")
			}
		}
	}

	// 执行入库
	if err := s.StockIn(productID, quantity, warehouse); err != nil {
		return err
	}

	// 更新采购明细的已收货数量
	if procurementID > 0 && procurementItemID > 0 {
		order, _ := s.procurementRepo.GetByID(procurementID)
		for _, item := range order.Items {
			if item.ID == procurementItemID {
				item.ReceivedQty += int(quantity)
				s.procurementRepo.UpdateItem(&item)
				break
			}
		}
		// 检查是否全部收货完毕
		allReceived := true
		for _, item := range order.Items {
			if item.ReceivedQty < item.Quantity {
				allReceived = false
				break
			}
		}
		// 如果全部收货，自动更新状态为 received
		if allReceived && order.Status != "received" {
			order.Status = "received"
			now := model.Date{Time: time.Now()}
			order.ActualDate = &now
			s.procurementRepo.Update(order)
		}
	}

	return nil
}

func (s *InventoryService) StockOut(productID, quantity uint, warehouse string) error {
	if warehouse == "" {
		warehouse = "默认仓库"
	}

	return database.DB.Transaction(func(tx *gorm.DB) error {
		// 获取产品名称
		var productName string
		if s.productRepo != nil {
			if product, err := s.productRepo.GetByID(productID); err == nil {
				productName = product.Name
			}
		}

		item, err := s.repo.GetByProductIDForUpdate(tx, productID, warehouse)
		if err != nil {
			return errors.New("库存不存在")
		}

		if item.AvailableQty < int(quantity) {
			return errors.New("库存不足")
		}

		beforeQty := item.Quantity
		item.Quantity -= int(quantity)
		item.AvailableQty -= int(quantity)
		s.updateStatus(item)

		if err := s.repo.UpdateInTransaction(tx, item); err != nil {
			return err
		}

		// 记录库存流水
		s.createLog(productID, productName, "out", int(quantity), beforeQty, item.Quantity, warehouse, "", "", "", "手动出库")

		return nil
	})
}

func (s *InventoryService) updateStatus(item *model.Inventory) {
	if item.Quantity <= 0 {
		item.Status = "low"
	} else if item.Product.ID > 0 && item.Product.MinStock > 0 && item.Quantity < item.Product.MinStock {
		item.Status = "low"
	} else if item.Product.ID > 0 && item.Product.MaxStock > 0 && item.Quantity > item.Product.MaxStock {
		item.Status = "over"
	} else {
		item.Status = "normal"
	}
}

func (s *InventoryService) GetStats() (map[string]interface{}, error) {
	// 统计库存总量（所有库存数量之和）
	totalQty, _ := s.repo.SumTotalQuantity()

	// 统计各状态的 SKU 数
	normal, _ := s.repo.CountByStatus("normal")
	low, _ := s.repo.CountByStatus("low")
	over, _ := s.repo.CountByStatus("over")

	// 统计总 SKU 数（有库存记录的产品数）
	totalSku, _ := s.repo.CountDistinctProducts()

	return map[string]interface{}{
		"total":    totalQty,  // 库存总量
		"normal":   normal,
		"low":      low,
		"over":     over,
		"totalSku": totalSku,
	}, nil
}

// LockStock 锁定库存（创建销售订单时调用）
func (s *InventoryService) LockStock(productID, quantity uint, warehouse string) error {
	if warehouse == "" {
		warehouse = "默认仓库"
	}

	return database.DB.Transaction(func(tx *gorm.DB) error {
		// 获取产品名称
		var productName string
		if s.productRepo != nil {
			if product, err := s.productRepo.GetByID(productID); err == nil {
				productName = product.Name
			}
		}

		item, err := s.repo.GetByProductIDForUpdate(tx, productID, warehouse)
		if err != nil {
			return errors.New("库存不存在")
		}

		if item.AvailableQty < int(quantity) {
			return fmt.Errorf("库存不足，可用库存: %d，需要: %d", item.AvailableQty, quantity)
		}

		beforeQty := item.Quantity
		item.LockedQty += int(quantity)
		item.AvailableQty -= int(quantity)
		s.updateStatus(item)

		if err := s.repo.UpdateInTransaction(tx, item); err != nil {
			return err
		}

		// 记录库存流水
		s.createLog(productID, productName, "lock", int(quantity), beforeQty, item.Quantity, warehouse, "", "", "", "销售锁定")

		return nil
	})
}

// UnlockStock 解锁库存（取消销售订单时调用）
func (s *InventoryService) UnlockStock(productID, quantity uint, warehouse string) error {
	if warehouse == "" {
		warehouse = "默认仓库"
	}

	return database.DB.Transaction(func(tx *gorm.DB) error {
		// 获取产品名称
		var productName string
		if s.productRepo != nil {
			if product, err := s.productRepo.GetByID(productID); err == nil {
				productName = product.Name
			}
		}

		item, err := s.repo.GetByProductIDForUpdate(tx, productID, warehouse)
		if err != nil {
			return errors.New("库存不存在")
		}

		if item.LockedQty < int(quantity) {
			quantity = uint(item.LockedQty) // 防止解锁超过已锁定数量
		}

		beforeQty := item.Quantity
		item.LockedQty -= int(quantity)
		item.AvailableQty += int(quantity)
		s.updateStatus(item)

		if err := s.repo.UpdateInTransaction(tx, item); err != nil {
			return err
		}

		// 记录库存流水
		s.createLog(productID, productName, "unlock", int(quantity), beforeQty, item.Quantity, warehouse, "", "", "", "取消锁定")

		return nil
	})
}

// DeductLockedStock 扣减已锁定的库存（销售订单确认出库时调用）
func (s *InventoryService) DeductLockedStock(productID, quantity uint, warehouse string) error {
	if warehouse == "" {
		warehouse = "默认仓库"
	}

	return database.DB.Transaction(func(tx *gorm.DB) error {
		// 获取产品名称
		var productName string
		if s.productRepo != nil {
			if product, err := s.productRepo.GetByID(productID); err == nil {
				productName = product.Name
			}
		}

		item, err := s.repo.GetByProductIDForUpdate(tx, productID, warehouse)
		if err != nil {
			return errors.New("库存不存在")
		}

		if item.LockedQty < int(quantity) {
			return fmt.Errorf("锁定库存不足，已锁定: %d，需要扣减: %d", item.LockedQty, quantity)
		}

		beforeQty := item.Quantity
		item.LockedQty -= int(quantity)
		item.Quantity -= int(quantity)
		s.updateStatus(item)

		if err := s.repo.UpdateInTransaction(tx, item); err != nil {
			return err
		}

		// 记录库存流水
		s.createLog(productID, productName, "out", int(quantity), beforeQty, item.Quantity, warehouse, "", "", "", "销售出库")

		return nil
	})
}

// CheckStock 检查库存是否充足
func (s *InventoryService) CheckStock(productID, quantity uint, warehouse string) (bool, int, error) {
	if warehouse == "" {
		warehouse = "默认仓库"
	}

	item, err := s.repo.GetByProductIDForUpdate(database.DB, productID, warehouse)
	if err != nil {
		return false, 0, errors.New("库存不存在")
	}

	return item.AvailableQty >= int(quantity), item.AvailableQty, nil
}

// createLog 创建库存流水记录
func (s *InventoryService) createLog(productID uint, productName string, logType string, quantity, beforeQty, afterQty int, warehouse, refType, refNo, operator, remark string) {
	if s.logRepo == nil {
		return
	}
	log := &model.InventoryLog{
		ProductID:   productID,
		ProductName: productName,
		Type:        logType,
		Quantity:    quantity,
		BeforeQty:   beforeQty,
		AfterQty:    afterQty,
		Warehouse:   warehouse,
		RefType:     refType,
		RefNo:       refNo,
		Operator:    operator,
		Remark:      remark,
		CreatedAt:   time.Now(),
	}
	s.logRepo.Create(log)
}

// GetLogs 获取库存流水记录
func (s *InventoryService) GetLogs(page, pageSize int, productID uint, logType string) ([]model.InventoryLog, int64, error) {
	if s.logRepo == nil {
		return []model.InventoryLog{}, 0, nil
	}
	return s.logRepo.List(page, pageSize, productID, logType)
}

// GetLowStockItems 获取低库存产品列表
func (s *InventoryService) GetLowStockItems(limit int) ([]model.Inventory, error) {
	return s.repo.GetLowStockItems(limit)
}

// GetInventoryDistribution 获取库存分布（按产品分类）
func (s *InventoryService) GetInventoryDistribution() ([]map[string]interface{}, error) {
	return s.repo.GetInventoryDistribution()
}
