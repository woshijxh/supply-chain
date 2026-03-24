package service

import (
	"fmt"
	"supply-chain-server/internal/model"
	"supply-chain-server/internal/repository"
	"time"
)

type ProcurementService struct {
	repo            *repository.ProcurementRepository
	inventoryService *InventoryService
}

func NewProcurementService(r *repository.ProcurementRepository, inv *InventoryService) *ProcurementService {
	return &ProcurementService{repo: r, inventoryService: inv}
}

func (s *ProcurementService) Create(order *model.ProcurementOrder) error {
	order.OrderDate = model.Date{Time: time.Now()}
	// 自动生成采购订单号: PO + 年月日 + 6位随机数
	order.OrderNo = fmt.Sprintf("PO%s%06d", time.Now().Format("20060102"), time.Now().Nanosecond()%1000000)
	return s.repo.Create(order)
}

func (s *ProcurementService) GetByID(id uint) (*model.ProcurementOrder, error) {
	return s.repo.GetByID(id)
}

func (s *ProcurementService) Update(order *model.ProcurementOrder) error {
	return s.repo.Update(order)
}

func (s *ProcurementService) Delete(id uint) error {
	return s.repo.Delete(id)
}

func (s *ProcurementService) List(page, pageSize int, status string) ([]model.ProcurementOrder, int64, error) {
	return s.repo.List(page, pageSize, status)
}

// CountPending 统计待处理采购订单数量
func (s *ProcurementService) CountPending() (int64, error) {
	return s.repo.CountPending()
}

func (s *ProcurementService) UpdateStatus(id uint, status, attachmentURL, remark string) error {
	order, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	order.Status = status
	if attachmentURL != "" {
		order.AttachmentURL = attachmentURL
	}
	if remark != "" {
		order.Remark = remark
	}
	if status == "received" {
		now := model.Date{Time: time.Now()}
		order.ActualDate = &now
		// 自动入库：遍历采购明细，为每个产品入库
		for _, item := range order.Items {
			if item.ProductID == 0 {
				continue
			}
			// 计算待入库数量
			pendingQty := item.Quantity - item.ReceivedQty
			if pendingQty <= 0 {
				continue
			}
			// 调用库存服务入库
			err := s.inventoryService.StockIn(item.ProductID, uint(pendingQty), order.Warehouse)
			if err != nil {
				return fmt.Errorf("入库失败: %w", err)
			}
			// 更新已收货数量
			item.ReceivedQty = item.Quantity
			s.repo.UpdateItem(&item)
		}
	}
	return s.repo.Update(order)
}

type SalesService struct {
	repo              *repository.SalesRepository
	inventoryService  *InventoryService
	logisticsService  *LogisticsService
}

func NewSalesService(r *repository.SalesRepository, inv *InventoryService, log *LogisticsService) *SalesService {
	return &SalesService{repo: r, inventoryService: inv, logisticsService: log}
}

func (s *SalesService) Create(order *model.SalesOrder) error {
	order.OrderDate = time.Now()
	// 自动生成销售订单号: SO + 年月日 + 6位随机数
	order.OrderNo = fmt.Sprintf("SO%s%06d", time.Now().Format("20060102"), time.Now().Nanosecond()%1000000)

	// 检查并锁定库存
	if s.inventoryService != nil {
		for _, item := range order.Items {
			if item.ProductID == 0 || item.Quantity <= 0 {
				continue
			}
			// 检查库存是否充足
			ok, available, err := s.inventoryService.CheckStock(uint(item.ProductID), uint(item.Quantity), "")
			if err != nil {
				return fmt.Errorf("产品ID %d 库存检查失败: %w", item.ProductID, err)
			}
			if !ok {
				return fmt.Errorf("产品ID %d 库存不足，可用库存: %d，需要: %d", item.ProductID, available, item.Quantity)
			}
		}

		// 锁定库存
		for _, item := range order.Items {
			if item.ProductID == 0 || item.Quantity <= 0 {
				continue
			}
			if err := s.inventoryService.LockStock(uint(item.ProductID), uint(item.Quantity), ""); err != nil {
				// 如果锁定失败，需要回滚之前锁定的库存
				for i := range order.Items {
					if i >= len(order.Items) || order.Items[i].ProductID == item.ProductID {
						break
					}
					if order.Items[i].ProductID > 0 && order.Items[i].Quantity > 0 {
						s.inventoryService.UnlockStock(uint(order.Items[i].ProductID), uint(order.Items[i].Quantity), "")
					}
				}
				return fmt.Errorf("锁定库存失败: %w", err)
			}
		}
	}

	// 设置初始状态
	order.Status = "pending"
	order.PaymentStatus = "pending"

	return s.repo.Create(order)
}

func (s *SalesService) GetByID(id uint) (*model.SalesOrder, error) {
	return s.repo.GetByID(id)
}

func (s *SalesService) Update(order *model.SalesOrder) error {
	return s.repo.Update(order)
}

func (s *SalesService) Delete(id uint) error {
	// 删除前先处理库存
	if s.inventoryService != nil {
		order, err := s.repo.GetByID(id)
		if err != nil {
			return err
		}
		// 根据订单状态处理库存
		switch order.Status {
		case "pending":
			// 待确认状态：解锁库存
			for _, item := range order.Items {
				if item.ProductID == 0 || item.Quantity <= 0 {
					continue
				}
				s.inventoryService.UnlockStock(uint(item.ProductID), uint(item.Quantity), "")
			}
		case "confirmed", "shipping":
			// 已确认/已发货状态：退货入库
			for _, item := range order.Items {
				if item.ProductID == 0 || item.Quantity <= 0 {
					continue
				}
				s.inventoryService.StockIn(uint(item.ProductID), uint(item.Quantity), "")
			}
		}
	}
	return s.repo.Delete(id)
}

func (s *SalesService) List(page, pageSize int, status string) ([]model.SalesOrder, int64, error) {
	return s.repo.List(page, pageSize, status)
}

// GetTodayStats 获取今日销售统计
func (s *SalesService) GetTodayStats() (float64, int64, error) {
	return s.repo.GetTodayStats()
}

// GetYesterdayStats 获取昨日销售统计
func (s *SalesService) GetYesterdayStats() (float64, int64, error) {
	return s.repo.GetYesterdayStats()
}

// CountByStatus 按状态统计订单数量
func (s *SalesService) CountByStatus(status string) (int64, error) {
	return s.repo.CountByStatus(status)
}

func (s *SalesService) UpdateStatus(id uint, status string) error {
	order, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	// 状态流转检查
	validTransitions := map[string][]string{
		"pending":   {"confirmed", "cancelled"},
		"confirmed": {"shipping", "cancelled"},
		"shipping":  {"completed", "cancelled"},
		"completed": {},
		"cancelled": {},
	}

	// 检查状态转换是否合法
	allowed, ok := validTransitions[order.Status]
	if !ok {
		return fmt.Errorf("当前状态 %s 不允许变更", order.Status)
	}

	isValidTransition := false
	for _, s := range allowed {
		if s == status {
			isValidTransition = true
			break
		}
	}
	if !isValidTransition {
		return fmt.Errorf("状态不允许从 %s 变更为 %s", order.Status, status)
	}

	// 处理库存和物流
	switch status {
	case "confirmed":
		// 确认订单：扣减锁定库存（出库）+ 自动创建物流订单
		if s.inventoryService != nil {
			for _, item := range order.Items {
				if item.ProductID == 0 || item.Quantity <= 0 {
					continue
				}
				if err := s.inventoryService.DeductLockedStock(uint(item.ProductID), uint(item.Quantity), ""); err != nil {
					return fmt.Errorf("扣减库存失败: %w", err)
				}
			}
		}
		// 自动创建物流订单
		if s.logisticsService != nil {
			logisticsOrder := &model.LogisticsOrder{
				SalesOrderID:   order.ID,
				SalesOrderNo:   order.OrderNo,
				ReceiverName:   order.CustomerName,
				ReceiverPhone:  order.CustomerPhone,
				ReceiverAddress: order.CustomerAddress,
				Status:         "pending",
			}
			if err := s.logisticsService.Create(logisticsOrder); err != nil {
				return fmt.Errorf("创建物流订单失败: %w", err)
			}
		}

	case "shipping":
		// 发货：更新物流状态为运输中
		if s.logisticsService != nil {
			if err := s.logisticsService.UpdateStatusBySalesOrderID(order.ID, "in_transit"); err != nil {
				// 物流订单可能不存在，忽略错误
			}
		}

	case "completed":
		// 完成：更新物流状态为已签收
		if s.logisticsService != nil {
			if err := s.logisticsService.UpdateStatusBySalesOrderID(order.ID, "delivered"); err != nil {
				// 物流订单可能不存在，忽略错误
			}
		}

	case "cancelled":
		// 取消订单：根据当前状态处理库存
		if s.inventoryService != nil {
			switch order.Status {
			case "pending":
				// 待确认状态取消：解锁库存
				for _, item := range order.Items {
					if item.ProductID == 0 || item.Quantity <= 0 {
						continue
					}
					s.inventoryService.UnlockStock(uint(item.ProductID), uint(item.Quantity), "")
				}
			case "confirmed", "shipping":
				// 已确认/已发货状态取消：退货入库
				for _, item := range order.Items {
					if item.ProductID == 0 || item.Quantity <= 0 {
						continue
					}
					s.inventoryService.StockIn(uint(item.ProductID), uint(item.Quantity), "")
				}
			}
		}
		// 取消物流订单
		if s.logisticsService != nil {
			s.logisticsService.UpdateStatusBySalesOrderID(order.ID, "cancelled")
		}
	}

	order.Status = status
	return s.repo.Update(order)
}

type LogisticsService struct {
	repo *repository.LogisticsRepository
}

func NewLogisticsService(r *repository.LogisticsRepository) *LogisticsService {
	return &LogisticsService{repo: r}
}

func (s *LogisticsService) Create(order *model.LogisticsOrder) error {
	// 自动生成物流单号: LOG + 年月日 + 6位随机数
	order.TrackingNo = fmt.Sprintf("LOG%s%06d", time.Now().Format("20060102"), time.Now().Nanosecond()%1000000)
	return s.repo.Create(order)
}

func (s *LogisticsService) GetByID(id uint) (*model.LogisticsOrder, error) {
	return s.repo.GetByID(id)
}

func (s *LogisticsService) Update(order *model.LogisticsOrder) error {
	return s.repo.Update(order)
}

func (s *LogisticsService) Delete(id uint) error {
	return s.repo.Delete(id)
}

func (s *LogisticsService) List(page, pageSize int, status string) ([]model.LogisticsOrder, int64, error) {
	return s.repo.List(page, pageSize, status)
}

func (s *LogisticsService) UpdateStatus(id uint, status string) error {
	order, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	order.Status = status
	if status == "delivered" {
		now := time.Now()
		order.ActualDelivery = &now
	}
	return s.repo.Update(order)
}

// UpdateStatusBySalesOrderID 根据销售订单ID更新物流状态
func (s *LogisticsService) UpdateStatusBySalesOrderID(salesOrderID uint, status string) error {
	return s.repo.UpdateStatusBySalesOrderID(salesOrderID, status)
}

// CountByStatus 按状态统计物流订单数量
func (s *LogisticsService) CountByStatus(status string) (int64, error) {
	return s.repo.CountByStatus(status)
}
