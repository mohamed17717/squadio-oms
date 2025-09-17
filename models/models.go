package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Enum types
type OrderStatus string

const (
	OrderStatusDraft                 OrderStatus = "draft"
	OrderStatusPendingPayment        OrderStatus = "pending_payment"
	OrderStatusPaid                  OrderStatus = "paid"
	OrderStatusFulfillmentInProgress OrderStatus = "fulfillment_in_progress"
	OrderStatusShipped               OrderStatus = "shipped"
	OrderStatusCompleted             OrderStatus = "completed"
	OrderStatusCancelled             OrderStatus = "cancelled"
)

type PaymentStatus string

const (
	PaymentStatusPending         PaymentStatus = "pending"
	PaymentStatusAuthorized      PaymentStatus = "authorized"
	PaymentStatusCaptured        PaymentStatus = "captured"
	PaymentStatusFailed          PaymentStatus = "failed"
	PaymentStatusRefunded        PaymentStatus = "refunded"
	PaymentStatusPartialRefunded PaymentStatus = "partial_refunded"
)

type ShipmentStatus string

const (
	ShipmentStatusPending   ShipmentStatus = "pending"
	ShipmentStatusPacked    ShipmentStatus = "packed"
	ShipmentStatusInTransit ShipmentStatus = "in_transit"
	ShipmentStatusDelivered ShipmentStatus = "delivered"
	ShipmentStatusFailed    ShipmentStatus = "failed"
)

type RefundStatus string

const (
	RefundStatusPending   RefundStatus = "pending"
	RefundStatusApproved  RefundStatus = "approved"
	RefundStatusRejected  RefundStatus = "rejected"
	RefundStatusProcessed RefundStatus = "processed"
)

// Product represents a product in the system
type Product struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Title       string    `gorm:"type:text;not null" json:"title" validate:"required"`
	Description *string   `gorm:"type:text" json:"description"`
	IsActive    bool      `gorm:"not null;default:true" json:"is_active"`
	CreatedAt   time.Time `gorm:"type:timestamptz;not null;default:now()" json:"created_at"`
	UpdatedAt   time.Time `gorm:"type:timestamptz;not null;default:now()" json:"updated_at"`

	// Relationships
	Variants []ProductVariant `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE" json:"variants,omitempty"`
}

// ProductVariant represents a specific variant of a product
type ProductVariant struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	ProductID  uuid.UUID `gorm:"type:uuid;not null" json:"product_id"`
	SKU        string    `gorm:"type:text;uniqueIndex;not null" json:"sku" validate:"required"`
	Attributes *string   `gorm:"type:jsonb" json:"attributes"` // JSON for color, size, etc.
	PriceMinor int       `gorm:"not null;check:price_minor >= 0" json:"price_minor" validate:"min=0"`
	Currency   string    `gorm:"type:char(3);not null" json:"currency" validate:"required,len=3"`
	IsActive   bool      `gorm:"not null;default:true" json:"is_active"`
	CreatedAt  time.Time `gorm:"type:timestamptz;not null;default:now()" json:"created_at"`
	UpdatedAt  time.Time `gorm:"type:timestamptz;not null;default:now()" json:"updated_at"`

	// Relationships
	Inventory Inventory `gorm:"foreignKey:VariantID;constraint:OnDelete:CASCADE" json:"inventory,omitempty"`
}

// Inventory represents stock levels for a product variant
type Inventory struct {
	VariantID   uuid.UUID `gorm:"type:uuid;primary_key" json:"variant_id"`
	QtyOnHand   int       `gorm:"not null;default:0;check:qty_on_hand >= 0" json:"qty_on_hand" validate:"min=0"`
	QtyReserved int       `gorm:"not null;default:0;check:qty_reserved >= 0" json:"qty_reserved" validate:"min=0"`
	UpdatedAt   time.Time `gorm:"type:timestamptz;not null;default:now()" json:"updated_at"`

	// Relationships - VariantID is the foreign key to ProductVariant
}

// Order represents a customer order
type Order struct {
	ID            uuid.UUID   `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	CustomerID    *uuid.UUID  `gorm:"type:uuid" json:"customer_id"`
	Status        OrderStatus `gorm:"type:order_status;not null;default:'draft'" json:"status"`
	SubtotalMinor int         `gorm:"not null;default:0;check:subtotal_minor >= 0" json:"subtotal_minor" validate:"min=0"`
	// DiscountMinor     int         `gorm:"not null;default:0;check:discount_minor >= 0" json:"discount_minor" validate:"min=0"`
	// ShippingMinor     int         `gorm:"not null;default:0;check:shipping_minor >= 0" json:"shipping_minor" validate:"min=0"`
	// TaxMinor          int         `gorm:"not null;default:0;check:tax_minor >= 0" json:"tax_minor" validate:"min=0"`
	TotalMinor int    `gorm:"not null;default:0;check:total_minor >= 0" json:"total_minor" validate:"min=0"`
	Currency   string `gorm:"type:char(3);not null" json:"currency" validate:"required,len=3"`
	// BillingAddressID  *uuid.UUID  `gorm:"type:uuid" json:"billing_address_id"`
	// ShippingAddressID *uuid.UUID  `gorm:"type:uuid" json:"shipping_address_id"`
	CreatedAt time.Time `gorm:"type:timestamptz;not null;default:now()" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamptz;not null;default:now()" json:"updated_at"`
	Version   int       `gorm:"not null;default:1" json:"version"` // Optimistic locking

	// Relationships
	Items    []OrderItem  `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE" json:"items,omitempty"`
	Payments []Payment    `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE" json:"payments,omitempty"`
	Refunds  []Refund     `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE" json:"refunds,omitempty"`
	Events   []OrderEvent `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE" json:"events,omitempty"`
}

// OrderItem represents an item within an order
type OrderItem struct {
	ID             uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	OrderID        uuid.UUID `gorm:"type:uuid;not null" json:"order_id"`
	VariantID      uuid.UUID `gorm:"type:uuid;not null" json:"variant_id"`
	Quantity       int       `gorm:"not null;check:quantity > 0" json:"quantity" validate:"min=1"`
	UnitPriceMinor int       `gorm:"not null;check:unit_price_minor >= 0" json:"unit_price_minor" validate:"min=0"`
	Currency       string    `gorm:"type:char(3);not null" json:"currency" validate:"required,len=3"`
	// TaxMinor       int       `gorm:"not null;default:0;check:tax_minor >= 0" json:"tax_minor" validate:"min=0"`
	// DiscountMinor  int       `gorm:"not null;default:0;check:discount_minor >= 0" json:"discount_minor" validate:"min=0"`
	LineTotalMinor int `gorm:"not null;check:line_total_minor >= 0" json:"line_total_minor" validate:"min=0"`

	// Relationships
	Order   Order          `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE" json:"order,omitempty"`
	Variant ProductVariant `gorm:"foreignKey:VariantID" json:"variant,omitempty"`
}

// Payment represents a payment for an order
type Payment struct {
	ID          uuid.UUID     `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	OrderID     uuid.UUID     `gorm:"type:uuid;not null" json:"order_id"`
	Provider    string        `gorm:"type:text;not null" json:"provider" validate:"required"`
	Status      PaymentStatus `gorm:"type:payment_status;not null;default:'pending'" json:"status"`
	AmountMinor int           `gorm:"not null;check:amount_minor >= 0" json:"amount_minor" validate:"min=0"`
	Currency    string        `gorm:"type:char(3);not null" json:"currency" validate:"required,len=3"`
	ExternalRef *string       `gorm:"type:text" json:"external_ref"` // Gateway payment_intent id
	CreatedAt   time.Time     `gorm:"type:timestamptz;not null;default:now()" json:"created_at"`
	UpdatedAt   time.Time     `gorm:"type:timestamptz;not null;default:now()" json:"updated_at"`

	// Relationships
	Order   Order    `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE" json:"order,omitempty"`
	Refunds []Refund `gorm:"foreignKey:PaymentID;constraint:OnDelete:SET NULL" json:"refunds,omitempty"`
}

// Refund represents a refund for an order
type Refund struct {
	ID          uuid.UUID    `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	OrderID     uuid.UUID    `gorm:"type:uuid;not null" json:"order_id"`
	PaymentID   *uuid.UUID   `gorm:"type:uuid" json:"payment_id"`
	Status      RefundStatus `gorm:"type:refund_status;not null;default:'pending'" json:"status"`
	AmountMinor int          `gorm:"not null;check:amount_minor >= 0" json:"amount_minor" validate:"min=0"`
	Reason      *string      `gorm:"type:text" json:"reason"`
	CreatedAt   time.Time    `gorm:"type:timestamptz;not null;default:now()" json:"created_at"`
	ProcessedAt *time.Time   `gorm:"type:timestamptz" json:"processed_at"`

	// Relationships
	Order   Order    `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE" json:"order,omitempty"`
	Payment *Payment `gorm:"foreignKey:PaymentID;constraint:OnDelete:SET NULL" json:"payment,omitempty"`
}

// OrderEvent represents an audit event for an order
type OrderEvent struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	OrderID   uuid.UUID `gorm:"type:uuid;not null" json:"order_id"`
	EventType string    `gorm:"type:text;not null" json:"event_type" validate:"required"`
	Payload   *string   `gorm:"type:jsonb" json:"payload"` // JSON payload
	CreatedAt time.Time `gorm:"type:timestamptz;not null;default:now()" json:"created_at"`

	// Relationships
	Order Order `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE" json:"order,omitempty"`
}

// TableName methods for custom table names (optional, GORM will use pluralized names by default)
func (Product) TableName() string {
	return "products"
}

func (ProductVariant) TableName() string {
	return "product_variants"
}

func (Inventory) TableName() string {
	return "inventory"
}

func (Order) TableName() string {
	return "orders"
}

func (OrderItem) TableName() string {
	return "order_items"
}

func (Payment) TableName() string {
	return "payments"
}

func (Refund) TableName() string {
	return "refunds"
}

func (OrderEvent) TableName() string {
	return "order_events"
}

// BeforeCreate hooks for setting default values
func (p *Product) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

func (pv *ProductVariant) BeforeCreate(tx *gorm.DB) error {
	if pv.ID == uuid.Nil {
		pv.ID = uuid.New()
	}
	return nil
}

func (o *Order) BeforeCreate(tx *gorm.DB) error {
	if o.ID == uuid.Nil {
		o.ID = uuid.New()
	}
	return nil
}

func (oi *OrderItem) BeforeCreate(tx *gorm.DB) error {
	if oi.ID == uuid.Nil {
		oi.ID = uuid.New()
	}
	return nil
}

func (p *Payment) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

func (r *Refund) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}

func (oe *OrderEvent) BeforeCreate(tx *gorm.DB) error {
	if oe.ID == uuid.Nil {
		oe.ID = uuid.New()
	}
	return nil
}

// BeforeUpdate hooks for updating timestamps
func (p *Product) BeforeUpdate(tx *gorm.DB) error {
	p.UpdatedAt = time.Now()
	return nil
}

func (pv *ProductVariant) BeforeUpdate(tx *gorm.DB) error {
	pv.UpdatedAt = time.Now()
	return nil
}

func (o *Order) BeforeUpdate(tx *gorm.DB) error {
	o.UpdatedAt = time.Now()
	return nil
}

func (i *Inventory) BeforeUpdate(tx *gorm.DB) error {
	i.UpdatedAt = time.Now()
	return nil
}

func (p *Payment) BeforeUpdate(tx *gorm.DB) error {
	p.UpdatedAt = time.Now()
	return nil
}

// AutoMigrate runs database migrations for all models
func AutoMigrate(db *gorm.DB) error {
	// Enable UUID extension
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error; err != nil {
		return err
	}

	// Create custom types
	if err := db.Exec(`
		DO $$ BEGIN
			CREATE TYPE order_status AS ENUM ('draft','pending_payment','paid','fulfillment_in_progress','shipped','completed','cancelled');
		EXCEPTION
			WHEN duplicate_object THEN null;
		END $$;
	`).Error; err != nil {
		return err
	}

	if err := db.Exec(`
		DO $$ BEGIN
			CREATE TYPE payment_status AS ENUM ('pending','authorized','captured','failed','refunded','partial_refunded');
		EXCEPTION
			WHEN duplicate_object THEN null;
		END $$;
	`).Error; err != nil {
		return err
	}

	if err := db.Exec(`
		DO $$ BEGIN
			CREATE TYPE refund_status AS ENUM ('pending','approved','rejected','processed');
		EXCEPTION
			WHEN duplicate_object THEN null;
		END $$;
	`).Error; err != nil {
		return err
	}

	// Auto migrate all models
	return db.AutoMigrate(
		&Product{},
		&ProductVariant{},
		&Inventory{},
		&Order{},
		&OrderItem{},
		&Payment{},
		&Refund{},
		&OrderEvent{},
	)
}

// CreateIndexes creates additional indexes for better performance
func CreateIndexes(db *gorm.DB) error {
	indexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_orders_status ON orders(status);",
		"CREATE INDEX IF NOT EXISTS idx_payments_order ON payments(order_id);",
		"CREATE INDEX IF NOT EXISTS idx_order_items_order ON order_items(order_id);",
		"CREATE INDEX IF NOT EXISTS idx_orders_open ON orders(status) WHERE status IN ('pending_payment','paid','fulfillment_in_progress');",
		"CREATE INDEX IF NOT EXISTS idx_product_variants_sku ON product_variants(sku);",
		"CREATE INDEX IF NOT EXISTS idx_orders_customer ON orders(customer_id);",
		"CREATE INDEX IF NOT EXISTS idx_order_events_order ON order_events(order_id);",
	}

	for _, indexSQL := range indexes {
		if err := db.Exec(indexSQL).Error; err != nil {
			return err
		}
	}

	return nil
}
