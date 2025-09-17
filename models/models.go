package models

import (
	"gorm.io/gorm"
)

// AutoMigrate runs database migrations for all models
func AutoMigrate(db *gorm.DB) error {
	// Enable UUID extension
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error; err != nil {
		return err
	}

	// Create enums
	if err := CreateEnums(db); err != nil {
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
		&Customer{},
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
