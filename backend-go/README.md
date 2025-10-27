# EV Warranty Management System - Backend API

A robust RESTful API backend service for managing Electric Vehicle warranty claims, built with Go. This service provides comprehensive warranty claim management, authentication, OAuth integration, and cloud-based file storage capabilities.

## Features

- **Warranty Claim Management** - Create, track, and manage EV warranty claims with detailed item tracking and attachment support
- **Authentication & Authorization** - Secure JWT-based authentication with refresh token support and role-based access control
- **OAuth Integration** - Google OAuth 2.0 integration for seamless user authentication
- **Cloud Storage** - Cloudinary integration for secure claim attachment storage and management
- **Database Migrations** - Version-controlled database schema management with PostgreSQL
- **API Documentation** - Auto-generated Swagger/OpenAPI documentation for all endpoints

## Prerequisites

- Go 1.24.0 or higher
- PostgreSQL 12+
- [golang-migrate](https://github.com/golang-migrate/migrate) CLI tool
- Google OAuth credentials (for OAuth features)
- Cloudinary account (for file uploads)

##  Tech Stack

- **Language**: Go 1.24
- **Web Framework**: Gin
- **Database**: PostgreSQL with GORM ORM
- **Authentication**: JWT (golang-jwt/jwt)
- **OAuth**: Google OAuth 2.0
- **Cloud Storage**: Cloudinary
- **API Documentation**: Swagger
- **Testing**: Ginkgo & Gomega
- **Mocking**: Mockery

## Installation

### 1. Clone the repository

```bash
cd backend-go
```

### 2. Install dependencies

```bash
go mod download
```

### 3. Set up environment variables

Create a `.env` file in the root directory:

```bash
cp .env.example .env
```

Edit `.env` with your configuration:

```env
PORT=8080
DATABASE_URL=postgres://username:password@localhost:5432/ev_warranty?sslmode=disable
LOG_LEVEL=info

# JWT Keys
PRIVATE_KEY_PATH=./keys/private.pem
PUBLIC_KEY_PATH=./keys/public.pem
ACCESS_TOKEN_TTL=15m
REFRESH_TOKEN_TTL=168h

# Google OAuth
GOOGLE_CLIENT_ID=your-google-client-id
GOOGLE_CLIENT_SECRET=your-google-client-secret
GOOGLE_REDIRECT_URL=http://localhost:8080/api/v1/auth/oauth/google/callback
FRONTEND_BASE_URL=http://localhost:3000

# Cloudinary
CLOUDINARY_URL=cloudinary://api_key:api_secret@cloud_name
CLOUDINARY_UPLOAD_FOLDER=ev-warranty-claim-attachment

# Admin Setup
ADMIN_EMAIL=admin@example.com
ADMIN_PASSWORD=Admin@123
```

### 4. Generate RSA keys for JWT

```bash
mkdir -p keys
ssh-keygen -t rsa -b 4096 -m PEM -f keys/private.pem -N ""
openssl rsa -in keys/private.pem -pubout -outform PEM -out keys/public.pem
```

### 5. Set up the database

Make sure PostgreSQL is running, then run migrations:

```bash
make db/migrations/up
```

## Usage

### Running the Application

Start the development server:

```bash
make run
```

The API will be available at `http://localhost:8080`

### API Documentation

Access the Swagger UI documentation:

```
http://localhost:8080/swagger/index.html
```

### Example API Requests

#### Register a new user

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "SecurePass123!",
    "full_name": "John Doe",
    "phone": "+1234567890"
  }'
```

#### Login

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "SecurePass123!"
  }'
```

#### Create a warranty claim (authenticated)

```bash
curl -X POST http://localhost:8080/api/v1/claims \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "customer_id": "customer-uuid",
    "description": "Battery performance degradation",
    "vehicle_id": "vehicle-uuid"
  }'
```

## 🔧 Makefile Commands

The project includes a comprehensive Makefile for common development tasks:

### Development

```bash
# Run the application
make run

# Run code quality checks
make check

# Format, vet, and test the code
make audit
```

### Database Management

```bash
# Connect to PostgreSQL database
make db/psql

# Create a new migration
make db/migrations/new name=migration_name

# Run all pending migrations (requires confirmation)
make db/migrations/up

# Rollback the last migration (requires confirmation)
make db/migrations/down

# Force migration to a specific version (requires confirmation)
make db/migrations/force version=1

# Check current migration version
make db/migrations/version
```

### Testing

```bash
# Generate test mocks
make test/mocks

# Run tests with coverage
make test/cover

# Generate HTML coverage report
make test/cover/html

# Show coverage by function
make test/cover/func
```

### Documentation

```bash
# Generate/update Swagger documentation
make swagger/gen
```

### Utility

```bash
# Confirm prompts (used internally by other commands)
make confirm
```

## ⚙️ Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Server port | `8080` |
| `DATABASE_URL` | PostgreSQL connection string | - |
| `LOG_LEVEL` | Logging level (debug, info, warn, error) | `info` |
| `PRIVATE_KEY_PATH` | Path to RSA private key | `./keys/private.pem` |
| `PUBLIC_KEY_PATH` | Path to RSA public key | `./keys/public.pem` |
| `ACCESS_TOKEN_TTL` | Access token expiration duration | `15m` |
| `REFRESH_TOKEN_TTL` | Refresh token expiration duration | `168h` |
| `GOOGLE_CLIENT_ID` | Google OAuth client ID | - |
| `GOOGLE_CLIENT_SECRET` | Google OAuth client secret | - |
| `GOOGLE_REDIRECT_URL` | OAuth callback URL | - |
| `FRONTEND_BASE_URL` | Frontend application URL | `http://localhost:3000` |
| `CLOUDINARY_URL` | Cloudinary connection URL | - |
| `CLOUDINARY_UPLOAD_FOLDER` | Upload folder name in Cloudinary | `ev-warranty-claim-attachment` |

## 📁 Project Structure

```
backend-go/
├── cmd/
│   └── server/           # Application entry point
├── internal/
│   ├── application/      # Business logic and services
│   ├── domain/          # Domain models and entities
│   ├── infrastructure/  # External services (DB, OAuth, Cloudinary)
│   ├── interfaces/      # HTTP handlers and API routes
│   └── security/        # Security utilities (JWT, password)
├── migrations/          # Database migrations
├── docs/               # Swagger documentation
├── pkg/                # Reusable packages
│   ├── logger/        # Logging utilities
│   └── mocks/         # Generated mocks for testing
├── keys/              # RSA keys for JWT (gitignored)
└── Makefile           # Build and development commands
```

## 🧪 Testing

Run the test suite:

```bash
# Run all tests
go test ./...

# Run tests with coverage
make test/cover

# Generate HTML coverage report
make test/cover/html
```

## 🐳 Docker Support

Build and run with Docker:

```bash
# Build the image
docker build -t ev-warranty-backend .

# Run the container
docker run -p 8080:8080 --env-file .env ev-warranty-backend
```

**Note**: Remember to keep your `.env` file secure and never commit it to version control. Always use `.env.example` as a template for required environment variables.
