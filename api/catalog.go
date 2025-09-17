package api

import (
	"oms-services/config"
	"oms-services/models"
	"oms-services/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Constants
const (
	TimeFormat = "2006-01-02T15:04:05Z07:00"
)

// Request DTOs
type ProductRequest struct {
	Title       string  `json:"title" binding:"required"`
	Description *string `json:"description"`
	IsActive    *bool   `json:"is_active"`
}

func ProductRequestToModel(p *ProductRequest) models.Product {
	return models.Product{
		Title:       p.Title,
		Description: p.Description,
		IsActive:    *p.IsActive,
	}
}

type VariantRequest struct {
	ProductID  uuid.UUID `json:"product_id" binding:"required"`
	SKU        string    `json:"sku" binding:"required"`
	Attributes *string   `json:"attributes"`
	PriceMinor int       `json:"price_minor" binding:"required,min=0"`
	Currency   string    `json:"currency" binding:"required,len=3"`
	IsActive   *bool     `json:"is_active"`
}

func VariantRequestToModel(v *VariantRequest) models.ProductVariant {
	return models.ProductVariant{
		ProductID:  v.ProductID,
		SKU:        v.SKU,
		Attributes: v.Attributes,
		PriceMinor: v.PriceMinor,
		Currency:   v.Currency,
		IsActive:   *v.IsActive,
	}
}

// Response DTOs

type VariantResponse struct {
	ID         uuid.UUID          `json:"id"`
	ProductID  uuid.UUID          `json:"product_id"`
	SKU        string             `json:"sku"`
	Attributes *string            `json:"attributes"`
	PriceMinor int                `json:"price_minor"`
	Currency   string             `json:"currency"`
	IsActive   bool               `json:"is_active"`
	CreatedAt  string             `json:"created_at"`
	UpdatedAt  string             `json:"updated_at"`
	Inventory  *InventoryResponse `json:"inventory,omitempty"`
}

type InventoryResponse struct {
	VariantID   uuid.UUID `json:"variant_id"`
	QtyOnHand   int       `json:"qty_on_hand"`
	QtyReserved int       `json:"qty_reserved"`
	UpdatedAt   string    `json:"updated_at"`
}

// RegisterCatalogRoutes registers all catalog routes
func RegisterCatalogRoutes() {
	api := config.Server.Group("/api/v1")

	productViewSet := utils.ViewSet[models.Product, ProductRequest, ProductRequest]{
		DB: config.DB,
		PerformCreateFunc: func(c *gin.Context, obj *models.Product) error {
			return nil
		},
		InputOfCreateToModel: ProductRequestToModel,
		InputOfUpdateToModel: ProductRequestToModel,
	}
	// Product routes
	api.GET("/products", productViewSet.List)
	api.POST("/products", productViewSet.Create)
	api.GET("/products/:id", productViewSet.Retrieve)
	api.PATCH("/products/:id", productViewSet.Update)

	variantViewSet := utils.ViewSet[models.ProductVariant, VariantRequest, VariantRequest]{
		DB: config.DB,
		PerformCreateFunc: func(c *gin.Context, obj *models.ProductVariant) error {
			return nil
		},
		InputOfCreateToModel: VariantRequestToModel,
		InputOfUpdateToModel: VariantRequestToModel,
	}

	// Variant routes
	api.GET("/variants", variantViewSet.List)
	api.POST("/variants", variantViewSet.Create)
	api.GET("/variants/:id", variantViewSet.Retrieve)
	api.PATCH("/variants/:id", variantViewSet.Update)

}
