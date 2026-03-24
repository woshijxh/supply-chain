package service

import (
	"fmt"
	"supply-chain-server/internal/model"
	"supply-chain-server/internal/repository"
	"time"
)

type SalesReturnService struct {
	repo            *repository.SalesReturnRepository
	salesRepo       *repository.SalesRepository
	inventoryService *InventoryService
}

func NewSalesReturnService(r *repository.SalesReturnRepository, sr *repository.SalesRepository, is *InventoryService) *SalesReturnService {
	return &SalesReturnService{repo: r, salesRepo: sr, inventoryService: is}
}

func (s *SalesReturnService) Create(ret *model.SalesReturn) error {
	// 生成退货单号
	ret.ReturnNo = fmt.Sprintf("SR%s%06d", time.Now().Format("20060102"), time.Now().Nanosecond()%1000000)
	ret.Status = "pending"
	ret.RefundStatus = "pending"

	// 计算总金额
	totalAmount := 0.0
	for i := range ret.Items {
		ret.Items[i].Amount = float64(ret.Items[i].Quantity) * ret.Items[i].UnitPrice
		totalAmount += ret.Items[i].Amount
	}
	ret.TotalAmount = totalAmount

	return s.repo.Create(ret)
}

func (s *SalesReturnService) GetByID(id uint) (*model.SalesReturn, error) {
	return s.repo.GetByID(id)
}

func (s *SalesReturnService) Update(ret *model.SalesReturn) error {
	return s.repo.Update(ret)
}

func (s *SalesReturnService) Delete(id uint) error {
	return s.repo.Delete(id)
}

func (s *SalesReturnService) List(page, pageSize int, status string) ([]model.SalesReturn, int64, error) {
	return s.repo.List(page, pageSize, status)
}

// Approve 批准退货
func (s *SalesReturnService) Approve(id uint) error {
	ret, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	if ret.Status != "pending" {
		return fmt.Errorf("只有待处理状态的退货单可以批准")
	}

	ret.Status = "approved"
	return s.repo.Update(ret)
}

// Reject 拒绝退货
func (s *SalesReturnService) Reject(id uint, reason string) error {
	ret, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	if ret.Status != "pending" {
		return fmt.Errorf("只有待处理状态的退货单可以拒绝")
	}

	ret.Status = "rejected"
	if reason != "" {
		ret.Remark = reason
	}
	return s.repo.Update(ret)
}

// Complete 完成退货（入库+退款）
func (s *SalesReturnService) Complete(id uint, refundAmount float64) error {
	ret, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	if ret.Status != "approved" {
		return fmt.Errorf("只有已批准状态的退货单可以完成")
	}

	// 退货入库
	if s.inventoryService != nil {
		for _, item := range ret.Items {
			if item.ProductID > 0 && item.Quantity > 0 {
				s.inventoryService.StockIn(item.ProductID, uint(item.Quantity), "")
			}
		}
	}

	ret.Status = "completed"
	ret.RefundStatus = "refunded"
	ret.RefundAmount = refundAmount
	return s.repo.Update(ret)
}

// CountByStatus 按状态统计退货数量
func (s *SalesReturnService) CountByStatus(status string) (int64, error) {
	return s.repo.CountByStatus(status)
}

type ProcurementReturnService struct {
	repo             *repository.ProcurementReturnRepository
	procurementRepo  *repository.ProcurementRepository
	inventoryService *InventoryService
}

func NewProcurementReturnService(r *repository.ProcurementReturnRepository, pr *repository.ProcurementRepository, is *InventoryService) *ProcurementReturnService {
	return &ProcurementReturnService{repo: r, procurementRepo: pr, inventoryService: is}
}

func (s *ProcurementReturnService) Create(ret *model.ProcurementReturn) error {
	// 生成退货单号
	ret.ReturnNo = fmt.Sprintf("PR%s%06d", time.Now().Format("20060102"), time.Now().Nanosecond()%1000000)
	ret.Status = "pending"
	ret.RefundStatus = "pending"

	// 计算总金额
	totalAmount := 0.0
	for i := range ret.Items {
		ret.Items[i].Amount = float64(ret.Items[i].Quantity) * ret.Items[i].UnitPrice
		totalAmount += ret.Items[i].Amount
	}
	ret.TotalAmount = totalAmount

	return s.repo.Create(ret)
}

func (s *ProcurementReturnService) GetByID(id uint) (*model.ProcurementReturn, error) {
	return s.repo.GetByID(id)
}

func (s *ProcurementReturnService) Update(ret *model.ProcurementReturn) error {
	return s.repo.Update(ret)
}

func (s *ProcurementReturnService) Delete(id uint) error {
	return s.repo.Delete(id)
}

func (s *ProcurementReturnService) List(page, pageSize int, status string) ([]model.ProcurementReturn, int64, error) {
	return s.repo.List(page, pageSize, status)
}

// Approve 批准退货
func (s *ProcurementReturnService) Approve(id uint) error {
	ret, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	if ret.Status != "pending" {
		return fmt.Errorf("只有待处理状态的退货单可以批准")
	}

	ret.Status = "approved"
	return s.repo.Update(ret)
}

// Reject 拒绝退货
func (s *ProcurementReturnService) Reject(id uint, reason string) error {
	ret, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	if ret.Status != "pending" {
		return fmt.Errorf("只有待处理状态的退货单可以拒绝")
	}

	ret.Status = "rejected"
	if reason != "" {
		ret.Remark = reason
	}
	return s.repo.Update(ret)
}

// Complete 完成退货（出库+退款）
func (s *ProcurementReturnService) Complete(id uint, refundAmount float64) error {
	ret, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	if ret.Status != "approved" {
		return fmt.Errorf("只有已批准状态的退货单可以完成")
	}

	// 退货出库
	if s.inventoryService != nil {
		for _, item := range ret.Items {
			if item.ProductID > 0 && item.Quantity > 0 {
				s.inventoryService.StockOut(item.ProductID, uint(item.Quantity), "")
			}
		}
	}

	ret.Status = "completed"
	ret.RefundStatus = "refunded"
	ret.RefundAmount = refundAmount
	return s.repo.Update(ret)
}