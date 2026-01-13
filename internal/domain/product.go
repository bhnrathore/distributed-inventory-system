package domain

import (
	"errors"
	"time"
)

// Product represents a product in the inventory system
type Product struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	SKU         string    `json:"sku"`
	Price       float64   `json:"price"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Validate checks if the product data is valid
func (p *Product) Validate() error {
	if p.Name == "" {
		return errors.New("product name cannot be empty")
	}
	if p.SKU == "" {
		return errors.New("product SKU cannot be empty")
	}
	if p.Price < 0 {
		return errors.New("product price cannot be negative")
	}
	return nil
}

// InventoryItem represents the stock level for a product
type InventoryItem struct {
	ID        string    `json:"id"`
	ProductID string    `json:"product_id"`
	Quantity  int64     `json:"quantity"`
	Reserved  int64     `json:"reserved"`
	Location  string    `json:"location"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// AvailableQuantity returns the available (non-reserved) quantity
func (i *InventoryItem) AvailableQuantity() int64 {
	available := i.Quantity - i.Reserved
	if available < 0 {
		return 0
	}
	return available
}

// Validate checks if the inventory item data is valid
func (i *InventoryItem) Validate() error {
	if i.ProductID == "" {
		return errors.New("product_id cannot be empty")
	}
	if i.Quantity < 0 {
		return errors.New("quantity cannot be negative")
	}
	if i.Reserved < 0 {
		return errors.New("reserved quantity cannot be negative")
	}
	if i.Reserved > i.Quantity {
		return errors.New("reserved quantity cannot exceed total quantity")
	}
	if i.Location == "" {
		return errors.New("location cannot be empty")
	}
	return nil
}

// Transaction represents a stock movement transaction
type Transaction struct {
	ID          string    `json:"id"`
	InventoryID string    `json:"inventory_id"`
	ProductID   string    `json:"product_id"`
	Type        string    `json:"type"` // "IN", "OUT", "RETURN", "RESERVE", "UNRESERVE"
	Quantity    int64     `json:"quantity"`
	Reference   string    `json:"reference"` // e.g., order ID, return ID
	Notes       string    `json:"notes"`
	CreatedAt   time.Time `json:"created_at"`
}

// Validate checks if the transaction data is valid
func (t *Transaction) Validate() error {
	if t.InventoryID == "" {
		return errors.New("inventory_id cannot be empty")
	}
	if t.ProductID == "" {
		return errors.New("product_id cannot be empty")
	}
	if t.Quantity <= 0 {
		return errors.New("quantity must be positive")
	}
	validTypes := map[string]bool{
		"IN":        true,
		"OUT":       true,
		"RETURN":    true,
		"RESERVE":   true,
		"UNRESERVE": true,
	}
	if !validTypes[t.Type] {
		return errors.New("invalid transaction type")
	}
	return nil
}
