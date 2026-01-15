package domain

import (
	"testing"
	"time"
)

func TestProductValidation(t *testing.T) {
	tests := []struct {
		name    string
		product *Product
		wantErr bool
	}{
		{
			name: "Valid product",
			product: &Product{
				Name:        "Laptop",
				SKU:         "LAP001",
				Description: "Gaming Laptop",
				Price:       1500.00,
			},
			wantErr: false,
		},
		{
			name: "Missing name",
			product: &Product{
				Name:        "",
				SKU:         "LAP001",
				Description: "Gaming Laptop",
				Price:       1500.00,
			},
			wantErr: true,
		},
		{
			name: "Missing SKU",
			product: &Product{
				Name:        "Laptop",
				SKU:         "",
				Description: "Gaming Laptop",
				Price:       1500.00,
			},
			wantErr: true,
		},
		{
			name: "Negative price",
			product: &Product{
				Name:        "Laptop",
				SKU:         "LAP001",
				Description: "Gaming Laptop",
				Price:       -100.00,
			},
			wantErr: true,
		},
		{
			name: "Zero price",
			product: &Product{
				Name:        "Laptop",
				SKU:         "LAP001",
				Description: "Gaming Laptop",
				Price:       0,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.product.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInventoryItemAvailableQuantity(t *testing.T) {
	tests := []struct {
		name     string
		item     *InventoryItem
		expected int64
	}{
		{
			name: "Available quantity with no reserved items",
			item: &InventoryItem{
				Quantity: 100,
				Reserved: 0,
			},
			expected: 100,
		},
		{
			name: "Available quantity with some reserved items",
			item: &InventoryItem{
				Quantity: 100,
				Reserved: 30,
			},
			expected: 70,
		},
		{
			name: "No available quantity (all reserved)",
			item: &InventoryItem{
				Quantity: 100,
				Reserved: 100,
			},
			expected: 0,
		},
		{
			name: "Negative available quantity (invalid state)",
			item: &InventoryItem{
				Quantity: 50,
				Reserved: 100,
			},
			expected: 0,
		},
		{
			name: "Zero quantities",
			item: &InventoryItem{
				Quantity: 0,
				Reserved: 0,
			},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.item.AvailableQuantity()
			if got != tt.expected {
				t.Errorf("AvailableQuantity() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestInventoryItemValidation(t *testing.T) {
	tests := []struct {
		name    string
		item    *InventoryItem
		wantErr bool
	}{
		{
			name: "Valid inventory item",
			item: &InventoryItem{
				ProductID: "prod-1",
				Quantity:  100,
				Reserved:  10,
				Location:  "Warehouse A",
			},
			wantErr: false,
		},
		{
			name: "Missing product ID",
			item: &InventoryItem{
				ProductID: "",
				Quantity:  100,
				Reserved:  10,
				Location:  "Warehouse A",
			},
			wantErr: true,
		},
		{
			name: "Negative quantity",
			item: &InventoryItem{
				ProductID: "prod-1",
				Quantity:  -10,
				Reserved:  0,
				Location:  "Warehouse A",
			},
			wantErr: true,
		},
		{
			name: "Negative reserved",
			item: &InventoryItem{
				ProductID: "prod-1",
				Quantity:  100,
				Reserved:  -5,
				Location:  "Warehouse A",
			},
			wantErr: true,
		},
		{
			name: "Reserved exceeds quantity",
			item: &InventoryItem{
				ProductID: "prod-1",
				Quantity:  50,
				Reserved:  100,
				Location:  "Warehouse A",
			},
			wantErr: true,
		},
		{
			name: "Missing location",
			item: &InventoryItem{
				ProductID: "prod-1",
				Quantity:  100,
				Reserved:  10,
				Location:  "",
			},
			wantErr: true,
		},
		{
			name: "Zero quantity",
			item: &InventoryItem{
				ProductID: "prod-1",
				Quantity:  0,
				Reserved:  0,
				Location:  "Warehouse A",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.item.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTransactionValidation(t *testing.T) {
	tests := []struct {
		name    string
		tx      *Transaction
		wantErr bool
	}{
		{
			name: "Valid IN transaction",
			tx: &Transaction{
				InventoryID: "inv-1",
				ProductID:   "prod-1",
				Type:        "IN",
				Quantity:    50,
				Reference:   "PO-001",
				CreatedAt:   time.Now(),
			},
			wantErr: false,
		},
		{
			name: "Valid OUT transaction",
			tx: &Transaction{
				InventoryID: "inv-1",
				ProductID:   "prod-1",
				Type:        "OUT",
				Quantity:    20,
				Reference:   "ORDER-001",
				CreatedAt:   time.Now(),
			},
			wantErr: false,
		},
		{
			name: "Valid RESERVE transaction",
			tx: &Transaction{
				InventoryID: "inv-1",
				ProductID:   "prod-1",
				Type:        "RESERVE",
				Quantity:    10,
				Reference:   "ORDER-002",
				CreatedAt:   time.Now(),
			},
			wantErr: false,
		},
		{
			name: "Valid UNRESERVE transaction",
			tx: &Transaction{
				InventoryID: "inv-1",
				ProductID:   "prod-1",
				Type:        "UNRESERVE",
				Quantity:    5,
				Reference:   "ORDER-002",
				CreatedAt:   time.Now(),
			},
			wantErr: false,
		},
		{
			name: "Valid RETURN transaction",
			tx: &Transaction{
				InventoryID: "inv-1",
				ProductID:   "prod-1",
				Type:        "RETURN",
				Quantity:    5,
				Reference:   "RETURN-001",
				CreatedAt:   time.Now(),
			},
			wantErr: false,
		},
		{
			name: "Missing inventory ID",
			tx: &Transaction{
				InventoryID: "",
				ProductID:   "prod-1",
				Type:        "IN",
				Quantity:    50,
				CreatedAt:   time.Now(),
			},
			wantErr: true,
		},
		{
			name: "Missing product ID",
			tx: &Transaction{
				InventoryID: "inv-1",
				ProductID:   "",
				Type:        "IN",
				Quantity:    50,
				CreatedAt:   time.Now(),
			},
			wantErr: true,
		},
		{
			name: "Invalid transaction type",
			tx: &Transaction{
				InventoryID: "inv-1",
				ProductID:   "prod-1",
				Type:        "INVALID",
				Quantity:    50,
				CreatedAt:   time.Now(),
			},
			wantErr: true,
		},
		{
			name: "Zero quantity",
			tx: &Transaction{
				InventoryID: "inv-1",
				ProductID:   "prod-1",
				Type:        "IN",
				Quantity:    0,
				CreatedAt:   time.Now(),
			},
			wantErr: true,
		},
		{
			name: "Negative quantity",
			tx: &Transaction{
				InventoryID: "inv-1",
				ProductID:   "prod-1",
				Type:        "IN",
				Quantity:    -10,
				CreatedAt:   time.Now(),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.tx.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
