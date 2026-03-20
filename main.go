package main

import (
	"log"
	"net/http"
	"os"

	"inventory/data"
	"inventory/handlers"

	"github.com/gorilla/mux"
)

func main() {
	// Initialize database
	db := data.NewCSVDB()
	if err := db.Initialize(); err != nil {
		log.Printf("Warning: Failed to initialize database: %v", err)
	}

	// Create handler
	handler := handlers.NewHandler(db)

	// Create router
	router := mux.NewRouter()

	// Apply middleware
	router.Use(handlers.CORSMiddleware)
	router.Use(handlers.LoggingMiddleware)

	// API routes
	api := router.PathPrefix("/api").Subrouter()

	// Product routes
	api.HandleFunc("/products", handler.GetProducts).Methods("GET")
	api.HandleFunc("/products", handler.CreateProduct).Methods("POST")
	api.HandleFunc("/products/{id}", handler.GetProduct).Methods("GET")
	api.HandleFunc("/products/{id}", handler.UpdateProduct).Methods("PUT")
	api.HandleFunc("/products/{id}", handler.DeleteProduct).Methods("DELETE")

	// Category routes
	api.HandleFunc("/categories", handler.GetCategories).Methods("GET")
	api.HandleFunc("/categories", handler.CreateCategory).Methods("POST")

	// Order routes
	api.HandleFunc("/orders", handler.GetOrders).Methods("GET")
	api.HandleFunc("/orders", handler.CreateOrder).Methods("POST")

	// Dashboard routes
	api.HandleFunc("/dashboard/stats", handler.GetDashboardStats).Methods("GET")
	api.HandleFunc("/dashboard/top-products", handler.GetTopProducts).Methods("GET")

	// Serve static files (frontend)
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./")))

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on http://localhost:%s", port)
	log.Printf("API endpoints:")
	log.Printf("  GET    /api/products")
	log.Printf("  POST   /api/products")
	log.Printf("  GET    /api/products/{id}")
	log.Printf("  PUT    /api/products/{id}")
	log.Printf("  DELETE /api/products/{id}")
	log.Printf("  GET    /api/categories")
	log.Printf("  POST   /api/categories")
	log.Printf("  GET    /api/orders")
	log.Printf("  POST   /api/orders")
	log.Printf("  GET    /api/dashboard/stats")
	log.Printf("  GET    /api/dashboard/top-products")

	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatal(err)
	}
}
