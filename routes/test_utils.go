package routes

import (
	"REST_API/auth"
	"REST_API/db"
	"database/sql"
	"testing"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

// CommonTestDB manages shared test database setup and cleanup
type CommonTestDB struct {
	originalDB *sql.DB
	testDB     *sql.DB
	t          *testing.T
}

// SetupTestDB creates a comprehensive in-memory database for all HTTP tests
func SetupTestDB(t *testing.T) *CommonTestDB {
	// Store the original DB
	originalDB := db.DB

	// Create an in-memory database
	testDB, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	// Enable foreign key constraints
	_, err = testDB.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		t.Fatalf("Failed to enable foreign keys: %v", err)
	}

	// Create all tables with a proper schema
	createTables(t, testDB)

	// Create test users with hashed passwords
	createTestUsers(t, testDB)

	// Replace the global DB with test DB
	db.DB = testDB

	return &CommonTestDB{
		originalDB: originalDB,
		testDB:     testDB,
		t:          t,
	}
}

// Cleanup restores the original database connection
func (ctdb *CommonTestDB) Cleanup() {
	if ctdb.testDB != nil {
		err := ctdb.testDB.Close()
		if err != nil {
			return
		}
	}
	db.DB = ctdb.originalDB
}

// createTables creates all required tables for testing
func createTables(t *testing.T, testDB *sql.DB) {
	// Users' table
	createUsersTable := `
		CREATE TABLE users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			email TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL
		)`

	_, err := testDB.Exec(createUsersTable)
	if err != nil {
		t.Fatalf("Failed to create users table: %v", err)
	}

	// Events table
	createEventsTable := `
		CREATE TABLE events (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			description TEXT NOT NULL,
			location TEXT NOT NULL,
			date_time DATETIME NOT NULL,
			user_id INTEGER,
			FOREIGN KEY(user_id) REFERENCES users(id)
		)`

	_, err = testDB.Exec(createEventsTable)
	if err != nil {
		t.Fatalf("Failed to create events table: %v", err)
	}

	// Registration table
	createRegistrationTable := `
		CREATE TABLE registrations (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			event_id INTEGER,
			user_id INTEGER,
			FOREIGN KEY(event_id) REFERENCES events(id),
			FOREIGN KEY(user_id) REFERENCES users(id),
			UNIQUE(event_id, user_id)
		)`

	_, err = testDB.Exec(createRegistrationTable)
	if err != nil {
		t.Fatalf("Failed to create registrations table: %v", err)
	}
}

// createTestUsers creates standard test users for consistent testing
func createTestUsers(t *testing.T, testDB *sql.DB) {
	testUsers := []struct {
		email    string
		password string
	}{
		{"testuser@example.com", "testpassword"},
		{"user1@example.com", "testpassword"},
		{"user2@example.com", "testpassword2"},
		{"logintest@example.com", "testpassword123"},
	}

	for _, user := range testUsers {
		hashedPassword, err := auth.HashPassword(user.password)
		if err != nil {
			t.Fatalf("Failed to hash password for %s: %v", user.email, err)
		}

		_, err = testDB.Exec("INSERT INTO users (email, password) VALUES (?, ?)", user.email, hashedPassword)
		if err != nil {
			t.Fatalf("Failed to create test user %s: %v", user.email, err)
		}
	}
}

// SetupTestRouter creates a test Gin router with all routes configured
func SetupTestRouter() *gin.Engine {
	// Set gin to test mode to reduce output
	gin.SetMode(gin.TestMode)

	// Create a new gin engine
	router := gin.New()

	// Register all routes
	RegisterRoutes(router)

	return router
}

// GenerateTestJWT creates a valid JWT token for testing authenticated routes
func GenerateTestJWT(t *testing.T, userID int64, email string) string {
	token, err := auth.GenerateToken(email, userID)
	if err != nil {
		t.Fatalf("Failed to generate test token: %v", err)
	}
	return token
}

// TestUserCredentials provides standard test user credentials
type TestUserCredentials struct {
	ID       int64
	Email    string
	Password string
}

// GetTestUsers returns predefined test users with their credentials
func GetTestUsers() map[string]TestUserCredentials {
	return map[string]TestUserCredentials{
		"testuser": {
			ID:       1,
			Email:    "testuser@example.com",
			Password: "testpassword",
		},
		"user1": {
			ID:       2,
			Email:    "user1@example.com",
			Password: "testpassword",
		},
		"user2": {
			ID:       3,
			Email:    "user2@example.com",
			Password: "testpassword2",
		},
		"logintest": {
			ID:       4,
			Email:    "logintest@example.com",
			Password: "testpassword123",
		},
	}
}
