# About this repository

I created this repository just to learn Go and build a personal website API. educational purposes only.
and I created this with the help of LLM or AI, so I can learn the structure and purpose or explanation of each function/file. so don't hesitate to use this repo

# Personal Website API

A Go-based REST API microservice for managing a personal website with authentication, user management, and content management features. It provides endpoints for creating, retrieving, updating, and deleting blog posts, users, roles, personal information, contact information, and social links.

## Features

- 🔐 JWT Authentication
- 👥 User Management (CRUD)
- 🎭 Role-Based Access Control
- 📝 Blog Post Management
- 🛡️ Secure Password Hashing
- 📚 Database Migrations
- 🎯 Clean Architecture

## Tech Stack

- **Language:** Go 1.21+
- **Framework:** Gin
- **Database:** PostgreSQL
- **ORM:** GORM
- **Authentication:** JWT
- **Migration:** golang-migrate

## Project Structure

```
personal-api/
├── api/            # HTTP API specific code
│   └── v1/         # Version 1 of the API
├── cmd/            # Main applications
├── configs/        # Configuration files
├── internal/       # Private application code
├── migrations/     # Database migrations
└── pkg/           # Public libraries
```

## Prerequisites

1. Go 1.21 or higher
2. PostgreSQL 12 or higher
3. golang-migrate CLI

## Setup Instructions

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd personal-api
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Install golang-migrate**
   ```bash
   # Using Go (Recommended)
   go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

   # Alternatively, download from GitHub releases:
   # https://github.com/golang-migrate/migrate/releases
   ```

4. **Set up the environment variables**
   ```bash
   # Copy the example env file
   cp .env.example .env

   # Edit the .env file with your database credentials
   ```

5. **Create the database**
   ```sql
   CREATE DATABASE personal_website;
   ```

6. **Run database migrations**
   ```bash
   # Create a new migration
   migrate create -ext sql -dir migrations -seq migration_name

   # Run migrations up
   migrate -path migrations -database "postgresql://username:password@localhost:5432/personal_website?sslmode=disable" up

   # Run migrations down
   migrate -path migrations -database "postgresql://username:password@localhost:5432/personal_website?sslmode=disable" down

   # Go down specific versions
   migrate -path migrations -database "postgresql://username:password@localhost:5432/personal_website?sslmode=disable" down 1

   # Force a specific version
   migrate -path migrations -database "postgresql://username:password@localhost:5432/personal_website?sslmode=disable" force VERSION
   ```

## Running the Application

1. **Start the server**
   ```bash
   go run cmd/api/main.go
   ```

2. **The server will start at `http://localhost:8080`**

## API Endpoints

### Public Routes

#### Authentication
- `POST /api/v1/auth/register` - Register a new user
- `POST /api/v1/auth/login` - Login user

#### Public Information
- `GET /api/v1/personal` - Get personal information
- `GET /api/v1/contact` - Get contact information
- `GET /api/v1/social` - Get all social links
- `GET /api/v1/social/:id` - Get specific social link

### Protected Routes (Requires Authentication)

#### Posts Management
- `POST /api/v1/posts` - Create a new post
- `GET /api/v1/posts` - Get all posts
- `GET /api/v1/posts/user` - Get current user's posts
- `GET /api/v1/posts/:id` - Get specific post
- `PATCH /api/v1/posts/:id` - Update specific post
- `DELETE /api/v1/posts/:id` - Delete specific post

#### User Management (Admin Only)
- `GET /api/v1/users` - List all users
- `GET /api/v1/users/:id` - Get user details
- `PATCH /api/v1/users/:id` - Update user
- `DELETE /api/v1/users/:id` - Delete user

#### Role Management (Admin Only)
- `POST /api/v1/roles` - Create new role
- `GET /api/v1/roles` - List all roles
- `GET /api/v1/roles/:id` - Get role details
- `PATCH /api/v1/roles/:id` - Update role
- `DELETE /api/v1/roles/:id` - Delete role

#### Personal Information Management (Admin Only)
- `PATCH /api/v1/personal` - Update personal information

#### Contact Information Management (Admin Only)
- `PATCH /api/v1/contact` - Update contact information

#### Social Links Management (Admin Only)
- `POST /api/v1/social` - Create new social link
- `PATCH /api/v1/social/:id` - Update social link
- `DELETE /api/v1/social/:id` - Delete social link

## Authentication

The API uses JWT (JSON Web Tokens) for authentication. To access protected routes:

1. Register a new user or login to get a JWT token
2. Include the token in the Authorization header:
   ```
   Authorization: Bearer <your-jwt-token>
   ```

## Database Migrations

### Running Migrations

1. **Create a new migration**
   ```bash
   migrate create -ext sql -dir migrations -seq migration_name
   ```

2. **Apply migrations**
   ```bash
   # Apply all migrations
   migrate -path migrations -database "postgresql://postgres:postgres@localhost:5432/personal_website?sslmode=disable" up

   # Apply specific number of migrations
   migrate -path migrations -database "postgresql://postgres:postgres@localhost:5432/personal_website?sslmode=disable" up N
   ```

3. **Rollback migrations**
   ```bash
   # Rollback all migrations
   migrate -path migrations -database "postgresql://postgres:postgres@localhost:5432/personal_website?sslmode=disable" down

   # Rollback specific number of migrations
   migrate -path migrations -database "postgresql://postgres:postgres@localhost:5432/personal_website?sslmode=disable" down N
   ```

4. **Other commands**
   ```bash
   # Check current version
   migrate -path migrations -database "postgresql://postgres:postgres@localhost:5432/personal_website?sslmode=disable" version

   # Force a specific version
   migrate -path migrations -database "postgresql://postgres:postgres@localhost:5432/personal_website?sslmode=disable" force VERSION
   ```

### Migration Files Structure
Each migration consists of two files:
- `NNNN_name.up.sql`: Contains the changes to apply
- `NNNN_name.down.sql`: Contains the changes to rollback

## Development

### Code Structure
- `api/v1/handler/` - HTTP request handlers
- `internal/entity/` - Domain models
- `internal/repository/` - Data access layer
- `internal/service/` - Business logic
- `pkg/middleware/` - HTTP middleware
- `pkg/utils/` - Utility functions

### Development Commands

### Common Commands
```bash
# Update Go dependencies
go mod tidy

# Kill the running server process (Windows)
taskkill /F /IM main.exe
```

### Adding New Features
1. Create necessary database migrations
2. Add domain models in `internal/entity/`
3. Implement repository interfaces in `internal/repository/`
4. Add business logic in `internal/service/`
5. Create HTTP handlers in `api/v1/handler/`
6. Register routes in `api/v1/routes/`

## License

MIT License
