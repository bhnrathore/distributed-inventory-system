package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bhnrathore/distributed-inventory-system/internal/domain"
	"github.com/bhnrathore/distributed-inventory-system/internal/service"
)

// MockProductRepository implements ProductRepository interface for testing
type MockProductRepository struct {
	products map[string]*domain.Product
}

func NewMockProductRepository() *MockProductRepository {
	return &MockProductRepository{
		products: make(map[string]*domain.Product),
	}
}

func (m *MockProductRepository) Create(ctx context.Context, product *domain.Product) error {
	if product.ID == "" {
		product.ID = "test-id-" + product.SKU
	}
	m.products[product.ID] = product
	return nil
}

func (m *MockProductRepository) GetByID(ctx context.Context, id string) (*domain.Product, error) {
	if p, ok := m.products[id]; ok {
		return p, nil
	}
	return nil, nil
}

func (m *MockProductRepository) GetBySKU(ctx context.Context, sku string) (*domain.Product, error) {
	for _, p := range m.products {
		if p.SKU == sku {
			return p, nil
		}
	}
	return nil, nil
}

func (m *MockProductRepository) List(ctx context.Context, limit, offset int) ([]*domain.Product, error) {
	var products []*domain.Product
	for _, p := range m.products {
		products = append(products, p)
	}
	return products, nil
}

func (m *MockProductRepository) Update(ctx context.Context, product *domain.Product) error {
	m.products[product.ID] = product
	return nil
}

func (m *MockProductRepository) Delete(ctx context.Context, id string) error {
	delete(m.products, id)
	return nil
}

func (m *MockProductRepository) Count(ctx context.Context) (int64, error) {
	return int64(len(m.products)), nil
}

// MockInventoryRepository implements InventoryRepository interface for testing
type MockInventoryRepository struct {
	items map[string]*domain.InventoryItem
}

func NewMockInventoryRepository() *MockInventoryRepository {
	return &MockInventoryRepository{
		items: make(map[string]*domain.InventoryItem),
	}
}

func (m *MockInventoryRepository) Create(ctx context.Context, item *domain.InventoryItem) error {
	if item.ID == "" {
		item.ID = "inv-" + item.ProductID
	}
	m.items[item.ID] = item
	return nil
}

func (m *MockInventoryRepository) GetByID(ctx context.Context, id string) (*domain.InventoryItem, error) {
	if i, ok := m.items[id]; ok {
		return i, nil
	}
	return nil, nil
}

func (m *MockInventoryRepository) GetByProductID(ctx context.Context, productID string) (*domain.InventoryItem, error) {
	for _, i := range m.items {
		if i.ProductID == productID {
			return i, nil
		}
	}
	return nil, nil
}

func (m *MockInventoryRepository) List(ctx context.Context, limit, offset int) ([]*domain.InventoryItem, error) {
	var items []*domain.InventoryItem
	for _, i := range m.items {
		items = append(items, i)
	}
	return items, nil
}

func (m *MockInventoryRepository) Update(ctx context.Context, item *domain.InventoryItem) error {
	m.items[item.ID] = item
	return nil
}

func (m *MockInventoryRepository) Delete(ctx context.Context, id string) error {
	delete(m.items, id)
	return nil
}

func (m *MockInventoryRepository) UpdateQuantity(ctx context.Context, inventoryID string, quantityDelta, reservedDelta int64) error {
	if i, ok := m.items[inventoryID]; ok {
		i.Quantity += quantityDelta
		i.Reserved += reservedDelta
		return nil
	}
	return nil
}

// MockTransactionRepository implements TransactionRepository interface for testing
type MockTransactionRepository struct {
	transactions map[string]*domain.Transaction
}

func NewMockTransactionRepository() *MockTransactionRepository {
	return &MockTransactionRepository{
		transactions: make(map[string]*domain.Transaction),
	}
}

func (m *MockTransactionRepository) Create(ctx context.Context, transaction *domain.Transaction) error {
	if transaction.ID == "" {
		transaction.ID = "tx-" + transaction.Reference
	}
	m.transactions[transaction.ID] = transaction
	return nil
}

func (m *MockTransactionRepository) GetByID(ctx context.Context, id string) (*domain.Transaction, error) {
	if t, ok := m.transactions[id]; ok {
		return t, nil
	}
	return nil, nil
}

func (m *MockTransactionRepository) GetByInventoryID(ctx context.Context, inventoryID string, limit, offset int) ([]*domain.Transaction, error) {
	var txs []*domain.Transaction
	for _, t := range m.transactions {
		if t.InventoryID == inventoryID {
			txs = append(txs, t)
		}
	}
	return txs, nil
}

func (m *MockTransactionRepository) GetByProductID(ctx context.Context, productID string, limit, offset int) ([]*domain.Transaction, error) {
	var txs []*domain.Transaction
	for _, t := range m.transactions {
		if t.ProductID == productID {
			txs = append(txs, t)
		}
	}
	return txs, nil
}

func (m *MockTransactionRepository) List(ctx context.Context, limit, offset int) ([]*domain.Transaction, error) {
	var txs []*domain.Transaction
	for _, t := range m.transactions {
		txs = append(txs, t)
	}
	return txs, nil
}

func (m *MockTransactionRepository) Count(ctx context.Context) (int64, error) {
	return int64(len(m.transactions)), nil
}

// Tests

func TestHealthHandler(t *testing.T) {
	productRepo := NewMockProductRepository()
	inventoryRepo := NewMockInventoryRepository()
	transactionRepo := NewMockTransactionRepository()
	invService := service.NewInventoryService(productRepo, inventoryRepo, transactionRepo)
	handler := NewHandler(invService)

	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.HealthHandler(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Errorf("failed to decode response: %v", err)
	}
}

func TestHealthHandlerMethodNotAllowed(t *testing.T) {
	productRepo := NewMockProductRepository()
	inventoryRepo := NewMockInventoryRepository()
	transactionRepo := NewMockTransactionRepository()
	invService := service.NewInventoryService(productRepo, inventoryRepo, transactionRepo)
	handler := NewHandler(invService)

	req, err := http.NewRequest("POST", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.HealthHandler(rr, req)

	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusMethodNotAllowed)
	}
}

func TestCreateProductHandler(t *testing.T) {
	productRepo := NewMockProductRepository()
	inventoryRepo := NewMockInventoryRepository()
	transactionRepo := NewMockTransactionRepository()
	invService := service.NewInventoryService(productRepo, inventoryRepo, transactionRepo)
	handler := NewHandler(invService)

	reqBody := CreateProductRequest{
		Name:            "Laptop",
		Description:     "Gaming Laptop",
		SKU:             "LAP001",
		Price:           1500.00,
		Location:        "Warehouse A",
		InitialQuantity: 50,
	}

	body, _ := json.Marshal(reqBody)
	req, err := http.NewRequest("POST", "/products", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.CreateProductHandler(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}
}

func TestCreateProductHandlerInvalidRequest(t *testing.T) {
	productRepo := NewMockProductRepository()
	inventoryRepo := NewMockInventoryRepository()
	transactionRepo := NewMockTransactionRepository()
	invService := service.NewInventoryService(productRepo, inventoryRepo, transactionRepo)
	handler := NewHandler(invService)

	req, err := http.NewRequest("POST", "/products", bytes.NewBuffer([]byte("invalid json")))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.CreateProductHandler(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}

func TestCreateProductHandlerMethodNotAllowed(t *testing.T) {
	productRepo := NewMockProductRepository()
	inventoryRepo := NewMockInventoryRepository()
	transactionRepo := NewMockTransactionRepository()
	invService := service.NewInventoryService(productRepo, inventoryRepo, transactionRepo)
	handler := NewHandler(invService)

	req, err := http.NewRequest("GET", "/products", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.CreateProductHandler(rr, req)

	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusMethodNotAllowed)
	}
}
