package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"inventory/data"
	"inventory/models"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// Handler holds the database reference
type Handler struct {
	DB *data.CSVDB
}

// NewHandler creates a new handler
func NewHandler(db *data.CSVDB) *Handler {
	return &Handler{DB: db}
}

// Response represents a standard API response
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// Product Handlers

func (h *Handler) GetProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.DB.GetAllProducts()
	if err != nil {
		sendError(w, "Failed to get products", http.StatusInternalServerError)
		return
	}
	sendSuccess(w, products)
}

func (h *Handler) GetProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	product, err := h.DB.GetProduct(id)
	if err != nil {
		sendError(w, "Failed to get product", http.StatusInternalServerError)
		return
	}
	if product == nil {
		sendError(w, "Product not found", http.StatusNotFound)
		return
	}
	sendSuccess(w, product)
}

func (h *Handler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		sendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var product models.Product
	if err := json.Unmarshal(body, &product); err != nil {
		sendError(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	product.ID = uuid.New().String()
	product.CreatedAt = time.Now()
	product.UpdatedAt = time.Now()

	if err := h.DB.AddProduct(product); err != nil {
		sendError(w, "Failed to create product", http.StatusInternalServerError)
		return
	}

	sendSuccessWithMessage(w, product, "Product created successfully")
}

func (h *Handler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	body, err := io.ReadAll(r.Body)
	if err != nil {
		sendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var product models.Product
	if err := json.Unmarshal(body, &product); err != nil {
		sendError(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	product.ID = id
	product.UpdatedAt = time.Now()

	if err := h.DB.UpdateProduct(product); err != nil {
		sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sendSuccessWithMessage(w, product, "Product updated successfully")
}

func (h *Handler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.DB.DeleteProduct(id); err != nil {
		sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sendSuccessWithMessage(w, nil, "Product deleted successfully")
}

// Category Handlers

func (h *Handler) GetCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.DB.GetAllCategories()
	if err != nil {
		sendError(w, "Failed to get categories", http.StatusInternalServerError)
		return
	}
	sendSuccess(w, categories)
}

func (h *Handler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		sendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var category models.Category
	if err := json.Unmarshal(body, &category); err != nil {
		sendError(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	category.ID = uuid.New().String()
	category.CreatedAt = time.Now()

	if err := h.DB.AddCategory(category); err != nil {
		sendError(w, "Failed to create category", http.StatusInternalServerError)
		return
	}

	sendSuccessWithMessage(w, category, "Category created successfully")
}

// Order Handlers

func (h *Handler) GetOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := h.DB.GetAllOrders()
	if err != nil {
		sendError(w, "Failed to get orders", http.StatusInternalServerError)
		return
	}
	sendSuccess(w, orders)
}

func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		sendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var order models.Order
	if err := json.Unmarshal(body, &order); err != nil {
		sendError(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	order.ID = uuid.New().String()
	order.CreatedAt = time.Now()

	// Calculate total price based on product price
	if order.ProductID != "" {
		product, err := h.DB.GetProduct(order.ProductID)
		if err == nil && product != nil {
			order.TotalPrice = product.Price * float64(order.Quantity)
		}
	}

	if err := h.DB.AddOrder(order); err != nil {
		sendError(w, "Failed to create order", http.StatusInternalServerError)
		return
	}

	sendSuccessWithMessage(w, order, "Order created successfully")
}

// Dashboard Handler

func (h *Handler) GetDashboardStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.DB.GetDashboardStats()
	if err != nil {
		sendError(w, "Failed to get dashboard stats", http.StatusInternalServerError)
		return
	}
	sendSuccess(w, stats)
}

// Top Products Handler (for charts)
func (h *Handler) GetTopProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.DB.GetAllProducts()
	if err != nil {
		sendError(w, "Failed to get products", http.StatusInternalServerError)
		return
	}

	// Return top 5 products by quantity
	type TopProduct struct {
		Name     string `json:"name"`
		Quantity int    `json:"quantity"`
	}

	// Simple sort by quantity descending
	n := len(products)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			if products[j].Quantity < products[j+1].Quantity {
				products[j], products[j+1] = products[j+1], products[j]
			}
		}
	}

	topProducts := make([]TopProduct, 0)
	count := 5
	if len(products) < count {
		count = len(products)
	}
	for i := 0; i < count; i++ {
		topProducts = append(topProducts, TopProduct{
			Name:     products[i].Name,
			Quantity: products[i].Quantity,
		})
	}

	sendSuccess(w, topProducts)
}

// Helper functions

func sendSuccess(w http.ResponseWriter, data interface{}) {
	sendJSON(w, Response{Success: true, Data: data}, http.StatusOK)
}

func sendSuccessWithMessage(w http.ResponseWriter, data interface{}, message string) {
	sendJSON(w, Response{Success: true, Message: message, Data: data}, http.StatusOK)
}

func sendError(w http.ResponseWriter, message string, statusCode int) {
	sendJSON(w, Response{Success: false, Message: message}, statusCode)
}

func sendJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// parseQueryInt parses an integer query parameter
func parseQueryInt(r *http.Request, key string, defaultValue int) int {
	value := r.URL.Query().Get(key)
	if value == "" {
		return defaultValue
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return intValue
}

// CORS middleware
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		next.ServeHTTP(w, r)
	})
}

// Logging middleware
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("[%s] %s %s\n", time.Now().Format(time.RFC3339), r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}
