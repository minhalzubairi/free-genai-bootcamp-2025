# Language Learning Portal - Backend API

A Go-based REST API for a language learning portal, featuring word groups, study activities, and progress tracking.

## Features

- Word groups management
- Study activities
- Study sessions tracking
- Progress dashboard
- Word reviews
- SQLite database

## Getting Started

### Prerequisites

- Go 1.21 or higher
- SQLite3

### Installation

1. Clone the repository:
```bash
git clone https://github.com/yourusername/lang-portal-backend.git
cd lang-portal-backend
```

2. Install dependencies:
```bash
go mod download
```

3. Run migrations and seed data:
```bash
go run mage.go Migrate
go run mage.go Seed
```

4. Start the server:
```bash
go run cmd/server/main.go
```

The server will start on port 8080.

## API Documentation

See [API.md](API.md) for detailed API documentation.

## Project Structure

```
.
├── cmd/
│   └── server/          # Application entry point
├── internal/
│   ├── handlers/        # HTTP request handlers
│   ├── models/          # Data models
│   └── service/         # Business logic
├── migrations/          # Database migrations
└── seeds/              # Seed data
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.
