# Buku

A modern, self-hosted bookmark management application built with Go and Svelte.

## Overview

Buku (meaning "book" in Indonesian/Malay) is a full-featured bookmark manager that allows you to save, organize, search, and manage your favorite URLs. It features a responsive web interface, robust authentication, and a clean REST API.

## Features

### Core Features
- **User Authentication**: Secure JWT-based authentication with refresh tokens
- **Password Reset**: Token-based password reset functionality
- **Admin Role**: First registered user automatically becomes admin
- **Categories**: Organize bookmarks into custom categories
- **Bookmark Management**: Save URLs with descriptions and organize by category
- **Search**: Full-text search across URLs and descriptions
- **Input Sanitization**: XSS protection with HTML sanitization

### Security Features
- **Argon2id**: Modern password hashing algorithm
- **CORS Support**: Configurable cross-origin resource sharing
- **Content-Type Validation**: JSON validation middleware
- **XSS Protection**: Input sanitization using bluemonday
- **Token Expiration**: Short-lived access tokens (15 min) with refresh tokens (7 days)

### Technical Features
- **Type-Safe SQL**: Generated using sqlc
- **Database Migrations**: Automatic schema migrations with Goose
- **Graceful Shutdown**: Proper server shutdown handling
- **Middleware Chain**: CORS, logging, and validation middleware

## Technology Stack

### Backend
- **Go** 1.26.1
- **SQLite** (modernc.org/sqlite) - Embedded database
- **Goose** - Database migrations
- **sqlc** - Type-safe SQL code generation
- **JWT** (golang-jwt/jwt) - Authentication tokens
- **Argon2id** (alexedwards/argon2id) - Password hashing
- **Bluemonday** - HTML sanitization

### Frontend
- **Svelte** 5.x - Reactive UI framework
- **Vite** - Build tool and dev server
- **Tailwind CSS** 4.x - Utility-first CSS
- **DaisyUI** 5.x - Component library
- **TypeScript** - Type safety

## Project Structure

```
buku/
├── main.go                     # Application entry point
├── server.go                   # HTTP server and routing
├── db.go                       # Database connection and migrations
├── go.mod                      # Go dependencies
├── sqlc.yaml                   # sqlc configuration
├── sql/
│   ├── schema/                 # Database migrations
│   │   ├── 01_users.sql
│   │   ├── 02_refresh_tokens.sql
│   │   ├── 03_categories.sql
│   │   ├── 04_urls.sql
│   │   └── 05_password_reset_tokens.sql
│   └── queries/                # SQL queries for sqlc
│       ├── users.sql
│       ├── categories.sql
│       ├── urls.sql
│       ├── refresh_tokens.sql
│       └── password_reset_tokens.sql
├── internal/
│   ├── handlers/               # HTTP request handlers
│   │   ├── user.go            # Auth & user handlers
│   │   ├── category.go        # Category CRUD
│   │   └── url.go             # Bookmark CRUD & search
│   ├── middlewares/           # HTTP middlewares
│   │   ├── auth.go            # JWT authentication
│   │   ├── cors.go            # CORS handling
│   │   ├── logger.go          # Request logging
│   │   ├── should_json.go     # Content-Type validation
│   │   └── chain.go           # Middleware chaining
│   ├── models/                # Data models
│   │   ├── config.go
│   │   ├── users.go
│   │   ├── categories.go
│   │   └── urls.go
│   ├── database/              # Generated database code (sqlc)
│   └── utils/                 # Utility functions
│       ├── response.go
│       └── sanitize.go        # Input sanitization
├── ui/                        # Frontend application
│   ├── src/
│   │   ├── App.svelte         # Main application
│   │   ├── lib/
│   │   │   └── api.ts         # API client
│   │   └── components/
│   │       ├── Login.svelte
│   │       ├── Register.svelte
│   │       ├── Sidebar.svelte
│   │       ├── BookmarkForm.svelte
│   │       └── BookmarkList.svelte
│   ├── dist/                  # Built files (embedded in Go binary)
│   └── package.json
└── data/
    └── buku.db                # SQLite database
```

## Prerequisites

- **Go** 1.26.1 or later
- **Node.js** 18+ and **bun** or **npm** (for UI development)
- **sqlc** (optional, for regenerating database code)

## Quick Start

### 1. Clone and Setup

```bash
git clone <repository-url>
cd buku
```

### 2. Install Dependencies

```bash
# Go dependencies
go mod download

# UI dependencies
cd ui
bun install
cd ..
```

### 3. Build and Run

```bash
# Build the UI
make build-ui

# Or manually:
cd ui && bun run build && cd ..

# Set environment variable
export JWT_SECRET="your-secret-key-minimum-32-characters-long"

# Run the application
go run .
```

The server will start on `http://localhost:8080`

## Development

### Backend Development

```bash
# Run with auto-reload (requires air)
air

# Or manually:
go run .
```

### Frontend Development

```bash
cd ui
bun run dev
```

The UI dev server runs on `http://localhost:5173` with proxy to backend.

### Database Migrations

Migrations run automatically on startup. To add new migrations:

```bash
# Create a new migration file
echo "-- +goose Up" > sql/schema/06_new_feature.sql

# Regenerate Go code after modifying queries
sqlc generate
```

## API Documentation

### Authentication

#### Register
```http
POST /api/auth/register
Content-Type: application/json

{
  "username": "john_doe",
  "password": "secure_password",
  "name": "John Doe"
}
```

Response:
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "abc123...",
  "expires_in": 900,
  "user": {
    "id": 1,
    "username": "john_doe",
    "name": "John Doe"
  }
}
```

#### Login
```http
POST /api/auth/login
Content-Type: application/json

{
  "username": "john_doe",
  "password": "secure_password"
}
```

#### Password Reset
```http
POST /api/auth/forgot-password
Content-Type: application/json

{
  "username": "john_doe"
}
```

```http
POST /api/auth/reset-password
Content-Type: application/json

{
  "token": "reset_token_here",
  "password": "new_password"
}
```

### Categories

#### List Categories
```http
GET /api/category
Authorization: Bearer <token>
```

#### Create Category
```http
POST /api/category
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "Technology",
  "description": "Tech-related bookmarks"
}
```

#### Update Category
```http
PUT /api/category/{id}
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "Technology & Programming",
  "description": "Updated description"
}
```

#### Delete Category
```http
DELETE /api/category/{id}
Authorization: Bearer <token>
```

### Bookmarks (URLs)

#### List Bookmarks
```http
GET /api/url
Authorization: Bearer <token>

# Optional query parameters:
GET /api/url?category_id=1
GET /api/url?search=github
GET /api/url?category_id=1&search=programming
```

#### Create Bookmark
```http
POST /api/url
Authorization: Bearer <token>
Content-Type: application/json

{
  "url": "https://github.com",
  "description": "GitHub - Code repository",
  "category_id": 1
}
```

#### Update Bookmark
```http
PUT /api/url/{id}
Authorization: Bearer <token>
Content-Type: application/json

{
  "url": "https://github.com",
  "description": "Updated description",
  "category_id": 1
}
```

#### Delete Bookmark
```http
DELETE /api/url/{id}
Authorization: Bearer <token>
```

### Profile

#### Get Profile
```http
GET /api/profile
Authorization: Bearer <token>
```

### Admin

#### Admin Dashboard (Admin only)
```http
GET /api/admin/dashboard
Authorization: Bearer <token>
```

## Database Schema

### Users
| Column | Type | Description |
|--------|------|-------------|
| id | INTEGER | Primary key |
| username | VARCHAR(255) | Unique username |
| password | VARCHAR(255) | Argon2id hashed password |
| name | VARCHAR(255) | Display name |
| is_admin | INTEGER | Admin flag (0/1) |
| created_at | DATETIME | Creation timestamp |
| updated_at | DATETIME | Last update timestamp |

### Categories
| Column | Type | Description |
|--------|------|-------------|
| id | INTEGER | Primary key |
| name | VARCHAR(255) | Category name |
| description | TEXT | Category description |
| user_id | INTEGER | Owner reference |
| created_at | DATETIME | Creation timestamp |
| updated_at | DATETIME | Last update timestamp |

### URLs
| Column | Type | Description |
|--------|------|-------------|
| id | INTEGER | Primary key |
| url | VARCHAR(100) | Bookmark URL |
| description | TEXT | Bookmark description |
| category_id | INTEGER | Category reference |
| user_id | INTEGER | Owner reference |
| created_at | DATETIME | Creation timestamp |
| updated_at | DATETIME | Last update timestamp |

### Refresh Tokens
| Column | Type | Description |
|--------|------|-------------|
| id | INTEGER | Primary key |
| user_id | INTEGER | User reference |
| token_hash | VARCHAR(255) | Token hash |
| expires_at | DATETIME | Expiration timestamp |
| revoked | INTEGER | Revocation flag |
| created_at | DATETIME | Creation timestamp |

### Password Reset Tokens
| Column | Type | Description |
|--------|------|-------------|
| id | INTEGER | Primary key |
| user_id | INTEGER | User reference |
| token_hash | VARCHAR(255) | Token hash |
| expires_at | DATETIME | Expiration timestamp (1 hour) |
| used | BOOLEAN | Usage flag |
| created_at | DATETIME | Creation timestamp |

## Configuration

### Environment Variables

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `JWT_SECRET` | Yes | - | Secret key for JWT signing (min 32 chars) |
| `PORT` | No | `:8080` | Server port |

### Security Notes

1. **JWT Secret**: Must be at least 32 characters long and kept secure
2. **First User**: The first registered user automatically becomes admin
3. **Password Policy**: Minimum 6 characters required
4. **Token Expiry**: Access tokens expire in 15 minutes, refresh tokens in 7 days
5. **XSS Protection**: All user inputs are sanitized using bluemonday

## UI Components

The frontend includes:

- **Login/Register Forms**: Authentication with validation
- **Sidebar**: Category navigation with search
- **Bookmark Form**: Add and edit bookmarks
- **Bookmark List**: Display bookmarks with edit/delete actions
- **Responsive Design**: Works on desktop and mobile

## Testing

### API Testing with cURL

```bash
# Register a new user
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"test","password":"test123","name":"Test User"}'

# Login
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"test","password":"test123"}'

# Create a category (replace <token> with actual token)
curl -X POST http://localhost:8080/api/category \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"name":"Technology"}'

# Search bookmarks
curl "http://localhost:8080/api/url?search=github" \
  -H "Authorization: Bearer <token>"
```

## Deployment

### Docker (Optional)

```dockerfile
# Dockerfile
FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY . .
RUN cd ui && npm install && npm run build
RUN go build -o buku .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/buku .
ENV JWT_SECRET=${JWT_SECRET}
EXPOSE 8080
CMD ["./buku"]
```

### Production Build

```bash
# Build UI for production
cd ui && bun run build && cd ..

# Build Go binary
go build -ldflags="-s -w" -o buku .

# Run
export JWT_SECRET="your-production-secret"
./buku
```

## License

MIT License - see LICENSE file for details

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Acknowledgments

- **DaisyUI** for the beautiful UI components
- **sqlc** for type-safe SQL
- **Argon2id** for secure password hashing
- **Svelte** for the reactive frontend framework
