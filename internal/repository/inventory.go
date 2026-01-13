package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/big1111111/distributed-inventory-system/internal/domain"
	"github.com/google/uuid"
)

// PostgresInventoryRepository implements InventoryRepository using PostgreSQL
type PostgresInventoryRepository struct {
	db *sql.DB
}

// NewPostgresInventoryRepository creates a new PostgresInventoryRepository
func NewPostgresInventoryRepository(db *sql.DB) *PostgresInventoryRepository {
	return &PostgresInventoryRepository{db: db}
}

// Create inserts a new inventory item
func (r *PostgresInventoryRepository) Create(ctx context.Context, item *domain.InventoryItem) error {
	if err := item.Validate(); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	item.ID = uuid.New().String()
	now := time.Now()
	item.CreatedAt = now
	item.UpdatedAt = now

	query := `
		INSERT INTO inventory (id, product_id, quantity, reserved, location, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.db.ExecContext(ctx, query,
		item.ID, item.ProductID, item.Quantity, item.Reserved, item.Location,
		item.CreatedAt, item.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create inventory item: %w", err)
	}

	return nil
}

// GetByID retrieves an inventory item by ID
func (r *PostgresInventoryRepository) GetByID(ctx context.Context, id string) (*domain.InventoryItem, error) {
	query := `
		SELECT id, product_id, quantity, reserved, location, created_at, updated_at
		FROM inventory WHERE id = $1
	`

	item := &domain.InventoryItem{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&item.ID, &item.ProductID, &item.Quantity, &item.Reserved, &item.Location,
		&item.CreatedAt, &item.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("inventory item not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get inventory item: %w", err)
	}

	return item, nil
}

// GetByProductID retrieves inventory for a specific product
func (r *PostgresInventoryRepository) GetByProductID(ctx context.Context, productID string) (*domain.InventoryItem, error) {
	query := `
		SELECT id, product_id, quantity, reserved, location, created_at, updated_at
		FROM inventory WHERE product_id = $1
	`

	item := &domain.InventoryItem{}
	err := r.db.QueryRowContext(ctx, query, productID).Scan(
		&item.ID, &item.ProductID, &item.Quantity, &item.Reserved, &item.Location,
		&item.CreatedAt, &item.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("inventory item not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get inventory item: %w", err)
	}

	return item, nil
}

// List retrieves a paginated list of inventory items
func (r *PostgresInventoryRepository) List(ctx context.Context, limit, offset int) ([]*domain.InventoryItem, error) {
	query := `
		SELECT id, product_id, quantity, reserved, location, created_at, updated_at
		FROM inventory
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list inventory items: %w", err)
	}
	defer rows.Close()

	var items []*domain.InventoryItem
	for rows.Next() {
		item := &domain.InventoryItem{}
		if err := rows.Scan(
			&item.ID, &item.ProductID, &item.Quantity, &item.Reserved, &item.Location,
			&item.CreatedAt, &item.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan inventory item: %w", err)
		}
		items = append(items, item)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating inventory items: %w", err)
	}

	return items, nil
}

// Update updates an existing inventory item
func (r *PostgresInventoryRepository) Update(ctx context.Context, item *domain.InventoryItem) error {
	if err := item.Validate(); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	item.UpdatedAt = time.Now()

	query := `
		UPDATE inventory
		SET quantity = $1, reserved = $2, location = $3, updated_at = $4
		WHERE id = $5
	`

	result, err := r.db.ExecContext(ctx, query,
		item.Quantity, item.Reserved, item.Location, item.UpdatedAt, item.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update inventory item: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rows == 0 {
		return errors.New("inventory item not found")
	}

	return nil
}

// Delete deletes an inventory item
func (r *PostgresInventoryRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM inventory WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete inventory item: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rows == 0 {
		return errors.New("inventory item not found")
	}

	return nil
}

// UpdateQuantity updates the quantity and reserved quantities atomically
func (r *PostgresInventoryRepository) UpdateQuantity(ctx context.Context, inventoryID string, quantityDelta, reservedDelta int64) error {
	query := `
		UPDATE inventory
		SET quantity = quantity + $1, reserved = reserved + $2, updated_at = $3
		WHERE id = $4 AND (quantity + $1) >= 0 AND (reserved + $2) >= 0 AND (quantity + $1 - reserved - $2) >= 0
	`

	result, err := r.db.ExecContext(ctx, query, quantityDelta, reservedDelta, time.Now(), inventoryID)
	if err != nil {
		return fmt.Errorf("failed to update quantity: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rows == 0 {
		return errors.New("quantity update failed: invalid operation or item not found")
	}

	return nil
}
