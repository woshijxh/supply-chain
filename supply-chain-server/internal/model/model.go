package model

import (
	"context"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Date 自定义日期类型，只返回年月日
type Date struct {
	time.Time
}

func (d Date) MarshalJSON() ([]byte, error) {
	if d.IsZero() {
		return []byte("null"), nil
	}
	return []byte(`"` + d.Format("2006-01-02") + `"`), nil
}

func (d *Date) UnmarshalJSON(data []byte) error {
	str := string(data)
	if str == "null" {
		*d = Date{}
		return nil
	}
	str = str[1 : len(str)-1] // 去掉引号
	parsed, err := time.Parse("2006-01-02", str)
	if err != nil {
		return err
	}
	*d = Date{parsed}
	return nil
}

// GORM 扫描支持
func (d *Date) Scan(value interface{}) error {
	if value == nil {
		*d = Date{}
		return nil
	}
	t, ok := value.(time.Time)
	if !ok {
		return nil
	}
	*d = Date{t}
	return nil
}

func (d Date) Value() (interface{}, error) {
	if d.IsZero() {
		return nil, nil
	}
	return d.Time, nil
}

// GORM 数据类型
func (Date) GormDataType() string {
	return "date"
}

// GORM 创建时自动转换
func (d Date) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	if d.IsZero() {
		return clause.Expr{SQL: "NULL"}
	}
	return clause.Expr{SQL: "?", Vars: []interface{}{d.Format("2006-01-02")}}
}

type User struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Username     string         `gorm:"uniqueIndex;size:50;not null" json:"username"`
	Password     string         `gorm:"size:255;not null" json:"-"`
	Email        string         `gorm:"uniqueIndex;size:100" json:"email"`
	Phone        string         `gorm:"size:20" json:"phone"`
	Role         string         `gorm:"size:20;default:operator" json:"role"`
	Department   string         `gorm:"size:50" json:"department"`
	Position     string         `gorm:"size:50" json:"position"`
	Avatar       string         `gorm:"size:255" json:"avatar"`
	Status       int8           `gorm:"default:1" json:"status"`
	CreatedAt    time.Time      `json:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
	Roles        []Role         `gorm:"many2many:user_roles" json:"roles,omitempty"`
	Profile      *UserProfile   `gorm:"foreignKey:UserID" json:"profile,omitempty"`
}

type Supplier struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Code         string         `gorm:"uniqueIndex;size:20;not null" json:"code"`
	Name         string         `gorm:"size:100;not null" json:"name" validate:"required,max=100"`
	Contact      string         `gorm:"size:50" json:"contact" validate:"required,max=50"`
	Phone        string         `gorm:"size:20" json:"phone" validate:"required,max=20"`
	Email        string         `gorm:"size:100" json:"email" validate:"omitempty,email,max=100"`
	Address      string         `gorm:"size:255" json:"address" validate:"max=255"`
	Level        string         `gorm:"size:1;default:B" json:"level" validate:"omitempty,oneof=A B C"`
	Category     string         `gorm:"size:50" json:"category" validate:"max=50"`
	PaymentTerms string         `gorm:"size:100" json:"paymentTerms" validate:"max=100"`
	BankName     string         `gorm:"size:100" json:"bankName" validate:"max=100"`
	BankAccount  string         `gorm:"size:50" json:"bankAccount" validate:"max=50"`
	TaxNumber    string         `gorm:"size:50" json:"taxNumber" validate:"max=50"`
	Rating       float64        `gorm:"default:4.0" json:"rating"`
	Status       int8           `gorm:"default:1" json:"status" validate:"oneof=0 1"`
	Remark       string         `gorm:"type:text" json:"remark"`
	Department   string         `gorm:"size:50" json:"department"`
	CreatedBy    *uint          `json:"createdBy"`
	CreatedAt    time.Time      `json:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

type Product struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Code        string         `gorm:"uniqueIndex;size:20;not null" json:"code"`
	Name        string         `gorm:"size:100;not null" json:"name"`
	SKU         string         `gorm:"uniqueIndex;size:50" json:"sku"`
	Category    string         `gorm:"size:50" json:"category"`
	Unit        string         `gorm:"size:20" json:"unit"`
	CostPrice   float64        `json:"costPrice"`
	SalePrice   float64        `json:"salePrice"`
	MinStock    int            `gorm:"default:0" json:"minStock"`
	MaxStock    int            `gorm:"default:0" json:"maxStock"`
	Description string         `gorm:"type:text" json:"description"`
	Status      int8           `gorm:"default:1" json:"status"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type Inventory struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	ProductID    uint           `gorm:"index;not null" json:"productId"`
	Product      Product        `gorm:"foreignKey:ProductID" json:"product"`
	Warehouse    string         `gorm:"size:50;not null" json:"warehouse"`
	Quantity     int            `gorm:"default:0" json:"quantity"`
	AvailableQty int            `gorm:"default:0" json:"availableQty"`
	LockedQty    int            `gorm:"default:0" json:"lockedQty"`
	Location     string         `gorm:"size:50" json:"location"`
	BatchNo      string         `gorm:"size:50" json:"batchNo"`
	ExpiryDate   *time.Time     `json:"expiryDate"`
	Status       string         `gorm:"size:20;default:normal" json:"status"`
	CreatedAt    time.Time      `json:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

type ProcurementOrder struct {
	ID            uint              `gorm:"primaryKey" json:"id"`
	OrderNo       string            `gorm:"uniqueIndex;size:30;not null" json:"orderNo"`
	SupplierID    uint              `gorm:"index;not null" json:"supplierId"`
	Supplier      Supplier          `gorm:"foreignKey:SupplierID" json:"supplier"`
	TotalAmount   float64           `json:"totalAmount"`
	OrderDate     Date              `gorm:"type:date" json:"orderDate"`
	ExpectedDate  *Date             `gorm:"type:date" json:"expectedDate"`
	ActualDate    *Date             `gorm:"type:date" json:"actualDate"`
	Status        string            `gorm:"size:20;default:pending" json:"status"`
	Warehouse     string            `gorm:"size:50" json:"warehouse"`
	Remark        string            `gorm:"type:text" json:"remark"`
	AttachmentURL string            `gorm:"size:500" json:"attachmentUrl"`
	Department    string            `gorm:"size:50" json:"department"`
	CreatedBy     *uint             `json:"createdBy"`
	Items         []ProcurementItem `gorm:"foreignKey:OrderID" json:"items"`
	CreatedAt     time.Time         `json:"createdAt"`
	UpdatedAt     time.Time         `json:"updatedAt"`
	DeletedAt     gorm.DeletedAt    `gorm:"index" json:"-"`
}

type ProcurementItem struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	OrderID     uint      `gorm:"index;not null" json:"orderId"`
	ProductID   uint      `gorm:"index" json:"productId"`
	Product     Product   `gorm:"foreignKey:ProductID" json:"product"`
	ProductName string    `gorm:"size:100" json:"productName"`
	Quantity    int       `json:"quantity"`
	Unit        string    `gorm:"size:20" json:"unit"`
	UnitPrice   float64   `json:"unitPrice"`
	Amount      float64   `json:"amount"`
	ReceivedQty int       `gorm:"default:0" json:"receivedQty"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type SalesOrder struct {
	ID              uint             `gorm:"primaryKey" json:"id"`
	OrderNo         string           `gorm:"uniqueIndex;size:30;not null" json:"orderNo"`
	CustomerID      uint             `gorm:"index" json:"customerId"`
	CustomerName    string           `gorm:"size:100;not null" json:"customerName"`
	CustomerPhone   string           `gorm:"size:20" json:"customerPhone"`
	CustomerAddress string           `gorm:"size:255" json:"customerAddress"`
	TotalAmount     float64          `json:"totalAmount"`
	Discount        float64          `json:"discount"`
	Tax             float64          `json:"tax"`
	ShippingFee     float64          `json:"shippingFee"`
	OrderDate       time.Time        `json:"orderDate"`
	DeliveryDate    *time.Time       `json:"deliveryDate"`
	Status          string           `gorm:"size:20;default:pending" json:"status"`
	PaymentMethod   string           `gorm:"size:20" json:"paymentMethod"`
	PaymentStatus   string           `gorm:"size:20;default:pending" json:"paymentStatus"`
	Remark          string           `gorm:"type:text" json:"remark"`
	Department      string           `gorm:"size:50" json:"department"`
	CreatedBy       *uint            `json:"createdBy"`
	Items           []SalesOrderItem `gorm:"foreignKey:OrderID" json:"items"`
	CreatedAt       time.Time        `json:"createdAt"`
	UpdatedAt       time.Time        `json:"updatedAt"`
	DeletedAt       gorm.DeletedAt   `gorm:"index" json:"-"`
}

type SalesOrderItem struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	OrderID     uint      `gorm:"index;not null" json:"orderId"`
	ProductID   uint      `gorm:"index" json:"productId"`
	ProductName string    `gorm:"size:100" json:"productName"`
	Quantity    int       `json:"quantity"`
	Unit        string    `gorm:"size:20" json:"unit"`
	UnitPrice   float64   `json:"unitPrice"`
	Amount      float64   `json:"amount"`
	Discount    float64   `json:"discount"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type LogisticsOrder struct {
	ID                uint                `gorm:"primaryKey" json:"id"`
	TrackingNo        string              `gorm:"uniqueIndex;size:30;not null" json:"trackingNo"`
	Carrier           string              `gorm:"size:50" json:"carrier"`
	SalesOrderID      uint                `gorm:"index" json:"salesOrderId"`
	SalesOrderNo      string              `gorm:"size:30" json:"salesOrderNo"`
	Status            string              `gorm:"size:20;default:pending" json:"status"`
	SenderName        string              `gorm:"size:50" json:"senderName"`
	SenderPhone       string              `gorm:"size:20" json:"senderPhone"`
	SenderAddress     string              `gorm:"size:255" json:"senderAddress"`
	ReceiverName      string              `gorm:"size:50" json:"receiverName"`
	ReceiverPhone     string              `gorm:"size:20" json:"receiverPhone"`
	ReceiverAddress   string              `gorm:"size:255" json:"receiverAddress"`
	Weight            float64             `json:"weight"`
	ShippingFee       float64             `json:"shippingFee"`
	EstimatedDelivery *time.Time          `json:"estimatedDelivery"`
	ActualDelivery    *time.Time          `json:"actualDelivery"`
	Department        string              `gorm:"size:50" json:"department"`
	CreatedBy         *uint               `json:"createdBy"`
	Timeline          []LogisticsTimeline `gorm:"foreignKey:LogisticsID" json:"timeline"`
	CreatedAt         time.Time           `json:"createdAt"`
	UpdatedAt         time.Time           `json:"updatedAt"`
	DeletedAt         gorm.DeletedAt      `gorm:"index" json:"-"`
}

type LogisticsTimeline struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	LogisticsID uint      `gorm:"index;not null" json:"logisticsId"`
	Time        time.Time `json:"time"`
	Status      string    `gorm:"size:50" json:"status"`
	Location    string    `gorm:"size:100" json:"location"`
	Description string    `gorm:"size:255" json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
}

// ========== RBAC 模型 ==========

// Role 角色表
type Role struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"size:50;not null;uniqueIndex" json:"name"`
	Code        string         `gorm:"size:50" json:"code"`
	Description string         `gorm:"size:255" json:"description"`
	Status      int8           `gorm:"default:1" json:"status"`
	Permissions []Permission  `gorm:"many2many:role_permissions" json:"permissions,omitempty"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// Permission 权限表
type Permission struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"size:100;not null;uniqueIndex" json:"name"`
	Code        string    `gorm:"size:100" json:"code"`          // 权限编码
	Type        string    `gorm:"size:20;default:api" json:"type"` // 权限类型: api, menu, button, data
	Description string    `gorm:"size:255" json:"description"`
	Status      int8      `gorm:"default:1" json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
}

// UserProfile 用户扩展信息（用于行级权限）
type UserProfile struct {
	UserID      uint      `gorm:"primaryKey" json:"userId"`
	Department  string    `gorm:"size:50" json:"department"`
	Position    string    `gorm:"size:50" json:"position"`
	CreatedBy   *uint     `gorm" json:"createdBy"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	User        User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// UserRole 用户-角色关联表
type UserRole struct {
	UserID uint `gorm:"primaryKey"`
	RoleID uint `gorm:"primaryKey"`
	User   User `gorm:"foreignKey:UserID"`
	Role   Role `gorm:"foreignKey:RoleID"`
}

// UserPermission 用户-权限直接关联表
type UserPermission struct {
	UserID       uint       `gorm:"primaryKey"`
	PermissionID uint       `gorm:"primaryKey"`
	User         User       `gorm:"foreignKey:UserID"`
	Permission   Permission `gorm:"foreignKey:PermissionID"`
}
