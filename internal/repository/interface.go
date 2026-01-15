package repository

import (
	"context"
	"github.com/bhnrathore/distributed-inventory-system/internal/domain"
)

// ProductRepository defines the interface for product data operations
type ProductRepository interface {
	Create(ctx context.Context, product *domain.Product) error
	GetByID(ctx context.Context, id string) (*domain.Product, error)
	GetBySKU(ctx context.Context, sku string) (*domain.Product, error)
	List(ctx context.Context, limit, offset int) ([]*domain.Product, error)
	Update(ctx context.Context, product *domain.Product) error
	Delete(ctx context.Context, id string) error
	Count(ctx context.Context) (int64, error)
}

// InventoryRepository defines the interface for inventory data operations
type InventoryRepository interface {
	Create(ctx context.Context, item *domain.InventoryItem) error
	GetByID(ctx context.Context, id string) (*domain.InventoryItem, error)
	GetByProductID(ctx context.Context, productID string) (*domain.InventoryItem, error)
	List(ctx context.Context, limit, offset int) ([]*domain.InventoryItem, error)
	Update(ctx context.Context, item *domain.InventoryItem) error
	Delete(ctx context.Context, id string) error
	UpdateQuantity(ctx context.Context, inventoryID string, quantityDelta, reservedDelta int64) error
}

// TransactionRepository defines the interface for transaction data operations
type TransactionRepository interface {
	Create(ctx context.Context, transaction *domain.Transaction) error
	GetByID(ctx context.Context, id string) (*domain.Transaction, error)
	GetByInventoryID(ctx context.Context, inventoryID string, limit, offset int) ([]*domain.Transaction, error)
	GetByProductID(ctx context.Context, productID string, limit, offset int) ([]*domain.Transaction, error)
	List(ctx context.Context, limit, offset int) ([]*domain.Transaction, error)
	Count(ctx context.Context) (int64, error)
}
