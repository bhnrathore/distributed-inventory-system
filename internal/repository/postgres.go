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

// PostgresProductRepository implements ProductRepository using PostgreSQL
type PostgresProductRepository struct {
	db *sql.DB
}

// NewPostgresProductRepository creates a new PostgresProductRepository
func NewPostgresProductRepository(db *sql.DB) *PostgresProductRepository {
	return &PostgresProductRepository{db: db}
}

// Create inserts a new product
func (r *PostgresProductRepository) Create(ctx context.Context, product *domain.Product) error {
	if err := product.Validate(); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	product.ID = uuid.New().String()
	now := time.Now()
	product.CreatedAt = now
	product.UpdatedAt = now

	query := `
		INSERT INTO products (id, name, description, sku, price, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.db.ExecContext(ctx, query,
		product.ID, product.Name, product.Description, product.SKU, product.Price,
		product.CreatedAt, product.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create product: %w", err)
	}

	return nil
}

// GetByID retrieves a product by ID
func (r *PostgresProductRepository) GetByID(ctx context.Context, id string) (*domain.Product, error) {
	query := `
		SELECT id, name, description, sku, price, created_at, updated_at
		FROM products WHERE id = $1
	`

	product := &domain.Product{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&product.ID, &product.Name, &product.Description, &product.SKU,
		&product.Price, &product.CreatedAt, &product.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("product not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	return product, nil
}

// GetBySKU retrieves a product by SKU
func (r *PostgresProductRepository) GetBySKU(ctx context.Context, sku string) (*domain.Product, error) {
	query := `
		SELECT id, name, description, sku, price, created_at, updated_at
		FROM products WHERE sku = $1
	`

	product := &domain.Product{}
	err := r.db.QueryRowContext(ctx, query, sku).Scan(
		&product.ID, &product.Name, &product.Description, &product.SKU,
		&product.Price, &product.CreatedAt, &product.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("product not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	return product, nil
}

// List retrieves a paginated list of products
func (r *PostgresProductRepository) List(ctx context.Context, limit, offset int) ([]*domain.Product, error) {
	query := `
		SELECT id, name, description, sku, price, created_at, updated_at
		FROM products
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list products: %w", err)
	}
	defer rows.Close()

	var products []*domain.Product
	for rows.Next() {
		product := &domain.Product{}
		if err := rows.Scan(
			&product.ID, &product.Name, &product.Description, &product.SKU,
			&product.Price, &product.CreatedAt, &product.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan product: %w", err)
		}
		products = append(products, product)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating products: %w", err)
	}

	return products, nil
}

// Update updates an existing product
func (r *PostgresProductRepository) Update(ctx context.Context, product *domain.Product) error {
	if err := product.Validate(); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	product.UpdatedAt = time.Now()

	query := `
		UPDATE products
		SET name = $1, description = $2, sku = $3, price = $4, updated_at = $5
		WHERE id = $6
	`

	result, err := r.db.ExecContext(ctx, query,
		product.Name, product.Description, product.SKU, product.Price,
		product.UpdatedAt, product.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update product: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rows == 0 {
		return errors.New("product not found")
	}

	return nil
}

// Delete deletes a product
func (r *PostgresProductRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM products WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rows == 0 {
		return errors.New("product not found")
	}

	return nil
}

// Count returns the total number of products
func (r *PostgresProductRepository) Count(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM products`

	var count int64
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count products: %w", err)
	}

	return count, nil
}
