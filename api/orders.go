package api

import (
	"oms-services/config"
	"oms-services/models"
	"oms-services/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Request DTOs
type OrderRequest struct {
	CustomerID *uuid.UUID `json:"customer_id"`
	Currency   string     `json:"currency" binding:"required,len=3"`
	Status     *string    `json:"status"`
}

func OrderRequestToModel(o *OrderRequest) models.Order {
	return models.Order{
		CustomerID: o.CustomerID,
		Status:     models.OrderStatusDraft,
		Currency:   o.Currency,
	}
}

type OrderItemRequest struct {
	OrderID   uuid.UUID `json:"order_id" binding:"required"`
	VariantID uuid.UUID `json:"variant_id" binding:"required"`
	Quantity  int       `json:"quantity" binding:"required,min=1"`
}

func OrderItemRequestToModel(o *OrderItemRequest) models.OrderItem {
	return models.OrderItem{
		OrderID:   o.OrderID,
		VariantID: o.VariantID,
		Quantity:  o.Quantity,
	}
}

// RegisterOrderRoutes registers all order routes
func RegisterOrderRoutes() {
	api := config.Server.Group("/api/v1")

	// Order routes
	orderViewSet := utils.ViewSet[models.Order, OrderRequest, OrderRequest]{
		DB: config.DB,
		PerformCreateFunc: func(c *gin.Context, obj *models.Order) error {
			return nil
		},
		InputOfCreateToModel: OrderRequestToModel,
		InputOfUpdateToModel: OrderRequestToModel,
	}
	api.POST("/orders", orderViewSet.Create)
	api.GET("/orders", orderViewSet.List)
	api.GET("/orders/:id", orderViewSet.Retrieve)
	api.PATCH("/orders/:id", orderViewSet.Update)

	// Order items routes
	orderItemViewSet := utils.ViewSet[models.OrderItem, OrderItemRequest, OrderItemRequest]{
		DB: config.DB,
		PerformCreateFunc: func(c *gin.Context, obj *models.OrderItem) error {
			return nil
		},
		InputOfCreateToModel: OrderItemRequestToModel,
		InputOfUpdateToModel: OrderItemRequestToModel,
	}
	api.GET("/orders/items", orderItemViewSet.List)
	api.POST("/orders/items", orderItemViewSet.Create)
	api.PATCH("/orders/items/:item_id", orderItemViewSet.Update)
	api.DELETE("/orders/items/:item_id", orderItemViewSet.Delete)

	// Checkout routes
	// api.POST("/orders/:id/checkout/preview", CheckoutPreview)
	// api.POST("/orders/:id/checkout/confirm", CheckoutConfirm)

	// Refunds routes
	// api.POST("/orders/:id/refunds", CreateRefund)

	// Events routes
	// api.GET("/orders/:id/events", GetOrderEvents)
}
