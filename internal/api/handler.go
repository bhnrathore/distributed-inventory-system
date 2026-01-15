package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/bhnrathore/distributed-inventory-system/internal/domain"
	"github.com/bhnrathore/distributed-inventory-system/internal/service"
)

// Handler holds references to services
type Handler struct {
	inventoryService *service.InventoryService
}

// NewHandler creates a new API handler
func NewHandler(inventoryService *service.InventoryService) *Handler {
	return &Handler{
		inventoryService: inventoryService,
	}
}

// CreateProductRequest represents a product creation request
type CreateProductRequest struct {
	Name            string  `json:"name"`
	Description     string  `json:"description"`
	SKU             string  `json:"sku"`
	Price           float64 `json:"price"`
	Location        string  `json:"location"`
	InitialQuantity int64   `json:"initial_quantity"`
}

// UpdateProductRequest represents a product update request
type UpdateProductRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

// StockOperationRequest represents a stock operation request
type StockOperationRequest struct {
	Quantity  int64  `json:"quantity"`
	Reference string `json:"reference"`
	Notes     string `json:"notes"`
}

// HealthHandler handles health check requests
func (h *Handler) HealthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Only GET is allowed")
		return
	}

	WriteSuccess(w, http.StatusOK, "Service is healthy", map[string]string{
		"status": "ok",
	})
}

// CreateProductHandler handles product creation
func (h *Handler) CreateProductHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Only POST is allowed")
		return
	}

	var req CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	product := &domain.Product{
		Name:        req.Name,
		Description: req.Description,
		SKU:         req.SKU,
		Price:       req.Price,
	}

	if err := h.inventoryService.CreateProduct(r.Context(), product, req.Location, req.InitialQuantity); err != nil {
		WriteError(w, http.StatusInternalServerError, "CREATION_FAILED", err.Error())
		return
	}

	WriteSuccess(w, http.StatusCreated, "Product created successfully", product)
}

// GetProductHandler handles retrieving a product
func (h *Handler) GetProductHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Only GET is allowed")
		return
	}

	productID := strings.TrimPrefix(r.URL.Path, "/api/products/")
	productID = strings.TrimSuffix(productID, "/")

	product, inventory, err := h.inventoryService.GetProduct(r.Context(), productID)
	if err != nil {
		WriteError(w, http.StatusNotFound, "NOT_FOUND", err.Error())
		return
	}

	response := map[string]interface{}{
		"product":   product,
		"inventory": inventory,
	}

	WriteSuccess(w, http.StatusOK, "Product retrieved successfully", response)
}

// ListProductsHandler handles listing products
func (h *Handler) ListProductsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Only GET is allowed")
		return
	}

	limit := 10
	offset := 0

	if l := r.URL.Query().Get("limit"); l != "" {
		if parsedLimit, err := strconv.Atoi(l); err == nil {
			limit = parsedLimit
		}
	}

	if o := r.URL.Query().Get("offset"); o != "" {
		if parsedOffset, err := strconv.Atoi(o); err == nil {
			offset = parsedOffset
		}
	}

	products, err := h.inventoryService.ListProducts(r.Context(), limit, offset)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "LIST_FAILED", err.Error())
		return
	}

	WriteSuccess(w, http.StatusOK, "Products retrieved successfully", products)
}

// UpdateProductHandler handles product updates
func (h *Handler) UpdateProductHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Only PUT is allowed")
		return
	}

	productID := strings.TrimPrefix(r.URL.Path, "/api/products/")
	productID = strings.TrimSuffix(productID, "/")

	var req UpdateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	// Get existing product
	product, _, err := h.inventoryService.GetProduct(r.Context(), productID)
	if err != nil {
		WriteError(w, http.StatusNotFound, "NOT_FOUND", err.Error())
		return
	}

	// Update fields
	product.Name = req.Name
	product.Description = req.Description
	product.Price = req.Price

	if err := h.inventoryService.UpdateProduct(r.Context(), product); err != nil {
		WriteError(w, http.StatusInternalServerError, "UPDATE_FAILED", err.Error())
		return
	}

	WriteSuccess(w, http.StatusOK, "Product updated successfully", product)
}

// DeleteProductHandler handles product deletion
func (h *Handler) DeleteProductHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Only DELETE is allowed")
		return
	}

	productID := strings.TrimPrefix(r.URL.Path, "/api/products/")
	productID = strings.TrimSuffix(productID, "/")

	if err := h.inventoryService.DeleteProduct(r.Context(), productID); err != nil {
		WriteError(w, http.StatusInternalServerError, "DELETE_FAILED", err.Error())
		return
	}

	WriteSuccess(w, http.StatusOK, "Product deleted successfully", nil)
}

// AddStockHandler handles adding stock
func (h *Handler) AddStockHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Only POST is allowed")
		return
	}

	productID := strings.TrimPrefix(r.URL.Path, "/api/products/")
	productID = strings.TrimPrefix(productID, "/stock/add")
	productID = strings.TrimSuffix(productID, "/")

	var req StockOperationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	if err := h.inventoryService.AddStock(r.Context(), productID, req.Quantity, req.Reference); err != nil {
		WriteError(w, http.StatusInternalServerError, "OPERATION_FAILED", err.Error())
		return
	}

	WriteSuccess(w, http.StatusOK, "Stock added successfully", nil)
}

// RemoveStockHandler handles removing stock
func (h *Handler) RemoveStockHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Only POST is allowed")
		return
	}

	productID := strings.TrimPrefix(r.URL.Path, "/api/products/")
	productID = strings.TrimPrefix(productID, "/stock/remove")
	productID = strings.TrimSuffix(productID, "/")

	var req StockOperationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	if err := h.inventoryService.RemoveStock(r.Context(), productID, req.Quantity, req.Reference); err != nil {
		WriteError(w, http.StatusInternalServerError, "OPERATION_FAILED", err.Error())
		return
	}

	WriteSuccess(w, http.StatusOK, "Stock removed successfully", nil)
}

// ReserveStockHandler handles reserving stock
func (h *Handler) ReserveStockHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Only POST is allowed")
		return
	}

	productID := strings.TrimPrefix(r.URL.Path, "/api/products/")
	productID = strings.TrimPrefix(productID, "/stock/reserve")
	productID = strings.TrimSuffix(productID, "/")

	var req StockOperationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	if err := h.inventoryService.ReserveStock(r.Context(), productID, req.Quantity, req.Reference); err != nil {
		WriteError(w, http.StatusInternalServerError, "OPERATION_FAILED", err.Error())
		return
	}

	WriteSuccess(w, http.StatusOK, "Stock reserved successfully", nil)
}

// UnreserveStockHandler handles unreserving stock
func (h *Handler) UnreserveStockHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Only POST is allowed")
		return
	}

	productID := strings.TrimPrefix(r.URL.Path, "/api/products/")
	productID = strings.TrimPrefix(productID, "/stock/unreserve")
	productID = strings.TrimSuffix(productID, "/")

	var req StockOperationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	if err := h.inventoryService.UnreserveStock(r.Context(), productID, req.Quantity, req.Reference); err != nil {
		WriteError(w, http.StatusInternalServerError, "OPERATION_FAILED", err.Error())
		return
	}

	WriteSuccess(w, http.StatusOK, "Stock unreserved successfully", nil)
}

// GetInventoryHandler handles retrieving inventory details
func (h *Handler) GetInventoryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Only GET is allowed")
		return
	}

	productID := strings.TrimPrefix(r.URL.Path, "/api/products/")
	productID = strings.TrimSuffix(productID, "/inventory")
	productID = strings.TrimSuffix(productID, "/")

	inventory, err := h.inventoryService.GetInventory(r.Context(), productID)
	if err != nil {
		WriteError(w, http.StatusNotFound, "NOT_FOUND", err.Error())
		return
	}

	WriteSuccess(w, http.StatusOK, "Inventory retrieved successfully", inventory)
}

// GetTransactionsHandler handles retrieving transaction history
func (h *Handler) GetTransactionsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Only GET is allowed")
		return
	}

	productID := strings.TrimPrefix(r.URL.Path, "/api/products/")
	productID = strings.TrimSuffix(productID, "/transactions")
	productID = strings.TrimSuffix(productID, "/")

	limit := 10
	offset := 0

	if l := r.URL.Query().Get("limit"); l != "" {
		if parsedLimit, err := strconv.Atoi(l); err == nil {
			limit = parsedLimit
		}
	}

	if o := r.URL.Query().Get("offset"); o != "" {
		if parsedOffset, err := strconv.Atoi(o); err == nil {
			offset = parsedOffset
		}
	}

	transactions, err := h.inventoryService.ListTransactions(r.Context(), productID, limit, offset)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "RETRIEVAL_FAILED", err.Error())
		return
	}

	WriteSuccess(w, http.StatusOK, "Transactions retrieved successfully", transactions)
}
