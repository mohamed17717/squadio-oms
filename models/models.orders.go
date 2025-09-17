package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

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

func (o *Order) BeforeUpdate(tx *gorm.DB) error {
	o.UpdatedAt = time.Now()
	return nil
}

func (p *Payment) BeforeUpdate(tx *gorm.DB) error {
	p.UpdatedAt = time.Now()
	return nil
}
