# REST API - Event Management System

A comprehensive and efficient REST API built with Go and the Gin framework for managing events. This project provides full CRUD operations for events with persistent SQLite database storage.

## 🚀 Features

- **Complete CRUD Operations**: Create, Read, Update, and Delete events
- **Database Persistence**: SQLite database for data storage
- **RESTful Design**: Clean REST API endpoints following best practices
- **Structured Architecture**: Organized codebase with separate packages for routes, models, and database
- **JSON API**: RESTful API with JSON request/response format
- **Input Validation**: Built-in validation for required fields
- **Error Handling**: Comprehensive error handling and HTTP status codes
- **Database Connection Pooling**: Optimized database connections
- **Lightweight**: Fast and efficient using the Gin web framework

## 🛠️ Tech Stack

- **Language**: Go 1.25
- **Web Framework**: [Gin](https://github.com/gin-gonic/gin) v1.11.0
- **Database**: SQLite 3 with [go-sqlite3](https://github.com/mattn/go-sqlite3) driver
- **API Format**: JSON REST API
- **Architecture**: Clean separation of concerns with packages

## 📋 API Endpoints

### Get All Events
- **Endpoint**: `GET /events`
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

### Get Event by ID
- **Endpoint**: `GET /events/{id}`
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

### Create Event
- **Endpoint**: `POST /events`
- **Content-Type**: `application/json`
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

### Update Event
- **Endpoint**: `PUT /events/{id}`
- **Content-Type**: `application/json`
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

### Delete Event
- **Endpoint**: `DELETE /events/{id}`
- **Description**: Deletes an event by ID

**Response:**
```json
{
  "message": "Event deleted successfully"
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

### Testing the API

The project includes comprehensive HTTP test files in the `api-test/` directory:
- `create-event.http` - Test event creation
- `get-events.http` - Test getting all events and specific events by ID
- `update-events.http` - Test event updates
- `delete-events.http` - Test event deletion

You can use these with tools like:
- JetBrains HTTP Client (built into GoLand/IntelliJ IDEA)
- VS Code REST Client extension
- Postman
- cURL

#### Example cURL commands:

**Create an event:**
```bash
curl -X POST http://localhost:8080/events \
  -H "Content-Type: application/json" \
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

**Update an event:**
```bash
curl -X PUT http://localhost:8080/events/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Updated Event",
    "description": "Updated description",
    "location": "Updated location",
    "date_time": "2025-01-01T13:37:00.000Z"
  }'
```

**Delete an event:**
```bash
curl -X DELETE http://localhost:8080/events/1
```

## 📁 Project Structure

```
REST_API/
├── main.go              # Main application entry point
├── db/                  # Database package
│   └── db.go            # Database initialization and setup
├── models/              # Data models and business logic
│   └── event.go         # Event model with CRUD operations
├── routes/              # Route handlers
│   ├── events.go        # Event-related route handlers
│   └── routes.go        # Route registration
├── api-test/            # HTTP test files
│   ├── create-event.http # POST request tests
│   ├── get-events.http   # GET request tests
│   ├── update-events.http # PUT request tests
│   └── delete-events.http # DELETE request tests
├── api.db               # SQLite database file (auto-generated)
├── go.mod               # Go module dependencies
├── go.sum               # Dependency checksums
├── .gitignore           # Git ignore rules
└── README.md            # Project documentation
```

## 🎯 Event Model

The Event model includes the following fields stored in SQLite database:

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `id` | int64 | No | Auto-generated primary key (SQLite AUTOINCREMENT) |
| `name` | string | Yes | Event name |
| `description` | string | Yes | Event description |
| `location` | string | Yes | Event location |
| `date_time` | time.Time | Yes | Event date and time (SQLite DATETIME) |
| `user_id` | int | No | Auto-assigned user ID |

### Database Operations

The Event model supports full CRUD operations:
- **Create**: `Save()` method inserts new events into database
- **Read**: `GetAllEvents()` and `GetEventByID()` functions for querying
- **Update**: `Update()` method modifies existing events
- **Delete**: `Delete()` method removes events from database

## 🔮 Future Enhancements

- [ ] User authentication and authorization
- [ ] Event filtering and search capabilities
- [ ] Pagination for large event lists
- [ ] Input sanitization and advanced validation
- [ ] Unit and integration tests
- [ ] Docker containerization
- [ ] API documentation with Swagger
- [ ] Database migration system
- [ ] PostgreSQL/MySQL support
- [ ] Logging middleware
- [ ] Rate limiting
- [ ] CORS support for web frontends

## 🤝 Contributing

Feel free to fork this project and submit pull requests for any improvements.

## 📝 License

This project is open source and available under the [MIT License](https://opensource.org/license/mit).

## 👨‍💻 Author

**Eric Eklund** - [GitHub](https://github.com/Eric-Eklund)

---

*Built with ❤️ using Go and Gin*