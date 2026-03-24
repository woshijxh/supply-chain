package service

import (
	"errors"
	"sort"
	"strconv"
	"supply-chain-server/internal/model"
	"supply-chain-server/internal/repository"
)

type TraceService struct {
	traceRepo   *repository.TraceRepository
	productRepo *repository.ProductRepository
}

func NewTraceService(tr *repository.TraceRepository, pr *repository.ProductRepository) *TraceService {
	return &TraceService{traceRepo: tr, productRepo: pr}
}

// Trace 追溯商品流转记录
func (s *TraceService) Trace(code string) (*model.TraceResult, error) {
	// 1. 查找产品
	product, _, err := s.traceRepo.FindProductByCode(code)
	if err != nil {
		return nil, errors.New("未找到对应的产品信息，请检查编码是否正确")
	}

	// 2. 获取各类追溯数据
	result := &model.TraceResult{
		Procurement: []model.ProcurementTrace{},
		Inventory:   []model.InventoryTrace{},
		Sales:       []model.SalesTrace{},
		Logistics:   []model.LogisticsTrace{},
		Returns:     []model.ReturnTrace{},
		Timeline:    []model.TimelineEvent{},
	}

	// 产品信息
	totalStock, _ := s.traceRepo.GetTotalStock(product.ID)
	result.Product = &model.ProductInfo{
		ID:         product.ID,
		Code:       product.Code,
		Name:       product.Name,
		SKU:        product.SKU,
		Category:   product.Category,
		Unit:       product.Unit,
		CostPrice:  product.CostPrice,
		SalePrice:  product.SalePrice,
		TotalStock: totalStock,
	}

	// 采购记录
	result.Procurement, _ = s.traceRepo.GetProcurementTrace(product.ID)

	// 库存流水
	result.Inventory, _ = s.traceRepo.GetInventoryTrace(product.ID)

	// 销售记录
	result.Sales, _ = s.traceRepo.GetSalesTrace(product.ID)

	// 物流记录
	result.Logistics, _ = s.traceRepo.GetLogisticsTrace(product.ID)

	// 退货记录
	result.Returns, _ = s.traceRepo.GetReturnTrace(product.ID)

	// 3. 生成时间线
	result.Timeline = s.buildTimeline(result)

	return result, nil
}

// buildTimeline 构建时间线
func (s *TraceService) buildTimeline(result *model.TraceResult) []model.TimelineEvent {
	var events []model.TimelineEvent

	// 采购节点
	for _, p := range result.Procurement {
		events = append(events, model.TimelineEvent{
			Time:        p.OrderDate,
			Type:        "procurement",
			Action:      "采购入库",
			Title:       "采购订单: " + p.OrderNo,
			Description: "从 " + p.SupplierName + " 采购 " + strconv.Itoa(p.Quantity) + " 件",
			RefNo:       p.OrderNo,
			RefType:     "procurement",
			Quantity:    p.Quantity,
			Amount:      float64(p.Quantity) * p.UnitPrice,
		})
	}

	// 库存节点
	for _, i := range result.Inventory {
		action := s.getInventoryAction(i.Type)
		events = append(events, model.TimelineEvent{
			Time:        i.CreatedAt,
			Type:        "inventory",
			Action:      action,
			Title:       "库存变动: " + i.Warehouse,
			Description: action + " " + strconv.Itoa(i.Quantity) + " 件",
			RefNo:       i.RefNo,
			RefType:     i.RefType,
			Operator:    i.Operator,
			Quantity:    i.Quantity,
		})
	}

	// 销售节点
	for _, sl := range result.Sales {
		events = append(events, model.TimelineEvent{
			Time:        sl.OrderDate,
			Type:        "sales",
			Action:      "销售出库",
			Title:       "销售订单: " + sl.OrderNo,
			Description: "销售给 " + sl.CustomerName + " " + strconv.Itoa(sl.Quantity) + " 件",
			RefNo:       sl.OrderNo,
			RefType:     "sales",
			Quantity:    sl.Quantity,
			Amount:      float64(sl.Quantity) * sl.UnitPrice,
		})
	}

	// 物流节点
	for _, l := range result.Logistics {
		for _, t := range l.Timeline {
			events = append(events, model.TimelineEvent{
				Time:        t.Time,
				Type:        "logistics",
				Action:      t.Status,
				Title:       "物流追踪: " + l.TrackingNo,
				Description: t.Description,
				RefNo:       l.TrackingNo,
				RefType:     "logistics",
			})
		}
	}

	// 退货节点
	for _, r := range result.Returns {
		action := "退货入库"
		if r.Type == "procurement_return" {
			action = "退货出库"
		}
		events = append(events, model.TimelineEvent{
			Time:        r.CreatedAt,
			Type:        "return",
			Action:      action,
			Title:       "退货单: " + r.ReturnNo,
			Description: r.RelatedName + " 退货 " + strconv.Itoa(r.Quantity) + " 件，原因: " + r.Reason,
			RefNo:       r.ReturnNo,
			RefType:     r.Type,
			Quantity:    r.Quantity,
		})
	}

	// 按时间倒序排序
	sort.Slice(events, func(i, j int) bool {
		return events[i].Time.After(events[j].Time)
	})

	return events
}

func (s *TraceService) getInventoryAction(t string) string {
	switch t {
	case "in":
		return "入库"
	case "out":
		return "出库"
	case "lock":
		return "锁定"
	case "unlock":
		return "解锁"
	default:
		return t
	}
}