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

// PostgresTransactionRepository implements TransactionRepository using PostgreSQL
type PostgresTransactionRepository struct {
	db *sql.DB
}

// NewPostgresTransactionRepository creates a new PostgresTransactionRepository
func NewPostgresTransactionRepository(db *sql.DB) *PostgresTransactionRepository {
	return &PostgresTransactionRepository{db: db}
}

// Create inserts a new transaction
func (r *PostgresTransactionRepository) Create(ctx context.Context, transaction *domain.Transaction) error {
	if err := transaction.Validate(); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	transaction.ID = uuid.New().String()
	transaction.CreatedAt = time.Now()

	query := `
		INSERT INTO transactions (id, inventory_id, product_id, type, quantity, reference, notes, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.db.ExecContext(ctx, query,
		transaction.ID, transaction.InventoryID, transaction.ProductID, transaction.Type,
		transaction.Quantity, transaction.Reference, transaction.Notes, transaction.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}

	return nil
}

// GetByID retrieves a transaction by ID
func (r *PostgresTransactionRepository) GetByID(ctx context.Context, id string) (*domain.Transaction, error) {
	query := `
		SELECT id, inventory_id, product_id, type, quantity, reference, notes, created_at
		FROM transactions WHERE id = $1
	`

	transaction := &domain.Transaction{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&transaction.ID, &transaction.InventoryID, &transaction.ProductID, &transaction.Type,
		&transaction.Quantity, &transaction.Reference, &transaction.Notes, &transaction.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("transaction not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}

	return transaction, nil
}

// GetByInventoryID retrieves transactions for a specific inventory item
func (r *PostgresTransactionRepository) GetByInventoryID(ctx context.Context, inventoryID string, limit, offset int) ([]*domain.Transaction, error) {
	query := `
		SELECT id, inventory_id, product_id, type, quantity, reference, notes, created_at
		FROM transactions
		WHERE inventory_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, inventoryID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list transactions: %w", err)
	}
	defer rows.Close()

	var transactions []*domain.Transaction
	for rows.Next() {
		transaction := &domain.Transaction{}
		if err := rows.Scan(
			&transaction.ID, &transaction.InventoryID, &transaction.ProductID, &transaction.Type,
			&transaction.Quantity, &transaction.Reference, &transaction.Notes, &transaction.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan transaction: %w", err)
		}
		transactions = append(transactions, transaction)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating transactions: %w", err)
	}

	return transactions, nil
}

// GetByProductID retrieves transactions for a specific product
func (r *PostgresTransactionRepository) GetByProductID(ctx context.Context, productID string, limit, offset int) ([]*domain.Transaction, error) {
	query := `
		SELECT id, inventory_id, product_id, type, quantity, reference, notes, created_at
		FROM transactions
		WHERE product_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, productID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list transactions: %w", err)
	}
	defer rows.Close()

	var transactions []*domain.Transaction
	for rows.Next() {
		transaction := &domain.Transaction{}
		if err := rows.Scan(
			&transaction.ID, &transaction.InventoryID, &transaction.ProductID, &transaction.Type,
			&transaction.Quantity, &transaction.Reference, &transaction.Notes, &transaction.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan transaction: %w", err)
		}
		transactions = append(transactions, transaction)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating transactions: %w", err)
	}

	return transactions, nil
}

// List retrieves a paginated list of transactions
func (r *PostgresTransactionRepository) List(ctx context.Context, limit, offset int) ([]*domain.Transaction, error) {
	query := `
		SELECT id, inventory_id, product_id, type, quantity, reference, notes, created_at
		FROM transactions
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list transactions: %w", err)
	}
	defer rows.Close()

	var transactions []*domain.Transaction
	for rows.Next() {
		transaction := &domain.Transaction{}
		if err := rows.Scan(
			&transaction.ID, &transaction.InventoryID, &transaction.ProductID, &transaction.Type,
			&transaction.Quantity, &transaction.Reference, &transaction.Notes, &transaction.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan transaction: %w", err)
		}
		transactions = append(transactions, transaction)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating transactions: %w", err)
	}

	return transactions, nil
}

// Count returns the total number of transactions
func (r *PostgresTransactionRepository) Count(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM transactions`

	var count int64
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count transactions: %w", err)
	}

	return count, nil
}
