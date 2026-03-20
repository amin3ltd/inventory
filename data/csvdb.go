package data

import (
	"encoding/csv"
	"fmt"
	"os"
	"time"

	"inventory/models"

	"github.com/google/uuid"
)

// CSVDB handles CSV file operations
type CSVDB struct {
	ProductsFile  string
	CategoriesFile string
	OrdersFile    string
}

// NewCSVDB creates a new CSV database handler
func NewCSVDB() *CSVDB {
	return &CSVDB{
		ProductsFile:   "data/products.csv",
		CategoriesFile: "data/categories.csv",
		OrdersFile:     "data/orders.csv",
	}
}

// Initialize creates the CSV files with headers if they don't exist
func (db *CSVDB) Initialize() error {
	if err := db.createFileIfNotExists(db.ProductsFile, []string{"id", "name", "sku", "category", "quantity", "price", "min_stock", "description", "created_at", "updated_at"}); err != nil {
		return err
	}
	if err := db.createFileIfNotExists(db.CategoriesFile, []string{"id", "name", "description", "created_at"}); err != nil {
		return err
	}
	if err := db.createFileIfNotExists(db.OrdersFile, []string{"id", "type", "product_id", "quantity", "total_price", "status", "created_at"}); err != nil {
		return err
	}
	
	// Add sample categories if empty
	categories, _ := db.GetAllCategories()
	if len(categories) == 0 {
		db.AddCategory(models.Category{
			ID:          uuid.New().String(),
			Name:        "Electronics",
			Description: "Electronic devices and accessories",
			CreatedAt:   time.Now(),
		})
		db.AddCategory(models.Category{
			ID:          uuid.New().String(),
			Name:        "Office Supplies",
			Description: "Office and stationery items",
			CreatedAt:   time.Now(),
		})
		db.AddCategory(models.Category{
			ID:          uuid.New().String(),
			Name:        "Furniture",
			Description: "Office and home furniture",
			CreatedAt:   time.Now(),
		})
	}
	
	// Add sample products if empty
	products, _ := db.GetAllProducts()
	if len(products) == 0 {
		products := []models.Product{
			{
				ID: uuid.New().String(), Name: "Laptop", SKU: "LAP-001", Category: "Electronics",
				Quantity: 10, Price: 999.99, MinStock: 5, Description: "High-performance laptop",
				CreatedAt: time.Now(), UpdatedAt: time.Now(),
			},
			{
				ID: uuid.New().String(), Name: "Phone", SKU: "PHO-001", Category: "Electronics",
				Quantity: 25, Price: 699.99, MinStock: 10, Description: "Smartphone with advanced features",
				CreatedAt: time.Now(), UpdatedAt: time.Now(),
			},
			{
				ID: uuid.New().String(), Name: "Monitor", SKU: "MON-001", Category: "Electronics",
				Quantity: 15, Price: 299.99, MinStock: 5, Description: "27-inch 4K monitor",
				CreatedAt: time.Now(), UpdatedAt: time.Now(),
			},
			{
				ID: uuid.New().String(), Name: "Headphones", SKU: "HEA-001", Category: "Electronics",
				Quantity: 50, Price: 149.99, MinStock: 15, Description: "Wireless noise-cancelling headphones",
				CreatedAt: time.Now(), UpdatedAt: time.Now(),
			},
			{
				ID: uuid.New().String(), Name: "Camera", SKU: "CAM-001", Category: "Electronics",
				Quantity: 8, Price: 899.99, MinStock: 3, Description: "Professional DSLR camera",
				CreatedAt: time.Now(), UpdatedAt: time.Now(),
			},
		}
		for _, p := range products {
			db.AddProduct(p)
		}
	}
	
	return nil
}

func (db *CSVDB) createFileIfNotExists(filename string, headers []string) error {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		file, err := os.Create(filename)
		if err != nil {
			return err
		}
		defer file.Close()
		
		writer := csv.NewWriter(file)
		defer writer.Flush()
		
		if err := writer.Write(headers); err != nil {
			return err
		}
	}
	return nil
}

// Product CRUD operations

func (db *CSVDB) GetAllProducts() ([]models.Product, error) {
	file, err := os.Open(db.ProductsFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var products []models.Product
	for i, record := range records {
		if i == 0 {
			continue // skip header
		}
		if len(record) < 10 {
			continue
		}
		
		createdAt, _ := time.Parse(time.RFC3339, record[8])
		updatedAt, _ := time.Parse(time.RFC3339, record[9])
		
		product := models.Product{
			ID:          record[0],
			Name:        record[1],
			SKU:         record[2],
			Category:    record[3],
			Quantity:    parseInt(record[4]),
			Price:       parseFloat(record[5]),
			MinStock:    parseInt(record[6]),
			Description: record[7],
			CreatedAt:   createdAt,
			UpdatedAt:   updatedAt,
		}
		products = append(products, product)
	}
	return products, nil
}

func (db *CSVDB) GetProduct(id string) (*models.Product, error) {
	products, err := db.GetAllProducts()
	if err != nil {
		return nil, err
	}
	for _, p := range products {
		if p.ID == id {
			return &p, nil
		}
	}
	return nil, nil
}

func (db *CSVDB) AddProduct(product models.Product) error {
	file, err := os.OpenFile(db.ProductsFile, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	record := []string{
		product.ID,
		product.Name,
		product.SKU,
		product.Category,
		fmt.Sprintf("%d", product.Quantity),
		fmt.Sprintf("%.2f", product.Price),
		fmt.Sprintf("%d", product.MinStock),
		product.Description,
		product.CreatedAt.Format(time.RFC3339),
		product.UpdatedAt.Format(time.RFC3339),
	}

	if err := writer.Write(record); err != nil {
		return err
	}
	return nil
}

func (db *CSVDB) UpdateProduct(product models.Product) error {
	products, err := db.GetAllProducts()
	if err != nil {
		return err
	}

	for i, p := range products {
		if p.ID == product.ID {
			products[i] = product
			return db.writeProducts(products)
		}
	}
	return fmt.Errorf("product not found")
}

func (db *CSVDB) DeleteProduct(id string) error {
	products, err := db.GetAllProducts()
	if err != nil {
		return err
	}

	for i, p := range products {
		if p.ID == id {
			products = append(products[:i], products[i+1:]...)
			return db.writeProducts(products)
		}
	}
	return fmt.Errorf("product not found")
}

func (db *CSVDB) writeProducts(products []models.Product) error {
	file, err := os.Create(db.ProductsFile)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	if err := writer.Write([]string{"id", "name", "sku", "category", "quantity", "price", "min_stock", "description", "created_at", "updated_at"}); err != nil {
		return err
	}

	for _, p := range products {
		record := []string{
			p.ID,
			p.Name,
			p.SKU,
			p.Category,
			fmt.Sprintf("%d", p.Quantity),
			fmt.Sprintf("%.2f", p.Price),
			fmt.Sprintf("%d", p.MinStock),
			p.Description,
			p.CreatedAt.Format(time.RFC3339),
			p.UpdatedAt.Format(time.RFC3339),
		}
		if err := writer.Write(record); err != nil {
			return err
		}
	}
	return nil
}

// Category CRUD operations

func (db *CSVDB) GetAllCategories() ([]models.Category, error) {
	file, err := os.Open(db.CategoriesFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var categories []models.Category
	for i, record := range records {
		if i == 0 {
			continue
		}
		if len(record) < 4 {
			continue
		}
		
		createdAt, _ := time.Parse(time.RFC3339, record[3])
		
		category := models.Category{
			ID:          record[0],
			Name:        record[1],
			Description: record[2],
			CreatedAt:   createdAt,
		}
		categories = append(categories, category)
	}
	return categories, nil
}

func (db *CSVDB) AddCategory(category models.Category) error {
	file, err := os.OpenFile(db.CategoriesFile, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	record := []string{
		category.ID,
		category.Name,
		category.Description,
		category.CreatedAt.Format(time.RFC3339),
	}

	if err := writer.Write(record); err != nil {
		return err
	}
	return nil
}

// Order CRUD operations

func (db *CSVDB) GetAllOrders() ([]models.Order, error) {
	file, err := os.Open(db.OrdersFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var orders []models.Order
	for i, record := range records {
		if i == 0 {
			continue
		}
		if len(record) < 7 {
			continue
		}
		
		createdAt, _ := time.Parse(time.RFC3339, record[6])
		
		order := models.Order{
			ID:         record[0],
			Type:       record[1],
			ProductID:  record[2],
			Quantity:   parseInt(record[3]),
			TotalPrice: parseFloat(record[4]),
			Status:     record[5],
			CreatedAt:  createdAt,
		}
		orders = append(orders, order)
	}
	return orders, nil
}

func (db *CSVDB) AddOrder(order models.Order) error {
	file, err := os.OpenFile(db.OrdersFile, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	record := []string{
		order.ID,
		order.Type,
		order.ProductID,
		fmt.Sprintf("%d", order.Quantity),
		fmt.Sprintf("%.2f", order.TotalPrice),
		order.Status,
		order.CreatedAt.Format(time.RFC3339),
	}

	if err := writer.Write(record); err != nil {
		return err
	}
	
	// Update product quantity based on order type
	if order.Status == "completed" {
		products, err := db.GetAllProducts()
		if err == nil {
			for i, p := range products {
				if p.ID == order.ProductID {
					if order.Type == "purchase" {
						products[i].Quantity += order.Quantity
					} else if order.Type == "sales" {
						products[i].Quantity -= order.Quantity
					}
					products[i].UpdatedAt = time.Now()
					db.writeProducts(products)
					break
				}
			}
		}
	}
	
	return nil
}

// GetDashboardStats returns statistics for the dashboard
func (db *CSVDB) GetDashboardStats() (models.DashboardStats, error) {
	stats := models.DashboardStats{}
	
	products, err := db.GetAllProducts()
	if err != nil {
		return stats, err
	}
	stats.TotalProducts = len(products)
	
	// Count products below minimum stock
	for _, p := range products {
		if p.Quantity < p.MinStock {
			stats.InventoryAlerts++
		}
	}
	
	orders, err := db.GetAllOrders()
	if err != nil {
		return stats, err
	}
	
	for _, o := range orders {
		if o.Type == "purchase" {
			stats.PurchaseOrders++
		} else if o.Type == "sales" {
			stats.SalesOrders++
		}
	}
	
	return stats, nil
}

// Helper functions
func parseInt(s string) int {
	var n int
	fmt.Sscanf(s, "%d", &n)
	return n
}

func parseFloat(s string) float64 {
	var f float64
	fmt.Sscanf(s, "%f", &f)
	return f
}
