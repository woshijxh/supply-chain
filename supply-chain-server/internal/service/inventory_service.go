package service

import (
	"errors"
	"supply-chain-server/internal/model"
	"supply-chain-server/internal/repository"
	"supply-chain-server/pkg/database"
	"time"

	"gorm.io/gorm"
)

type InventoryService struct {
	repo        *repository.InventoryRepository
	procurementRepo *repository.ProcurementRepository
}

func NewInventoryService(r *repository.InventoryRepository, p *repository.ProcurementRepository) *InventoryService {
	return &InventoryService{repo: r, procurementRepo: p}
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
		item, err := s.repo.GetByProductIDForUpdate(tx, productID, warehouse)
		if err != nil {
			item = &model.Inventory{
				ProductID:    productID,
				Warehouse:    warehouse,
				Quantity:     int(quantity),
				AvailableQty: int(quantity),
				LockedQty:    0,
				Status:       "normal",
			}
			return s.repo.CreateInTransaction(tx, item)
		}

		item.Quantity += int(quantity)
		item.AvailableQty += int(quantity)
		s.updateStatus(item)

		return s.repo.UpdateInTransaction(tx, item)
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
			now := model.Date{time.Now()}
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
		item, err := s.repo.GetByProductIDForUpdate(tx, productID, warehouse)
		if err != nil {
			return errors.New("库存不存在")
		}

		if item.AvailableQty < int(quantity) {
			return errors.New("库存不足")
		}

		item.Quantity -= int(quantity)
		item.AvailableQty -= int(quantity)
		s.updateStatus(item)

		return s.repo.UpdateInTransaction(tx, item)
	})
}

func (s *InventoryService) updateStatus(item *model.Inventory) {
	if item.Quantity <= 0 {
		item.Status = "low"
	} else if item.Product.MinStock > 0 && item.Quantity < item.Product.MinStock {
		item.Status = "low"
	} else if item.Product.MaxStock > 0 && item.Quantity > item.Product.MaxStock {
		item.Status = "over"
	} else {
		item.Status = "normal"
	}
}

func (s *InventoryService) GetStats() (map[string]interface{}, error) {
	total, err := s.repo.CountAll()
	if err != nil {
		return nil, err
	}

	normal, _ := s.repo.CountByStatus("normal")
	low, _ := s.repo.CountByStatus("low")
	over, _ := s.repo.CountByStatus("over")

	return map[string]interface{}{
		"total":  total,
		"normal": normal,
		"low":    low,
		"over":   over,
	}, nil
}
