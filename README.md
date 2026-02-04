# English Learning Backend Service

This is the backend service for the English Learning application, built with Go. It strictly follows **Clean Architecture** principles to ensure modularity, testability, and maintainability.

## 🚀 Tech Stack

- **Language**: Go (Golang) 1.23+
- **Web Framework**: [Gin](https://github.com/gin-gonic/gin)
- **Database**: PostgreSQL
- **ORM**: [GORM](https://gorm.io)
- **Configuration**: [Godotenv](https://github.com/joho/godotenv) & [Viper](https://github.com/spf13/viper)
- **Logging**: [Zap](https://github.com/uber-go/zap)
- **Authentication**: JWT (Access & Refresh Tokens) with Revocable Sessions

## 🏛️ Architecture & Project Structure

The project implements **Clean Architecture** with a clear separation of concerns using **Ports and Adapters**.

### Directory Layout

```text
.
├── cmd/
│   └── server/          # Application entry point (main.go)
├── configs/             # Configuration files & loaders
├── internal/
│   ├── adapters/        # Interface Adapters (Implementation details)
│   │   ├── handlers/    # HTTP Handlers (Gin) - maps HTTP requests to Service calls
│   │   └── repositories/# Database Implementations (GORM) - implements core ports
│   └── core/            # Core Business Logic (The "Circle of Life")
│       ├── domain/      # Sub-Domain Entities (Data Structures Only)
│       ├── ports/       # Interfaces (Ports) - Defines contracts for Repositories & Services
│       └── services/    # Use Cases - Implements business logic using Ports
└── pkg/
    ├── logger/          # Global logger setup
    └── response/        # Standardized API Response helpers
```

### Architectural Rules

1.  **Dependency Rule**: Dependencies point **inwards**.
    - `adapters` depends on `core`.
    - `core` **never** depends on `adapters` or external frameworks.
2.  **Domain (`internal/core/domain`)**:
    - Contains **ONLY** pure data structures (structs).
    - **NO** interfaces, behavior, or methods that depend on external libraries.
    - Examples: `User`, `Session`, `TokenPair`.
3.  **Ports (`internal/core/ports`)**:
    - Contains **ALL** interfaces (contracts) for repositories and services.
    - Separating ports from domain prevents circular dependencies and keeps entities pure.
    - Examples: `UserRepository`, `SessionRepository`, `AuthService`.
4.  **Services (`internal/core/services`)**:
    - Implement business logic.
    - Depend ONLY on `ports`.
5.  **Adapters (`internal/adapters`)**:
    - **Handlers**: Parse JSON, validate input, call Service, return JSON.
    - **Repositories**: Implement `ports` using specific drivers (e.g., Postgres/GORM).

## ✨ Features

- **Standardized Responses**: Unified API response format `{ data, code, message }`.
- **Secure Authentication**:
  - JWT Access & Refresh Tokens.
  - **Revocable Sessions**: Sessions are tracked in DB.
  - **Session Rotation**: Refresh tokens are rotated on use.
  - **Session Tracking**: Captures User Agent and Client IP.
  - **Logout**: Revokes session immediately.
- **Robust Validation**: Request validation using `validator/v10`.
- **Configuration**: Environment-based config via `.env`.

## 🛠️ Getting Started

### Prerequisites

- Go 1.23+
- PostgreSQL
- Docker (optional)

### Installation

1.  **Clone the repository**:

    ```bash
    git clone <repo-url>
    cd english-learning
    ```

2.  **Environment Setup**:
    Create a `.env` file based on your config needs:

    ```env
    ENV=dev
    PORT=8080
    DB_DSN="host=localhost user=postgres password=password dbname=english_learning port=5432 sslmode=disable"
    JWT_SECRET=supersecretkey
    ACCESS_EXPIRY_HOUR=24
    REFRESH_EXPIRY_HOUR=168
    ```

3.  **Run the server**:

    ```bash
    go mod tidy
    go run cmd/server/main.go
    ```

    _Using Docker Compose:_

    ```bash
    docker-compose up --build
    ```

## 🔌 API Endpoints

### Auth

- `POST /auth/register`: Register new user.
- `POST /auth/login`: Login (Returns Access + Refresh Token).
- `POST /auth/refresh`: Rotate Refresh Token & Get new Access Token.
- `POST /auth/logout`: Revoke current session.

### Users

- `GET /users`: List users (Admin).
- `POST /users`: Create user manually.
- `GET /users/:id`: Get profile.
