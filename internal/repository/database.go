package repository

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// Database handles database connection and initialization
type Database struct {
	conn *sql.DB
}

// NewDatabase creates a new database connection
func NewDatabase(dsn string) (*Database, error) {
	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Verify the connection
	if err := conn.PingContext(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Set connection pool settings
	conn.SetMaxOpenConns(25)
	conn.SetMaxIdleConns(5)

	return &Database{conn: conn}, nil
}

// GetConnection returns the database connection
func (d *Database) GetConnection() *sql.DB {
	return d.conn
}

// Close closes the database connection
func (d *Database) Close() error {
	return d.conn.Close()
}

// InitSchema creates the database schema
func (d *Database) InitSchema(ctx context.Context) error {
	schema := `
	CREATE TABLE IF NOT EXISTS products (
		id VARCHAR(36) PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		description TEXT,
		sku VARCHAR(100) UNIQUE NOT NULL,
		price NUMERIC(10, 2) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS inventory (
		id VARCHAR(36) PRIMARY KEY,
		product_id VARCHAR(36) NOT NULL UNIQUE,
		quantity BIGINT NOT NULL DEFAULT 0,
		reserved BIGINT NOT NULL DEFAULT 0,
		location VARCHAR(255) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS transactions (
		id VARCHAR(36) PRIMARY KEY,
		inventory_id VARCHAR(36) NOT NULL,
		product_id VARCHAR(36) NOT NULL,
		type VARCHAR(20) NOT NULL,
		quantity BIGINT NOT NULL,
		reference VARCHAR(255),
		notes TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (inventory_id) REFERENCES inventory(id) ON DELETE CASCADE,
		FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE
	);

	CREATE INDEX IF NOT EXISTS idx_products_sku ON products(sku);
	CREATE INDEX IF NOT EXISTS idx_inventory_product_id ON inventory(product_id);
	CREATE INDEX IF NOT EXISTS idx_transactions_inventory_id ON transactions(inventory_id);
	CREATE INDEX IF NOT EXISTS idx_transactions_product_id ON transactions(product_id);
	CREATE INDEX IF NOT EXISTS idx_transactions_created_at ON transactions(created_at DESC);
	`

	_, err := d.conn.ExecContext(ctx, schema)
	if err != nil {
		return fmt.Errorf("failed to create schema: %w", err)
	}

	return nil
}
