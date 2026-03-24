package model

import "time"

// TraceResult 追溯结果
type TraceResult struct {
	Product     *ProductInfo       `json:"product"`     // 产品信息
	Procurement []ProcurementTrace `json:"procurement"` // 采购记录
	Inventory   []InventoryTrace   `json:"inventory"`   // 库存流水
	Sales       []SalesTrace       `json:"sales"`       // 销售记录
	Logistics   []LogisticsTrace   `json:"logistics"`   // 物流记录
	Returns     []ReturnTrace      `json:"returns"`     // 退货记录
	Timeline    []TimelineEvent    `json:"timeline"`    // 时间线
}

// ProductInfo 产品追溯信息
type ProductInfo struct {
	ID         uint    `json:"id"`
	Code       string  `json:"code"`
	Name       string  `json:"name"`
	SKU        string  `json:"sku"`
	Category   string  `json:"category"`
	Unit       string  `json:"unit"`
	CostPrice  float64 `json:"costPrice"`
	SalePrice  float64 `json:"salePrice"`
	TotalStock int     `json:"totalStock"` // 当前总库存
}

// ProcurementTrace 采购追溯节点
type ProcurementTrace struct {
	OrderID      uint       `json:"orderId"`
	OrderNo      string     `json:"orderNo"`
	SupplierID   uint       `json:"supplierId"`
	SupplierName string     `json:"supplierName"`
	Quantity     int        `json:"quantity"`
	ReceivedQty  int        `json:"receivedQty"`
	UnitPrice    float64    `json:"unitPrice"`
	OrderDate    time.Time  `json:"orderDate"`
	ActualDate   *time.Time `json:"actualDate"`
	Warehouse    string     `json:"warehouse"`
	Status       string     `json:"status"`
}

// InventoryTrace 库存追溯节点
type InventoryTrace struct {
	ID        uint      `json:"id"`
	Type      string    `json:"type"`      // in/out/lock/unlock
	Quantity  int       `json:"quantity"`
	BeforeQty int       `json:"beforeQty"`
	AfterQty  int       `json:"afterQty"`
	Warehouse string    `json:"warehouse"`
	BatchNo   string    `json:"batchNo"`
	RefType   string    `json:"refType"`
	RefNo     string    `json:"refNo"`
	Operator  string    `json:"operator"`
	Remark    string    `json:"remark"`
	CreatedAt time.Time `json:"createdAt"`
}

// SalesTrace 销售追溯节点
type SalesTrace struct {
	OrderID        uint       `json:"orderId"`
	OrderNo        string     `json:"orderNo"`
	CustomerName   string     `json:"customerName"`
	CustomerPhone  string     `json:"customerPhone"`
	Quantity       int        `json:"quantity"`
	UnitPrice      float64    `json:"unitPrice"`
	OrderDate      time.Time  `json:"orderDate"`
	DeliveryDate   *time.Time `json:"deliveryDate"`
	Status         string     `json:"status"`
	PaymentStatus  string     `json:"paymentStatus"`
	TrackingNo     string     `json:"trackingNo"`
	LogisticsStatus string    `json:"logisticsStatus"`
}

// LogisticsTrace 物流追溯节点
type LogisticsTrace struct {
	LogisticsID       uint                   `json:"logisticsId"`
	TrackingNo        string                 `json:"trackingNo"`
	Carrier           string                 `json:"carrier"`
	SalesOrderNo      string                 `json:"salesOrderNo"`
	ReceiverName      string                 `json:"receiverName"`
	ReceiverPhone     string                 `json:"receiverPhone"`
	Status            string                 `json:"status"`
	ShippingFee       float64                `json:"shippingFee"`
	EstimatedDelivery *time.Time             `json:"estimatedDelivery"`
	ActualDelivery    *time.Time             `json:"actualDelivery"`
	Timeline          []LogisticsTimelineItem `json:"timeline"`
}

// LogisticsTimelineItem 物流轨迹项
type LogisticsTimelineItem struct {
	Time        time.Time `json:"time"`
	Status      string    `json:"status"`
	Location    string    `json:"location"`
	Description string    `json:"description"`
}

// ReturnTrace 退货追溯节点
type ReturnTrace struct {
	ReturnID    uint      `json:"returnId"`
	ReturnNo    string    `json:"returnNo"`
	Type        string    `json:"type"` // sales_return/procurement_return
	OrderNo     string    `json:"orderNo"`
	RelatedName string    `json:"relatedName"` // 供应商或客户名
	Quantity    int       `json:"quantity"`
	UnitPrice   float64   `json:"unitPrice"`
	Reason      string    `json:"reason"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
}

// TimelineEvent 时间线事件
type TimelineEvent struct {
	Time        time.Time `json:"time"`
	Type        string    `json:"type"`        // procurement/inventory/sales/logistics/return
	Action      string    `json:"action"`      // 入库/出库/发货/签收/退货
	Title       string    `json:"title"`
	Description string    `json:"description"`
	RefNo       string    `json:"refNo"`
	RefType     string    `json:"refType"`
	Operator    string    `json:"operator"`
	Quantity    int       `json:"quantity"`
	Amount      float64   `json:"amount"`
}