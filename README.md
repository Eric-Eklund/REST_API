# REST API - Event Management System

A comprehensive and efficient REST API built with Go and the Gin framework for managing events with JWT-based user authentication. This project provides full CRUD operations for events, secure user registration/login, event registration system, and persistent SQLite database storage.

## 🚀 Features

- **JWT Authentication**: Complete JWT token-based authentication system with secure login/logout
- **User Management**: Secure user registration and login with bcrypt password hashing
- **Event Registration System**: Users can register/unregister for events with protected endpoints
- **Complete CRUD Operations**: Create, Read, Update, and Delete events
- **Protected Routes**: Authentication middleware protecting sensitive operations
- **Database Persistence**: SQLite database with relational schema and foreign keys
- **RESTful Design**: Clean REST API endpoints following best practices
- **Structured Architecture**: Organized codebase with separate packages for routes, models, database, and authentication
- **JSON API**: RESTful API with JSON request/response format
- **Input Validation**: Built-in validation for required fields
- **Password Security**: bcrypt hashing for secure password storage
- **Token Security**: JWT tokens with expiration and validation
- **Error Handling**: Comprehensive error handling and HTTP status codes
- **Database Connection Pooling**: Optimized database connections
- **Comprehensive Testing**: Full test suite with unit and integration tests
- **Test Utilities**: Reusable test helpers for consistent testing across components
- **Lightweight**: Fast and efficient using the Gin web framework

## 🛠️ Tech Stack

- **Language**: Go 1.25
- **Web Framework**: [Gin](https://github.com/gin-gonic/gin) v1.11.0
- **Database**: SQLite 3 with [go-sqlite3](https://github.com/mattn/go-sqlite3) driver
- **Authentication**: [JWT](https://github.com/golang-jwt/jwt/v5) for token-based authentication
- **Password Hashing**: [bcrypt](https://golang.org/x/crypto/bcrypt) for secure password storage
- **API Format**: JSON REST API
- **Architecture**: Clean separation of concerns with packages

## 📋 API Endpoints

### User Authentication

#### User Registration
- **Endpoint**: `POST /signup`
- **Content-Type**: `application/json`
- **Description**: Register a new user with email and password

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "your_secure_password"
}
```

**Response (Success):**
```json
{
  "message": "User created successfully"
}
```

**Response (Error):**
```json
{
  "error": "Invalid user data"
}
```

#### User Login
- **Endpoint**: `POST /login`
- **Content-Type**: `application/json`
- **Description**: Authenticate user with email and password

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "your_secure_password"
}
```

**Response (Success):**
```json
{
  "message": "User logged in successfully",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Response (Error):**
```json
{
  "error": "Invalid credentials"
}
```

### Event Management

#### Get All Events
- **Endpoint**: `GET /events`
- **Authentication**: Not required
- **Description**: Retrieves all events from the database

**Response:**
```json
[
  {
    "id": 1,
    "name": "Event Name",
    "description": "Event description",
    "location": "Event location",
    "date_time": "2025-01-01T13:37:00.000Z",
    "user_id": 1337
  }
]
```

#### Get Event by ID
- **Endpoint**: `GET /events/{id}`
- **Authentication**: Not required
- **Description**: Retrieves a specific event by its ID

**Response:**
```json
{
  "id": 1,
  "name": "Event Name",
  "description": "Event description",
  "location": "Event location",
  "date_time": "2025-01-01T13:37:00.000Z",
  "user_id": 1337
}
```

#### Create Event
- **Endpoint**: `POST /events`
- **Content-Type**: `application/json`
- **Authentication**: Required (JWT token)
- **Description**: Creates a new event and stores it in the database

**Request Body:**
```json
{
  "name": "Event Name",
  "description": "Event description",
  "location": "Event location",
  "date_time": "2025-01-01T13:37:00.000Z"
}
```

**Response:**
```json
{
  "id": 1,
  "name": "Event Name",
  "description": "Event description",
  "location": "Event location",
  "date_time": "2025-01-01T13:37:00.000Z",
  "user_id": 1337
}
```

#### Update Event
- **Endpoint**: `PUT /events/{id}`
- **Content-Type**: `application/json`
- **Authentication**: Required (JWT token)
- **Description**: Updates an existing event by ID

**Request Body:**
```json
{
  "name": "Updated Event Name",
  "description": "Updated description",
  "location": "Updated location",
  "date_time": "2025-01-01T13:37:00.000Z"
}
```

**Response:**
```json
{
  "message": "Event updated successfully"
}
```

#### Delete Event
- **Endpoint**: `DELETE /events/{id}`
- **Authentication**: Required (JWT token)
- **Description**: Deletes an event by ID

**Response:**
```json
{
  "message": "Event deleted successfully"
}
```

### Event Registration

#### Register for Event
- **Endpoint**: `POST /events/{id}/register`
- **Authentication**: Required (JWT token)
- **Description**: Register the authenticated user for a specific event

**Headers:**
```
Authorization: YOUR_JWT_TOKEN_HERE
```

**Response (Success):**
```json
{
  "message": "Event registered successfully"
}
```

**Response (Error):**
```json
{
  "error": "Event could not be registered"
}
```

#### Unregister from Event
- **Endpoint**: `DELETE /events/{id}/register`
- **Authentication**: Required (JWT token)
- **Description**: Unregister the authenticated user from a specific event

**Headers:**
```
Authorization: YOUR_JWT_TOKEN_HERE
```

**Response (Success):**
```json
{
  "message": "Event unregistered successfully"
}
```

**Response (Error):**
```json
{
  "error": "Event could not be unregistered"
}
```

## 🏃‍♂️ Getting Started

### Prerequisites

- Go 1.25 or higher
- Git

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/Eric-Eklund/REST_API.git
   cd REST_API
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Run the server:
   ```bash
   go run main.go
   ```

The server will start on `http://localhost:8080`

## 🧪 Testing

This project includes a comprehensive test suite covering all major functionality:

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run specific test package
go test ./routes
go test ./models

# Run tests with coverage
go test -cover ./...
```

### Test Structure

The test suite is organized with reusable utilities and comprehensive coverage:

- **`routes/test_utils.go`**: Shared test utilities including:
  - `SetupTestDB()`: In-memory SQLite database setup with sample data
  - `SetupTestRouter()`: Gin router configuration for testing
  - `GenerateTestJWT()`: JWT token generation for authenticated tests
  - `GetTestUsers()`: Standard test user credentials

- **Unit Tests**: Individual component testing for models and utilities
- **Integration Tests**: End-to-end API testing with database interactions
- **Authentication Tests**: JWT token validation and user authentication flows
- **Registration Tests**: Event registration/unregistration functionality with helper functions

### Test Coverage

- **User Authentication**: Registration, login, JWT validation
- **Event Management**: CRUD operations, validation, error handling
- **Event Registration**: User registration/unregistration for events
- **Database Operations**: Model methods, foreign key constraints, data integrity
- **API Endpoints**: HTTP status codes, JSON responses, authentication middleware
- **Error Scenarios**: Invalid inputs, unauthorized access, non-existent resources

### Testing the API

The project includes comprehensive HTTP test files in the `api-test/` directory:
- `create-event.http` - Test event creation
- `get-events.http` - Test getting all events and specific events by ID
- `update-events.http` - Test event updates
- `delete-events.http` - Test event deletion
- `create-user.http` - Test user registration
- `login.http` - Test user login
- `registration.http` - Test event registration
- `unregistration.http` - Test event unregistration

You can use these with tools like:
- JetBrains HTTP Client (built into GoLand/IntelliJ IDEA)
- VS Code REST Client extension
- Postman
- cURL

#### Example cURL commands:

**User Registration:**
```bash
curl -X POST http://localhost:8080/signup \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "secure_password"
  }'
```

**User Login:**
```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "secure_password"
  }'
```

**Create an event (requires authentication):**
```bash
curl -X POST http://localhost:8080/events \
  -H "Content-Type: application/json" \
  -H "Authorization: YOUR_JWT_TOKEN_HERE" \
  -d '{
    "name": "Sample Event",
    "description": "This is a test event",
    "location": "Stockholm, Sweden",
    "date_time": "2025-01-01T13:37:00.000Z"
  }'
```

**Get all events:**
```bash
curl http://localhost:8080/events
```

**Get a specific event:**
```bash
curl http://localhost:8080/events/1
```

**Update an event (requires authentication):**
```bash
curl -X PUT http://localhost:8080/events/1 \
  -H "Content-Type: application/json" \
  -H "Authorization: YOUR_JWT_TOKEN_HERE" \
  -d '{
    "name": "Updated Event",
    "description": "Updated description",
    "location": "Updated location",
    "date_time": "2025-01-01T13:37:00.000Z"
  }'
```

**Delete an event (requires authentication):**
```bash
curl -X DELETE http://localhost:8080/events/1 \
  -H "Authorization: YOUR_JWT_TOKEN_HERE"
```

**Register for an event (requires authentication):**
```bash
curl -X POST http://localhost:8080/events/1/register \
  -H "Authorization: YOUR_JWT_TOKEN_HERE"
```

**Unregister from an event (requires authentication):**
```bash
curl -X DELETE http://localhost:8080/events/1/register \
  -H "Authorization: YOUR_JWT_TOKEN_HERE"
```

## 📁 Project Structure

```
REST_API/
├── main.go              # Main application entry point
├── db/                  # Database package
│   └── db.go            # Database initialization and setup
├── models/              # Data models and business logic
│   ├── event.go         # Event model with CRUD operations
│   ├── event_test.go    # Event model unit tests
│   ├── user.go          # User model with authentication
│   └── user_test.go     # User model unit tests
├── routes/              # Route handlers
│   ├── events.go        # Event-related route handlers
│   ├── events_test.go   # Event route integration tests
│   ├── users.go         # User authentication route handlers
│   ├── users_test.go    # User authentication route tests
│   ├── register.go      # Event registration route handlers
│   ├── register_test.go # Event registration route tests
│   ├── routes.go        # Route registration and middleware setup
│   └── test_utils.go    # Shared test utilities and helpers
├── auth/                # Authentication package
│   ├── auth.go          # Authentication middleware
│   ├── hash.go          # Password hashing and validation
│   └── jwt.go           # JWT token generation and validation
├── api-test/            # HTTP test files
│   ├── create-event.http # Event POST request tests
│   ├── get-events.http   # Event GET request tests
│   ├── update-events.http # Event PUT request tests
│   ├── delete-events.http # Event DELETE request tests
│   ├── create-user.http  # User registration tests
│   ├── login.http        # User login tests
│   ├── registration.http # Event registration tests
│   └── unregistration.http # Event unregistration tests
├── api.db               # SQLite database file (auto-generated)
├── go.mod               # Go module dependencies
├── go.sum               # Dependency checksums
├── .gitignore           # Git ignore rules
└── README.md            # Project documentation
```

## 🎯 Data Models

### User Model

The User model handles authentication and user management:

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `id` | int64 | No | Auto-generated primary key (SQLite AUTOINCREMENT) |
| `email` | string | Yes | User email address (unique) |
| `password` | string | Yes | bcrypt hashed password |

**Database Operations:**
- **Registration**: `Save()` method creates new users with hashed passwords
- **Authentication**: `ValidateCredentials()` method verifies login credentials
- **JWT Integration**: Login returns JWT tokens for authenticated sessions
- **Security**: All passwords are hashed using bcrypt before storage

### Event Model

The Event model includes the following fields stored in SQLite database:

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `id` | int64 | No | Auto-generated primary key (SQLite AUTOINCREMENT) |
| `name` | string | Yes | Event name |
| `description` | string | Yes | Event description |
| `location` | string | Yes | Event location |
| `date_time` | time.Time | Yes | Event date and time (SQLite DATETIME) |
| `user_id` | int | No | Foreign key reference to users table |

**Database Operations:**
- **Create**: `Save()` method inserts new events into database (requires authentication)
- **Read**: `GetAllEvents()` and `GetEventByID()` functions for querying (public access)
- **Update**: `Update()` method modifies existing events (requires authentication)
- **Delete**: `Delete()` method removes events from database (requires authentication)
- **Registration**: `Register()` and `Unregister()` methods for event registration (requires authentication)

### Database Schema

The application uses SQLite with the following tables:

**Users Table:**
- Primary key: `id` (INTEGER AUTOINCREMENT)
- Unique constraint on `email`
- Password stored as bcrypt hash

**Events Table:**
- Primary key: `id` (INTEGER AUTOINCREMENT) 
- Foreign key: `user_id` references `users(id)`
- Proper relational integrity with foreign key constraints

**Event Registrations Table:**
- Composite primary key: `user_id`, `event_id`
- Foreign keys: `user_id` references `users(id)`, `event_id` references `events(id)`
- Manages many-to-many relationship between users and events

## 🔮 Future Enhancements

- [x] ~~User authentication and authorization~~ ✅ **Completed**
- [x] ~~JWT token-based authentication~~ ✅ **Completed**
- [x] ~~Event registration system~~ ✅ **Completed**
- [ ] User-specific event access control (only event creators can modify)
- [ ] Event filtering and search capabilities
- [ ] Pagination for large event lists
- [ ] Input sanitization and advanced validation
- [x] ~~Unit and integration tests~~ ✅ **Completed**
- [ ] Docker containerization
- [ ] API documentation with Swagger
- [ ] Database migration system
- [ ] PostgreSQL/MySQL support
- [ ] Logging middleware
- [ ] Rate limiting
- [ ] CORS support for web frontends
- [ ] Password reset functionality
- [ ] Email verification for user registration

## 🤝 Contributing

Feel free to fork this project and submit pull requests for any improvements.

## 📝 License

This project is open source and available under the [MIT License](https://opensource.org/license/mit).

## 👨‍💻 Author

**Eric Eklund** - [GitHub](https://github.com/Eric-Eklund)

---

*Built with ❤️ using Go and Gin*