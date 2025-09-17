package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
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

// BeforeUpdate hooks for updating timestamps
func (p *Product) BeforeUpdate(tx *gorm.DB) error {
	p.UpdatedAt = time.Now()
	return nil
}

func (pv *ProductVariant) BeforeUpdate(tx *gorm.DB) error {
	pv.UpdatedAt = time.Now()
	return nil
}

func (i *Inventory) BeforeUpdate(tx *gorm.DB) error {
	i.UpdatedAt = time.Now()
	return nil
}
