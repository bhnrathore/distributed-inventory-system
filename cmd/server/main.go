package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bhnrathore/distributed-inventory-system/internal/api"
	"github.com/bhnrathore/distributed-inventory-system/internal/repository"
	"github.com/bhnrathore/distributed-inventory-system/internal/service"
)

func main() {
	// Database connection string (from environment or use default for local development)
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/inventory?sslmode=disable"
	}

	// Initialize database
	log.Println("Connecting to database...")
	db, err := repository.NewDatabase(dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize schema
	log.Println("Initializing database schema...")
	if err := db.InitSchema(context.Background()); err != nil {
		log.Fatalf("Failed to initialize schema: %v", err)
	}

	// Initialize repositories
	dbConn := db.GetConnection()
	productRepo := repository.NewPostgresProductRepository(dbConn)
	inventoryRepo := repository.NewPostgresInventoryRepository(dbConn)
	transactionRepo := repository.NewPostgresTransactionRepository(dbConn)

	// Initialize service
	inventoryService := service.NewInventoryService(productRepo, inventoryRepo, transactionRepo)

	// Initialize API handler
	handler := api.NewHandler(inventoryService)

	// Setup routes
	mux := http.NewServeMux()

	// Health check endpoint
	mux.HandleFunc("/health", handler.HealthHandler)

	// Product list and creation
	mux.HandleFunc("GET /api/products", handler.ListProductsHandler)
	mux.HandleFunc("POST /api/products", handler.CreateProductHandler)

	// Product operations (get, update, delete, stock operations, inventory, transactions)
	mux.HandleFunc("/api/products/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		// Stock operations
		if contains(path, "/stock/add") && r.Method == http.MethodPost {
			handler.AddStockHandler(w, r)
		} else if contains(path, "/stock/remove") && r.Method == http.MethodPost {
			handler.RemoveStockHandler(w, r)
		} else if contains(path, "/stock/reserve") && r.Method == http.MethodPost {
			handler.ReserveStockHandler(w, r)
		} else if contains(path, "/stock/unreserve") && r.Method == http.MethodPost {
			handler.UnreserveStockHandler(w, r)
		} else if contains(path, "/inventory") && r.Method == http.MethodGet {
			handler.GetInventoryHandler(w, r)
		} else if contains(path, "/transactions") && r.Method == http.MethodGet {
			handler.GetTransactionsHandler(w, r)
		} else if r.Method == http.MethodGet {
			handler.GetProductHandler(w, r)
		} else if r.Method == http.MethodPut {
			handler.UpdateProductHandler(w, r)
		} else if r.Method == http.MethodDelete {
			handler.DeleteProductHandler(w, r)
		} else {
			api.WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		}
	})

	// Apply middleware
	var h http.Handler = mux
	h = api.RecoveryMiddleware(h)
	h = api.JSONResponseMiddleware(h)
	h = api.LoggingMiddleware(h)

	// Server setup
	server := &http.Server{
		Addr:         ":8080",
		Handler:      h,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		log.Println("Shutting down server...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Printf("Server shutdown error: %v", err)
		}
	}()

	log.Println("Starting server on :8080")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server error: %v", err)
	}

	log.Println("Server stopped")
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
