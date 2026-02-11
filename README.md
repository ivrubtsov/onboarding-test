# Go/Gin Sample REST API

A comprehensive REST API built with Go and the Gin framework, demonstrating best practices for building production-ready web services.
This branch has updates over the main Go branch

## Features

- **User Management**: Full CRUD operations for users with validation
- **Item Management**: Create and manage items linked to users
- **Data Validation**: Built-in request validation using struct tags
- **Query Parameters**: Pagination, filtering, and search capabilities
- **Path Parameters**: Resource identification with validation
- **Error Handling**: Consistent error responses
- **JSON Serialization**: Automatic request/response marshaling
- **Type Safety**: Strong typing throughout with Go's type system
- **RESTful Design**: Following REST principles and conventions

## Project Structure

```
.
├── main.go          # Main application with all handlers
├── go.mod           # Go module dependencies
├── Makefile         # Build and development tasks
└── README.md        # This file
```

## Prerequisites

- Go 1.21 or higher
- Make (optional, for using Makefile commands)

## Installation

1. **Install dependencies**:
```bash
make install
# or
go mod download
```

## Running the Application

### Option 1: Using Make
```bash
make run
```

### Option 2: Using Go directly
```bash
go run main.go
```

### Option 3: Build and run binary
```bash
make build
./bin/server
```

The API will be available at `http://localhost:8080`

## API Endpoints

### Root & Health
- `GET /` - Welcome message and API information
- `GET /health` - Health check endpoint

### Users
- `POST /users/` - Create a new user
- `GET /users/` - Get all users (with optional pagination and role filter)
  - Query params: `skip` (default: 0), `limit` (default: 10, max: 100), `role` (optional: admin/user/guest)
- `GET /users/:id` - Get specific user by ID
- `PUT /users/:id` - Update user information
- `DELETE /users/:id` - Delete user

### Items
- `POST /users/:id/items/` - Create item for a specific user
- `GET /users/:id/items/` - Get all items belonging to a user
- `GET /items/` - Get all items (with optional price filtering)
  - Query params: `skip`, `limit`, `min_price`, `max_price`
- `GET /items/:id` - Get specific item by ID

## Example Usage

### Create a User
```bash
curl -X POST http://localhost:8080/users/ \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "email": "john@example.com",
    "full_name": "John Doe",
    "password": "securepass123",
    "role": "user"
  }'
```

Response:
```json
{
  "id": 1,
  "username": "johndoe",
  "email": "john@example.com",
  "full_name": "John Doe",
  "role": "user",
  "created_at": "2024-01-15T10:30:00Z",
  "is_active": true
}
```

### Get All Users
```bash
curl http://localhost:8080/users/?skip=0&limit=10
```

### Get Users by Role
```bash
curl http://localhost:8080/users/?role=admin
```

### Get Specific User
```bash
curl http://localhost:8080/users/1
```

### Update User
```bash
curl -X PUT http://localhost:8080/users/1 \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe_updated",
    "email": "john.updated@example.com",
    "full_name": "John Doe Updated",
    "role": "admin"
  }'
```

### Delete User
```bash
curl -X DELETE http://localhost:8080/users/1
```

### Create an Item
```bash
curl -X POST http://localhost:8080/users/1/items/ \
  -H "Content-Type: application/json" \
  -d '{
    "name": "MacBook Pro",
    "description": "16-inch M3 Max",
    "price": 3499.99,
    "tax": 349.99
  }'
```

### Get User's Items
```bash
curl http://localhost:8080/users/1/items/
```

### Get All Items with Price Filter
```bash
curl "http://localhost:8080/items/?min_price=100&max_price=2000"
```

### Get Specific Item
```bash
curl http://localhost:8080/items/1
```

### Health Check
```bash
curl http://localhost:8080/health
```

## Validation Rules

### User Validation
- `username`: Required, 3-50 characters
- `email`: Required, valid email format
- `password`: Required (on creation), minimum 8 characters
- `role`: Required, must be one of: admin, user, guest
- `full_name`: Optional

### Item Validation
- `name`: Required, 1-100 characters
- `description`: Optional
- `price`: Required, must be greater than 0
- `tax`: Optional, must be >= 0 if provided

## Development

### Format Code
```bash
make fmt
```

### Run Tests
```bash
make test
```

### Run Linter
```bash
make lint
```

### Clean Build Artifacts
```bash
make clean
```

## Key Go/Gin Concepts Demonstrated

1. **Struct Tags**: Using `json` and `binding` tags for serialization and validation
2. **Gin Router**: Setting up routes and route groups
3. **Middleware**: Using Gin's default middleware (Logger, Recovery)
4. **Request Binding**: `ShouldBindJSON` for automatic parsing and validation
5. **HTTP Status Codes**: Proper status code usage (200, 201, 204, 400, 404)
6. **Error Handling**: Consistent error response format
7. **Query/Path Parameters**: Extracting and validating URL parameters
8. **Type Safety**: Leveraging Go's strong typing system
9. **Slices**: Using Go slices as in-memory data storage
10. **Pointers**: Using pointers for optional fields

## Next Steps for Production

To make this production-ready, consider adding:

- **Database Integration**: PostgreSQL, MySQL, or MongoDB
  - GORM for ORM
  - sqlx for direct SQL access
- **Authentication**: JWT tokens, OAuth2
- **Authorization**: Role-based access control (RBAC)
- **Middleware**: Custom middleware for logging, auth, CORS
- **Configuration**: Environment variables with viper
- **Logging**: Structured logging with zap or logrus
- **Testing**: Unit tests with testify
- **Documentation**: Swagger/OpenAPI with swaggo
- **Docker**: Containerization with multi-stage builds
- **Monitoring**: Prometheus metrics, health checks
- **Rate Limiting**: API rate limiting middleware
- **Caching**: Redis integration
- **Graceful Shutdown**: Handling SIGINT/SIGTERM signals

## Performance Tips

- Use `gin.SetMode(gin.ReleaseMode)` in production
- Add database connection pooling
- Implement caching for frequently accessed data
- Use goroutines for concurrent operations
- Add request timeouts
- Implement pagination for large datasets
- Use prepared statements for database queries

## Dependencies

- **Gin**: High-performance HTTP web framework
- **Validator**: Struct and field validation

## License

MIT

## Contributing

Feel free to submit issues and enhancement requests!
