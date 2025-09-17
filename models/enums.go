package models

import (
	"fmt"
	"strings"

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

func CreateEnumSQLQuery(typeName string, fields []string) string {
	query := fmt.Sprintf(`
		DO $$ BEGIN
			CREATE TYPE %s AS ENUM ('%s');
		EXCEPTION
			WHEN duplicate_object THEN null;
		END $$;`, typeName, strings.Join(fields, "','"))
	return query
}

func CreateEnums(db *gorm.DB) error {
	orderStatusFields := []string{"draft", "pending_payment", "paid", "fulfillment_in_progress", "shipped", "completed", "cancelled"}
	paymentStatusFields := []string{"pending", "authorized", "captured", "failed", "refunded", "partial_refunded"}
	refundStatusFields := []string{"pending", "approved", "rejected", "processed"}

	// Create custom types
	if err := db.Exec(CreateEnumSQLQuery("order_status", orderStatusFields)).Error; err != nil {
		return err
	}

	if err := db.Exec(CreateEnumSQLQuery("payment_status", paymentStatusFields)).Error; err != nil {
		return err
	}

	if err := db.Exec(CreateEnumSQLQuery("refund_status", refundStatusFields)).Error; err != nil {
		return err
	}

	return nil
}
