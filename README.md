# Distributed Inventory System

[![CI](https://github.com/bhnrathore/distributed-inventory-system/actions/workflows/ci.yml/badge.svg)](https://github.com/bhnrathore/distributed-inventory-system/actions)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.25.5-blue)](https://golang.org)

A high-performance, production-grade distributed inventory management system built in Go. Demonstrates enterprise-level software architecture, API design, database management, and Go best practices.

## Features

- **RESTful API**: Clean HTTP API for inventory operations
- **Product Management**: Create, update, list, and delete products
- **Stock Management**: Add, remove, reserve, and unreserve stock
- **Transaction History**: Track all inventory movements
- **Atomic Operations**: Thread-safe stock operations
- **PostgreSQL**: Robust relational database with proper indexing
- **Error Handling**: Comprehensive error handling and validation
- **Logging & Monitoring**: Structured logging and request tracking
- **Clean Architecture**: Domain-driven design with clear separation of concerns
- **Unit Tests**: Comprehensive test coverage

## Project Structure

```
.
├── cmd/
│   └── server/           # Application entry point
├── internal/
│   ├── api/             # HTTP handlers and middleware
│   ├── domain/          # Domain models and business logic entities
│   ├── repository/      # Data access layer
│   └── service/         # Business logic layer
├── docker-compose.yml   # Docker services for dependencies
├── .env.example         # Example environment variables
└── go.mod              # Go module definition
```

## Architecture

The system follows a clean, layered architecture:

1. **Domain Layer** (`internal/domain/`): Defines core entities (Product, InventoryItem, Transaction)
2. **Repository Layer** (`internal/repository/`): Handles data persistence with PostgreSQL
3. **Service Layer** (`internal/service/`): Contains business logic and validations
4. **API Layer** (`internal/api/`): HTTP handlers, middleware, and request/response formatting
5. **Main** (`cmd/server/main.go`): Orchestrates initialization and starts the server

## Prerequisites

- Go 1.25.5 or later
- PostgreSQL 12 or later (or use Docker)
- Docker & Docker Compose (optional)

## Setup

### 1. Install Dependencies

```bash
go mod download
```

### 2. Start PostgreSQL

Using Docker Compose:
```bash
docker-compose up -d
```

Or use an existing PostgreSQL instance and set the `DATABASE_URL` environment variable.

### 3. Build

```bash
go build -o bin/server ./cmd/server/
```

### 4. Run

```bash
./bin/server
```

The server will start on `http://localhost:8080`

## API Endpoints

### Health Check
- **GET** `/health` - Check server health

### Products
- **POST** `/api/products` - Create a new product
  ```json
  {
    "name": "Laptop",
    "description": "Gaming Laptop",
    "sku": "LAP001",
    "price": 1500.00,
    "location": "Warehouse A",
    "initial_quantity": 50
  }
  ```

- **GET** `/api/products` - List all products (supports pagination)
  - Query params: `limit=10&offset=0`

- **GET** `/api/products/{id}` - Get product details with inventory

- **PUT** `/api/products/{id}` - Update product
  ```json
  {
    "name": "Updated Name",
    "description": "Updated Description",
    "price": 1600.00
  }
  ```

- **DELETE** `/api/products/{id}` - Delete product

### Stock Operations
- **POST** `/api/products/{id}/stock/add` - Add stock
  ```json
  {
    "quantity": 20,
    "reference": "PO-001",
    "notes": "Purchase order from supplier"
  }
  ```

- **POST** `/api/products/{id}/stock/remove` - Remove stock
  ```json
  {
    "quantity": 5,
    "reference": "ORDER-123"
  }
  ```

- **POST** `/api/products/{id}/stock/reserve` - Reserve stock
  ```json
  {
    "quantity": 10,
    "reference": "ORDER-456"
  }
  ```

- **POST** `/api/products/{id}/stock/unreserve` - Unreserve stock
  ```json
  {
    "quantity": 5,
    "reference": "ORDER-456"
  }
  ```

### Inventory & History
- **GET** `/api/products/{id}/inventory` - Get inventory details

- **GET** `/api/products/{id}/transactions` - Get transaction history
  - Query params: `limit=10&offset=0`

## Testing

Run unit tests:
```bash
go test ./...
```

With coverage:
```bash
go test -cover ./...
```

## Design Patterns & Best Practices

1. **Domain-Driven Design**: Core entities in domain package
2. **Repository Pattern**: Abstract data access with interfaces
3. **Dependency Injection**: Services receive dependencies
4. **Error Handling**: Custom error messages with context
5. **Middleware**: Composable HTTP middleware for cross-cutting concerns
6. **Atomic Operations**: Database constraints ensure data consistency
7. **Logging**: Structured logging for debugging and monitoring
8. **Pagination**: Efficient list operations with limit/offset
9. **Validation**: Input validation at domain and API layers
10. **Graceful Shutdown**: Proper server shutdown handling

## Performance Considerations

- **Connection Pooling**: Configured database connection pool
- **Indexes**: Database indexes on frequently queried columns
- **Prepared Statements**: Parameterized queries prevent SQL injection
- **Context Usage**: Proper timeout handling with context
- **Minimal Dependencies**: Lean dependency list for fast compilation

## Future Enhancements

- [ ] Distributed caching with Redis
- [ ] Event-driven architecture with message queues
- [ ] Multi-warehouse support with cross-location transfers
- [ ] Real-time inventory updates with WebSockets
- [ ] Advanced analytics and reporting
- [ ] API authentication and authorization
- [ ] Rate limiting and API versioning
- [ ] Containerized deployment (Kubernetes)

## License

MIT License

