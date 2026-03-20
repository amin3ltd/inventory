package models

import "time"

// Product represents an inventory product
type Product struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	SKU         string    `json:"sku"`
	Category    string    `json:"category"`
	Quantity    int       `json:"quantity"`
	Price       float64   `json:"price"`
	MinStock    int       `json:"min_stock"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Category represents a product category
type Category struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

// Order represents a purchase or sales order
type Order struct {
	ID         string    `json:"id"`
	Type       string    `json:"type"` // "purchase" or "sales"
	ProductID  string    `json:"product_id"`
	Quantity   int       `json:"quantity"`
	TotalPrice float64   `json:"total_price"`
	Status     string    `json:"status"` // "pending", "completed", "cancelled"
	CreatedAt  time.Time `json:"created_at"`
}

// DashboardStats represents dashboard statistics
type DashboardStats struct {
	TotalProducts      int `json:"total_products"`
	PurchaseOrders     int `json:"purchase_orders"`
	SalesOrders        int `json:"sales_orders"`
	InventoryAlerts    int `json:"inventory_alerts"`
}
