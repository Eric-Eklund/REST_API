# REST API - Event Management System

A simple and efficient REST API built with Go and the Gin framework for managing events. This project provides endpoints to create and retrieve events with in-memory storage.

## 🚀 Features

- **Create Events**: Add new events with comprehensive details
- **Retrieve Events**: Get a list of all events
- **JSON API**: RESTful API with JSON request/response format
- **Input Validation**: Built-in validation for required fields
- **Lightweight**: Fast and efficient using the Gin web framework

## 🛠️ Tech Stack

- **Language**: Go 1.25
- **Web Framework**: [Gin](https://github.com/gin-gonic/gin) v1.11.0
- **Storage**: In-memory (for demonstration purposes)
- **API Format**: JSON REST API

## 📋 API Endpoints

### Create Event
- **Endpoint**: `POST /events`
- **Content-Type**: `application/json`
- **Description**: Creates a new event

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

### Get All Events
- **Endpoint**: `GET /events`
- **Description**: Retrieves all events

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

The project includes HTTP test files in the `api-test/` directory that you can use with tools like:
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

## 📁 Project Structure

```
REST_API/
├── main.go              # Main application entry point
├── models/              # Data models and business logic
│   └── event.go        # Event model and data operations
├── api-test/           # HTTP test files
│   ├── create-event.http
│   └── get-events.http
├── go.mod              # Go module dependencies
├── go.sum              # Dependency checksums
├── .gitignore          # Git ignore rules
└── README.md           # Project documentation
```

## 🎯 Event Model

The Event model includes the following fields:

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `id` | int | No | Auto-generated event ID |
| `name` | string | Yes | Event name |
| `description` | string | Yes | Event description |
| `location` | string | Yes | Event location |
| `date_time` | time.Time | Yes | Event date and time |
| `user_id` | int | No | Auto-assigned user ID |

## 🔮 Future Enhancements

- [ ] Database integration (PostgreSQL/MySQL)
- [ ] User authentication and authorization
- [ ] Event update and deletion endpoints
- [ ] Event filtering and search capabilities
- [ ] Pagination for large event lists
- [ ] Input sanitization and advanced validation
- [ ] Unit and integration tests
- [ ] Docker containerization
- [ ] API documentation with Swagger

## 🤝 Contributing

Feel free to fork this project and submit pull requests for any improvements.

## 📝 License

This project is open source and available under the [MIT License](LICENSE).

## 👨‍💻 Author

**Eric Eklund** - [GitHub](https://github.com/Eric-Eklund)

---

*Built with ❤️ using Go and Gin*