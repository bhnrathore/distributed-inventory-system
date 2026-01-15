package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/bhnrathore/distributed-inventory-system/internal/domain"
	"github.com/bhnrathore/distributed-inventory-system/internal/repository"
)

// InventoryService handles inventory business logic
type InventoryService struct {
	productRepo     repository.ProductRepository
	inventoryRepo   repository.InventoryRepository
	transactionRepo repository.TransactionRepository
}

// NewInventoryService creates a new InventoryService
func NewInventoryService(
	productRepo repository.ProductRepository,
	inventoryRepo repository.InventoryRepository,
	transactionRepo repository.TransactionRepository,
) *InventoryService {
	return &InventoryService{
		productRepo:     productRepo,
		inventoryRepo:   inventoryRepo,
		transactionRepo: transactionRepo,
	}
}

// CreateProduct creates a new product and initializes inventory
func (s *InventoryService) CreateProduct(ctx context.Context, product *domain.Product, location string, initialQuantity int64) error {
	if err := product.Validate(); err != nil {
		return fmt.Errorf("invalid product: %w", err)
	}

	// Create product
	if err := s.productRepo.Create(ctx, product); err != nil {
		return fmt.Errorf("failed to create product: %w", err)
	}

	// Create inventory item
	inventoryItem := &domain.InventoryItem{
		ProductID: product.ID,
		Quantity:  initialQuantity,
		Reserved:  0,
		Location:  location,
	}

	if err := s.inventoryRepo.Create(ctx, inventoryItem); err != nil {
		// Clean up product if inventory creation fails
		_ = s.productRepo.Delete(ctx, product.ID)
		return fmt.Errorf("failed to create inventory: %w", err)
	}

	// Record initial stock transaction
	if initialQuantity > 0 {
		transaction := &domain.Transaction{
			InventoryID: inventoryItem.ID,
			ProductID:   product.ID,
			Type:        "IN",
			Quantity:    initialQuantity,
			Reference:   "INITIAL_STOCK",
			Notes:       "Initial stock entry",
		}
		_ = s.transactionRepo.Create(ctx, transaction)
	}

	return nil
}

// GetProduct retrieves a product with its inventory details
func (s *InventoryService) GetProduct(ctx context.Context, productID string) (*domain.Product, *domain.InventoryItem, error) {
	product, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get product: %w", err)
	}

	inventory, err := s.inventoryRepo.GetByProductID(ctx, productID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get inventory: %w", err)
	}

	return product, inventory, nil
}

// ListProducts lists all products with pagination
func (s *InventoryService) ListProducts(ctx context.Context, limit, offset int) ([]*domain.Product, error) {
	products, err := s.productRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list products: %w", err)
	}
	return products, nil
}

// UpdateProduct updates product details
func (s *InventoryService) UpdateProduct(ctx context.Context, product *domain.Product) error {
	if err := product.Validate(); err != nil {
		return fmt.Errorf("invalid product: %w", err)
	}

	if err := s.productRepo.Update(ctx, product); err != nil {
		return fmt.Errorf("failed to update product: %w", err)
	}

	return nil
}

// AddStock adds stock to inventory
func (s *InventoryService) AddStock(ctx context.Context, productID string, quantity int64, reference string) error {
	if quantity <= 0 {
		return errors.New("quantity must be positive")
	}

	inventory, err := s.inventoryRepo.GetByProductID(ctx, productID)
	if err != nil {
		return fmt.Errorf("failed to get inventory: %w", err)
	}

	// Update quantity
	if err := s.inventoryRepo.UpdateQuantity(ctx, inventory.ID, quantity, 0); err != nil {
		return fmt.Errorf("failed to update quantity: %w", err)
	}

	// Record transaction
	transaction := &domain.Transaction{
		InventoryID: inventory.ID,
		ProductID:   productID,
		Type:        "IN",
		Quantity:    quantity,
		Reference:   reference,
		Notes:       "Stock addition",
	}

	if err := s.transactionRepo.Create(ctx, transaction); err != nil {
		return fmt.Errorf("failed to record transaction: %w", err)
	}

	return nil
}

// RemoveStock removes stock from inventory
func (s *InventoryService) RemoveStock(ctx context.Context, productID string, quantity int64, reference string) error {
	if quantity <= 0 {
		return errors.New("quantity must be positive")
	}

	inventory, err := s.inventoryRepo.GetByProductID(ctx, productID)
	if err != nil {
		return fmt.Errorf("failed to get inventory: %w", err)
	}

	// Check if enough stock is available
	if inventory.AvailableQuantity() < quantity {
		return errors.New("insufficient stock available")
	}

	// Update quantity
	if err := s.inventoryRepo.UpdateQuantity(ctx, inventory.ID, -quantity, 0); err != nil {
		return fmt.Errorf("failed to update quantity: %w", err)
	}

	// Record transaction
	transaction := &domain.Transaction{
		InventoryID: inventory.ID,
		ProductID:   productID,
		Type:        "OUT",
		Quantity:    quantity,
		Reference:   reference,
		Notes:       "Stock removal",
	}

	if err := s.transactionRepo.Create(ctx, transaction); err != nil {
		return fmt.Errorf("failed to record transaction: %w", err)
	}

	return nil
}

// ReserveStock reserves stock for an order
func (s *InventoryService) ReserveStock(ctx context.Context, productID string, quantity int64, reference string) error {
	if quantity <= 0 {
		return errors.New("quantity must be positive")
	}

	inventory, err := s.inventoryRepo.GetByProductID(ctx, productID)
	if err != nil {
		return fmt.Errorf("failed to get inventory: %w", err)
	}

	// Check if enough stock is available
	if inventory.AvailableQuantity() < quantity {
		return errors.New("insufficient stock available for reservation")
	}

	// Update reserved quantity
	if err := s.inventoryRepo.UpdateQuantity(ctx, inventory.ID, 0, quantity); err != nil {
		return fmt.Errorf("failed to reserve stock: %w", err)
	}

	// Record transaction
	transaction := &domain.Transaction{
		InventoryID: inventory.ID,
		ProductID:   productID,
		Type:        "RESERVE",
		Quantity:    quantity,
		Reference:   reference,
		Notes:       "Stock reservation",
	}

	if err := s.transactionRepo.Create(ctx, transaction); err != nil {
		return fmt.Errorf("failed to record transaction: %w", err)
	}

	return nil
}

// UnreserveStock releases reserved stock
func (s *InventoryService) UnreserveStock(ctx context.Context, productID string, quantity int64, reference string) error {
	if quantity <= 0 {
		return errors.New("quantity must be positive")
	}

	inventory, err := s.inventoryRepo.GetByProductID(ctx, productID)
	if err != nil {
		return fmt.Errorf("failed to get inventory: %w", err)
	}

	// Check if enough reserved stock exists
	if inventory.Reserved < quantity {
		return errors.New("insufficient reserved stock")
	}

	// Update reserved quantity
	if err := s.inventoryRepo.UpdateQuantity(ctx, inventory.ID, 0, -quantity); err != nil {
		return fmt.Errorf("failed to unreserve stock: %w", err)
	}

	// Record transaction
	transaction := &domain.Transaction{
		InventoryID: inventory.ID,
		ProductID:   productID,
		Type:        "UNRESERVE",
		Quantity:    quantity,
		Reference:   reference,
		Notes:       "Stock unreservation",
	}

	if err := s.transactionRepo.Create(ctx, transaction); err != nil {
		return fmt.Errorf("failed to record transaction: %w", err)
	}

	return nil
}

// GetInventory retrieves inventory details for a product
func (s *InventoryService) GetInventory(ctx context.Context, productID string) (*domain.InventoryItem, error) {
	inventory, err := s.inventoryRepo.GetByProductID(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to get inventory: %w", err)
	}
	return inventory, nil
}

// ListTransactions lists transactions for a product
func (s *InventoryService) ListTransactions(ctx context.Context, productID string, limit, offset int) ([]*domain.Transaction, error) {
	transactions, err := s.transactionRepo.GetByProductID(ctx, productID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list transactions: %w", err)
	}
	return transactions, nil
}

// DeleteProduct deletes a product and its inventory
func (s *InventoryService) DeleteProduct(ctx context.Context, productID string) error {
	// This will cascade delete inventory and transactions due to foreign keys
	if err := s.productRepo.Delete(ctx, productID); err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}
	return nil
}
