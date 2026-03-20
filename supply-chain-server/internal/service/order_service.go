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
	order.OrderDate = model.Date{time.Now()}
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
		now := model.Date{time.Now()}
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
	repo *repository.SalesRepository
}

func NewSalesService(r *repository.SalesRepository) *SalesService {
	return &SalesService{repo: r}
}

func (s *SalesService) Create(order *model.SalesOrder) error {
	order.OrderDate = time.Now()
	// 自动生成销售订单号: SO + 年月日 + 6位随机数
	order.OrderNo = fmt.Sprintf("SO%s%06d", time.Now().Format("20060102"), time.Now().Nanosecond()%1000000)
	return s.repo.Create(order)
}

func (s *SalesService) GetByID(id uint) (*model.SalesOrder, error) {
	return s.repo.GetByID(id)
}

func (s *SalesService) Update(order *model.SalesOrder) error {
	return s.repo.Update(order)
}

func (s *SalesService) Delete(id uint) error {
	return s.repo.Delete(id)
}

func (s *SalesService) List(page, pageSize int, status string) ([]model.SalesOrder, int64, error) {
	return s.repo.List(page, pageSize, status)
}

func (s *SalesService) UpdateStatus(id uint, status string) error {
	order, err := s.repo.GetByID(id)
	if err != nil {
		return err
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
