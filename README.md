# Wallet Service

## Tech Stack

- **Language**: Go 1.24+
- **Framework**: Gin (HTTP router)
- **Database**: PostgreSQL
- **Cache**: Redis
- **Migration**: Goose
- **Documentation**: OpenAPI 3.0 (Swagger)
- **Testing**: Testify + Mockery

## Quick Start

### Prerequisites

- Go 1.24+
- Docker & Docker Compose
- Make (for running commands)

### 1. Clone the Repository

```bash
git clone https://github.com/Shaheen-AlQaraghuli/wallet-go.git
cd wallet-go
```

### 2. Environment Setup

```bash
# Copy environment template
cp .env.example .env

# Edit .env with your configuration
nano .env
```

Example `.env` configuration:
```bash
APP_NAME=wallet-service
APP_DEBUG=true
APP_PORT=8080
DATABASE_DSN=postgres://postgres:example@localhost:5432/wallet?sslmode=disable
REDIS_URL=redis://localhost:6379
```

### 3. Install Dependencies & Tools

```bash
# Install all required tools (linter, swagger, mockery, etc.)
make install-tools

# Download Go dependencies
go mod download
```

### 4. Start Infrastructure Services

```bash
# Start PostgreSQL and Redis
make setup

# Wait a moment for services to be ready, then run migrations
make migrate-up
```

### 5. Run the Service

```bash
# Start the API server
make run

# Or run directly with go
go run ./cmd/server/.
```

The API will be available at `http://localhost:8080/api`

## API Documentation

### Generate Documentation

```bash
# Generate OpenAPI specification
make doc-gen
```

This creates the API documentation in `./docs`

### View Documentation
Go to http://localhost:8080/swagger/index.html


## Development

### Available Commands

```bash
# Development
make run                    # Start the server
make lint                   # Run linter with auto-fix

# Database
make setup                  # Start PostgreSQL and Redis
make db-up                  # Start PostgreSQL
make db-down                # Stop PostgreSQL  
make migrate-up             # Run migrations
make migrate-down           # Rollback last migration
make migrate-reset          # Reset all migrations
make migrate-create name=   # Create new migration

# Documentation
make doc-gen                # Generate API documentatio
```

### Project Structure

```
├── cmd/                    # Application entry point
├── internal/app/           # Application core
│   ├── cache/              # Redis caching layer
│   ├── controller/         # HTTP controllers
│   ├── models/             # Data models
│   ├── repositories/       # Data access layer
│   └── services/           # Business logic
├── pkg/                    # Public packages
│   ├── types/              # Type definitions
│   └── wallet/             # Client SDK
├── database/migrations/    # Database migrations
└── docs/                   # API documentation
```


## Example Usage

### Create a Wallet

```bash
curl -X POST http://localhost:8080/api/v1/wallets \
  -H "Content-Type: application/json" \
  -d '{
    "owner_id": "user-123",
    "currency": "USD"
  }'
```

### Create a Transaction

```bash
curl -X POST http://localhost:8080/api/v1/transactions \
  -H "Content-Type: application/json" \
  -d '{
    "wallet_id": "wallet-123",
    "amount": 1000,
    "type": "credit",
    "idempotency_key": "unique-key-123"
  }'
```
