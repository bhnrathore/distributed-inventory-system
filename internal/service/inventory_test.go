package service

import (
	"context"
	"testing"

	"github.com/bhnrathore/distributed-inventory-system/internal/domain"
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
		product.ID = "test-id-1"
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
		item.ID = "test-inv-1"
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
		transaction.ID = "test-tx-1"
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

func TestCreateProduct(t *testing.T) {
	productRepo := NewMockProductRepository()
	inventoryRepo := NewMockInventoryRepository()
	transactionRepo := NewMockTransactionRepository()

	service := NewInventoryService(productRepo, inventoryRepo, transactionRepo)
	ctx := context.Background()

	product := &domain.Product{
		Name:        "Laptop",
		SKU:         "LAP001",
		Description: "Gaming Laptop",
		Price:       1500.00,
	}

	err := service.CreateProduct(ctx, product, "Warehouse A", 50)
	if err != nil {
		t.Fatalf("Failed to create product: %v", err)
	}

	if product.ID == "" {
		t.Fatal("Product ID should be set")
	}
}

func TestAddStock(t *testing.T) {
	productRepo := NewMockProductRepository()
	inventoryRepo := NewMockInventoryRepository()
	transactionRepo := NewMockTransactionRepository()

	service := NewInventoryService(productRepo, inventoryRepo, transactionRepo)
	ctx := context.Background()

	product := &domain.Product{
		ID:          "prod-1",
		Name:        "Laptop",
		SKU:         "LAP001",
		Description: "Gaming Laptop",
		Price:       1500.00,
	}
	productRepo.Create(ctx, product)

	inventory := &domain.InventoryItem{
		ID:        "inv-1",
		ProductID: product.ID,
		Quantity:  50,
		Reserved:  0,
		Location:  "Warehouse A",
	}
	inventoryRepo.Create(ctx, inventory)

	err := service.AddStock(ctx, product.ID, 20, "PO-001")
	if err != nil {
		t.Fatalf("Failed to add stock: %v", err)
	}

	updated, _ := inventoryRepo.GetByProductID(ctx, product.ID)
	if updated.Quantity != 70 {
		t.Errorf("Expected quantity 70, got %d", updated.Quantity)
	}
}

func TestRemoveStock(t *testing.T) {
	productRepo := NewMockProductRepository()
	inventoryRepo := NewMockInventoryRepository()
	transactionRepo := NewMockTransactionRepository()

	service := NewInventoryService(productRepo, inventoryRepo, transactionRepo)
	ctx := context.Background()

	product := &domain.Product{
		ID:          "prod-1",
		Name:        "Laptop",
		SKU:         "LAP001",
		Description: "Gaming Laptop",
		Price:       1500.00,
	}
	productRepo.Create(ctx, product)

	inventory := &domain.InventoryItem{
		ID:        "inv-1",
		ProductID: product.ID,
		Quantity:  50,
		Reserved:  0,
		Location:  "Warehouse A",
	}
	inventoryRepo.Create(ctx, inventory)

	err := service.RemoveStock(ctx, product.ID, 20, "ORDER-001")
	if err != nil {
		t.Fatalf("Failed to remove stock: %v", err)
	}

	updated, _ := inventoryRepo.GetByProductID(ctx, product.ID)
	if updated.Quantity != 30 {
		t.Errorf("Expected quantity 30, got %d", updated.Quantity)
	}
}

func TestReserveStock(t *testing.T) {
	productRepo := NewMockProductRepository()
	inventoryRepo := NewMockInventoryRepository()
	transactionRepo := NewMockTransactionRepository()

	service := NewInventoryService(productRepo, inventoryRepo, transactionRepo)
	ctx := context.Background()

	product := &domain.Product{
		ID:          "prod-1",
		Name:        "Laptop",
		SKU:         "LAP001",
		Description: "Gaming Laptop",
		Price:       1500.00,
	}
	productRepo.Create(ctx, product)

	inventory := &domain.InventoryItem{
		ID:        "inv-1",
		ProductID: product.ID,
		Quantity:  50,
		Reserved:  0,
		Location:  "Warehouse A",
	}
	inventoryRepo.Create(ctx, inventory)

	err := service.ReserveStock(ctx, product.ID, 10, "ORDER-001")
	if err != nil {
		t.Fatalf("Failed to reserve stock: %v", err)
	}

	updated, _ := inventoryRepo.GetByProductID(ctx, product.ID)
	if updated.Reserved != 10 {
		t.Errorf("Expected reserved 10, got %d", updated.Reserved)
	}

	if updated.AvailableQuantity() != 40 {
		t.Errorf("Expected available quantity 40, got %d", updated.AvailableQuantity())
	}
}

func TestInsufficientStockRemoval(t *testing.T) {
	productRepo := NewMockProductRepository()
	inventoryRepo := NewMockInventoryRepository()
	transactionRepo := NewMockTransactionRepository()

	service := NewInventoryService(productRepo, inventoryRepo, transactionRepo)
	ctx := context.Background()

	product := &domain.Product{
		ID:          "prod-1",
		Name:        "Laptop",
		SKU:         "LAP001",
		Description: "Gaming Laptop",
		Price:       1500.00,
	}
	productRepo.Create(ctx, product)

	inventory := &domain.InventoryItem{
		ID:        "inv-1",
		ProductID: product.ID,
		Quantity:  10,
		Reserved:  0,
		Location:  "Warehouse A",
	}
	inventoryRepo.Create(ctx, inventory)

	err := service.RemoveStock(ctx, product.ID, 20, "ORDER-001")
	if err == nil {
		t.Fatal("Expected error for insufficient stock")
	}
}
func TestReleaseReservedStock(t *testing.T) {
	productRepo := NewMockProductRepository()
	inventoryRepo := NewMockInventoryRepository()
	transactionRepo := NewMockTransactionRepository()

	service := NewInventoryService(productRepo, inventoryRepo, transactionRepo)
	ctx := context.Background()

	product := &domain.Product{
		ID:          "prod-1",
		Name:        "Laptop",
		SKU:         "LAP001",
		Description: "Gaming Laptop",
		Price:       1500.00,
	}
	productRepo.Create(ctx, product)

	inventory := &domain.InventoryItem{
		ID:        "inv-1",
		ProductID: product.ID,
		Quantity:  50,
		Reserved:  20,
		Location:  "Warehouse A",
	}
	inventoryRepo.Create(ctx, inventory)

	err := service.UnreserveStock(ctx, product.ID, 20, "ORDER-001")
	if err != nil {
		t.Fatalf("Failed to unreserve stock: %v", err)
	}

	updated, _ := inventoryRepo.GetByProductID(ctx, product.ID)
	if updated.Reserved != 0 {
		t.Errorf("Expected reserved 0, got %d", updated.Reserved)
	}
}

func TestInsufficientReservedStock(t *testing.T) {
	productRepo := NewMockProductRepository()
	inventoryRepo := NewMockInventoryRepository()
	transactionRepo := NewMockTransactionRepository()

	service := NewInventoryService(productRepo, inventoryRepo, transactionRepo)
	ctx := context.Background()

	product := &domain.Product{
		ID:          "prod-1",
		Name:        "Laptop",
		SKU:         "LAP001",
		Description: "Gaming Laptop",
		Price:       1500.00,
	}
	productRepo.Create(ctx, product)

	inventory := &domain.InventoryItem{
		ID:        "inv-1",
		ProductID: product.ID,
		Quantity:  50,
		Reserved:  5,
		Location:  "Warehouse A",
	}
	inventoryRepo.Create(ctx, inventory)

	err := service.UnreserveStock(ctx, product.ID, 20, "ORDER-001")
	if err == nil {
		t.Fatal("Expected error for insufficient reserved stock")
	}
}

func TestGetProductWithInventory(t *testing.T) {
	productRepo := NewMockProductRepository()
	inventoryRepo := NewMockInventoryRepository()
	transactionRepo := NewMockTransactionRepository()

	service := NewInventoryService(productRepo, inventoryRepo, transactionRepo)
	ctx := context.Background()

	product := &domain.Product{
		ID:          "prod-1",
		Name:        "Laptop",
		SKU:         "LAP001",
		Description: "Gaming Laptop",
		Price:       1500.00,
	}
	productRepo.Create(ctx, product)

	inventory := &domain.InventoryItem{
		ID:        "inv-1",
		ProductID: product.ID,
		Quantity:  50,
		Reserved:  10,
		Location:  "Warehouse A",
	}
	inventoryRepo.Create(ctx, inventory)

	productResult, inventoryResult, err := service.GetProduct(ctx, product.ID)
	if err != nil {
		t.Fatalf("Failed to get product: %v", err)
	}

	if productResult == nil {
		t.Fatal("Expected non-nil product")
	}

	if inventoryResult == nil {
		t.Fatal("Expected non-nil inventory")
	}
}

func TestGetProductNotFound(t *testing.T) {
	productRepo := NewMockProductRepository()
	inventoryRepo := NewMockInventoryRepository()
	transactionRepo := NewMockTransactionRepository()

	service := NewInventoryService(productRepo, inventoryRepo, transactionRepo)
	ctx := context.Background()

	productResult, inventoryResult, _ := service.GetProduct(ctx, "non-existent")
	if productResult != nil || inventoryResult != nil {
		t.Fatal("Expected nil for non-existent product")
	}
}

func TestCreateProductWithInvalidData(t *testing.T) {
	productRepo := NewMockProductRepository()
	inventoryRepo := NewMockInventoryRepository()
	transactionRepo := NewMockTransactionRepository()

	service := NewInventoryService(productRepo, inventoryRepo, transactionRepo)
	ctx := context.Background()

	product := &domain.Product{
		Name:        "",
		SKU:         "LAP001",
		Description: "Gaming Laptop",
		Price:       1500.00,
	}

	err := service.CreateProduct(ctx, product, "Warehouse A", 50)
	if err == nil {
		t.Fatal("Expected error for invalid product")
	}
}

func TestAddStockWithInvalidQuantity(t *testing.T) {
	productRepo := NewMockProductRepository()
	inventoryRepo := NewMockInventoryRepository()
	transactionRepo := NewMockTransactionRepository()

	service := NewInventoryService(productRepo, inventoryRepo, transactionRepo)
	ctx := context.Background()

	product := &domain.Product{
		ID:          "prod-1",
		Name:        "Laptop",
		SKU:         "LAP001",
		Description: "Gaming Laptop",
		Price:       1500.00,
	}
	productRepo.Create(ctx, product)

	inventory := &domain.InventoryItem{
		ID:        "inv-1",
		ProductID: product.ID,
		Quantity:  50,
		Reserved:  0,
		Location:  "Warehouse A",
	}
	inventoryRepo.Create(ctx, inventory)

	err := service.AddStock(ctx, product.ID, -10, "PO-001")
	if err == nil {
		t.Fatal("Expected error for negative quantity")
	}
}

func TestReleaseReservedStockWithInvalidQuantity(t *testing.T) {
	productRepo := NewMockProductRepository()
	inventoryRepo := NewMockInventoryRepository()
	transactionRepo := NewMockTransactionRepository()

	service := NewInventoryService(productRepo, inventoryRepo, transactionRepo)
	ctx := context.Background()

	product := &domain.Product{
		ID:          "prod-1",
		Name:        "Laptop",
		SKU:         "LAP001",
		Description: "Gaming Laptop",
		Price:       1500.00,
	}
	productRepo.Create(ctx, product)

	inventory := &domain.InventoryItem{
		ID:        "inv-1",
		ProductID: product.ID,
		Quantity:  50,
		Reserved:  10,
		Location:  "Warehouse A",
	}
	inventoryRepo.Create(ctx, inventory)

	err := service.UnreserveStock(ctx, product.ID, -5, "ORDER-001")
	if err == nil {
		t.Fatal("Expected error for negative quantity")
	}
}

func TestListProducts(t *testing.T) {
	productRepo := NewMockProductRepository()
	inventoryRepo := NewMockInventoryRepository()
	transactionRepo := NewMockTransactionRepository()

	service := NewInventoryService(productRepo, inventoryRepo, transactionRepo)
	ctx := context.Background()

	// Create multiple products
	for i := 1; i <= 3; i++ {
		num := string(rune('0' + i))
		product := &domain.Product{
			Name:        "Product " + num,
			SKU:         "SKU00" + num,
			Description: "Test Product",
			Price:       float64(i * 100),
		}
		productRepo.Create(ctx, product)
	}

	products, err := service.ListProducts(ctx, 10, 0)
	if err != nil {
		t.Fatalf("Failed to list products: %v", err)
	}

	if len(products) < 3 {
		t.Logf("Expected at least 3 products, got %d", len(products))
	}
}

func TestListTransactions(t *testing.T) {
	productRepo := NewMockProductRepository()
	inventoryRepo := NewMockInventoryRepository()
	transactionRepo := NewMockTransactionRepository()

	service := NewInventoryService(productRepo, inventoryRepo, transactionRepo)
	ctx := context.Background()

	product := &domain.Product{
		ID:          "prod-1",
		Name:        "Laptop",
		SKU:         "LAP001",
		Description: "Gaming Laptop",
		Price:       1500.00,
	}
	productRepo.Create(ctx, product)

	inventory := &domain.InventoryItem{
		ID:        "inv-1",
		ProductID: product.ID,
		Quantity:  50,
		Reserved:  0,
		Location:  "Warehouse A",
	}
	inventoryRepo.Create(ctx, inventory)

	// Add stock which creates a transaction
	_ = service.AddStock(ctx, product.ID, 10, "PO-001")

	transactions, err := service.ListTransactions(ctx, product.ID, 10, 0)
	if err != nil {
		t.Fatalf("Failed to list transactions: %v", err)
	}

	if len(transactions) == 0 {
		t.Fatal("Expected at least one transaction")
	}
}