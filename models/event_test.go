package models

import (
	"REST_API/db"
	"database/sql"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func setupEventTestDB(t *testing.T) func() {
	originalDB := db.DB

	testDB, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	_, err = testDB.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		t.Fatalf("Failed to enable foreign keys: %v", err)
	}

	createUsersTable := `
		CREATE TABLE users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			email TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL
		)`

	_, err = testDB.Exec(createUsersTable)
	if err != nil {
		t.Fatalf("Failed to create users table: %v", err)
	}

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

	_, err = testDB.Exec("INSERT INTO users (email, password) VALUES (?, ?)", "testuser@example.com", "hashedpassword")
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	db.DB = testDB

	return func() {
		err := testDB.Close()
		if err != nil {
			return
		}
		db.DB = originalDB
	}
}

func TestEvent_Save(t *testing.T) {
	cleanup := setupEventTestDB(t)
	defer cleanup()

	tests := []struct {
		name    string
		event   Event
		wantErr bool
	}{
		{
			name: "Valid event",
			event: Event{
				Name:        "Test Conference",
				Description: "Annual tech conference",
				Location:    "Convention Center",
				DateTime:    time.Now().Add(24 * time.Hour),
				UserID:      1, // References to the test user we created
			},
			wantErr: false,
		},
		{
			name: "Valid event with empty fields (SQLite allows this)",
			event: Event{
				Name:        "",
				Description: "",
				Location:    "",
				DateTime:    time.Now(),
				UserID:      1,
			},
			wantErr: false,
		},
		{
			name: "Invalid event with non-existent user ID",
			event: Event{
				Name:        "Test Conference",
				Description: "Annual tech conference",
				Location:    "Convention Center",
				DateTime:    time.Now().Add(24 * time.Hour),
				UserID:      999,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.event.Save()
			if (err != nil) != tt.wantErr {
				t.Errorf("Save() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && tt.event.ID == 0 {
				t.Errorf("Save() should set event ID, but ID is still 0")
			}
		})
	}
}

func TestEvent_Update(t *testing.T) {
	cleanup := setupEventTestDB(t)
	defer cleanup()

	event := &Event{
		Name:        "Original Event",
		Description: "Original description",
		Location:    "Original location",
		DateTime:    time.Now().Add(24 * time.Hour),
		UserID:      1,
	}

	err := event.Save()
	if err != nil {
		t.Fatalf("Failed to create test event: %v", err)
	}

	tests := []struct {
		name     string
		updateFn func(*Event)
		wantErr  bool
	}{
		{
			name: "Valid update",
			updateFn: func(e *Event) {
				e.Name = "Updated Event"
				e.Description = "Updated description"
			},
			wantErr: false,
		},
		{
			name: "Update with empty name (allowed)",
			updateFn: func(e *Event) {
				e.Name = ""
			},
			wantErr: false,
		},
		{
			name: "Update with invalid user ID",
			updateFn: func(e *Event) {
				e.UserID = 999
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testEvent := *event
			tt.updateFn(&testEvent)

			err := testEvent.Update()
			if (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEvent_Delete(t *testing.T) {
	cleanup := setupEventTestDB(t)
	defer cleanup()

	event := &Event{
		Name:        "Event to Delete",
		Description: "This event will be deleted",
		Location:    "Test location",
		DateTime:    time.Now().Add(24 * time.Hour),
		UserID:      1,
	}

	err := event.Save()
	if err != nil {
		t.Fatalf("Failed to create test event: %v", err)
	}

	t.Run("Successful delete", func(t *testing.T) {
		err := event.Delete()
		if err != nil {
			t.Errorf("Delete() error = %v", err)
		}

		_, err = GetEventByID(event.ID)
		if err == nil {
			t.Error("Event should not exist after deletion")
		}
	})

	t.Run("Delete non-existent event", func(t *testing.T) {
		nonExistentEvent := &Event{ID: 999}
		err := nonExistentEvent.Delete()
		if err != nil {
			t.Errorf("Delete() error = %v", err)
		}
	})
}

func TestEvent_Register(t *testing.T) {
	cleanup := setupEventTestDB(t)
	defer cleanup()

	event := &Event{
		Name:        "Registration Test Event",
		Description: "Event for testing registration",
		Location:    "Test location",
		DateTime:    time.Now().Add(24 * time.Hour),
		UserID:      1,
	}

	err := event.Save()
	if err != nil {
		t.Fatalf("Failed to create test event: %v", err)
	}

	tests := []struct {
		name    string
		userID  int64
		wantErr bool
	}{
		{
			name:    "Valid registration",
			userID:  1,
			wantErr: false,
		},
		{
			name:    "Invalid user ID",
			userID:  999,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := event.Register(tt.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Register() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

	t.Run("Duplicate registration should fail", func(t *testing.T) {
		testEvent := &Event{
			Name:        "Duplicate Registration Test",
			Description: "Event for testing duplicate registration",
			Location:    "Test location",
			DateTime:    time.Now().Add(48 * time.Hour),
			UserID:      1,
		}
		err := testEvent.Save()
		if err != nil {
			t.Fatalf("Failed to create test event: %v", err)
		}

		err = testEvent.Register(1)
		if err != nil {
			t.Fatalf("First registration should succeed: %v", err)
		}

		err = testEvent.Register(1)
		if err == nil {
			t.Error("Duplicate registration should fail")
		}
	})
}

func TestEvent_Unregister(t *testing.T) {
	cleanup := setupEventTestDB(t)
	defer cleanup()

	event := &Event{
		Name:        "Unregistration Test Event",
		Description: "Event for testing unregistration",
		Location:    "Test location",
		DateTime:    time.Now().Add(24 * time.Hour),
		UserID:      1,
	}

	err := event.Save()
	if err != nil {
		t.Fatalf("Failed to create test event: %v", err)
	}

	err = event.Register(1)
	if err != nil {
		t.Fatalf("Failed to register user: %v", err)
	}

	t.Run("Successful unregistration", func(t *testing.T) {
		err := event.Unregister(1)
		if err != nil {
			t.Errorf("Unregister() error = %v", err)
		}
	})

	t.Run("Unregister non-registered user", func(t *testing.T) {
		err := event.Unregister(999)
		if err != nil {
			t.Errorf("Unregister() error = %v", err)
		}
	})
}

func TestGetAllEvents(t *testing.T) {
	cleanup := setupEventTestDB(t)
	defer cleanup()

	events := []*Event{
		{
			Name:        "Event 1",
			Description: "First event",
			Location:    "Location 1",
			DateTime:    time.Now().Add(24 * time.Hour),
			UserID:      1,
		},
		{
			Name:        "Event 2",
			Description: "Second event",
			Location:    "Location 2",
			DateTime:    time.Now().Add(48 * time.Hour),
			UserID:      1,
		},
	}

	for _, event := range events {
		err := event.Save()
		if err != nil {
			t.Fatalf("Failed to create test event: %v", err)
		}
	}

	t.Run("Get all events", func(t *testing.T) {
		allEvents, err := GetAllEvents()
		if err != nil {
			t.Errorf("GetAllEvents() error = %v", err)
			return
		}

		if len(allEvents) != 2 {
			t.Errorf("GetAllEvents() returned %d events, want 2", len(allEvents))
		}
	})
}

func TestGetEventByID(t *testing.T) {
	cleanup := setupEventTestDB(t)
	defer cleanup()

	event := &Event{
		Name:        "Test Event for Retrieval",
		Description: "Event for testing retrieval by ID",
		Location:    "Test location",
		DateTime:    time.Now().Add(24 * time.Hour),
		UserID:      1,
	}

	err := event.Save()
	if err != nil {
		t.Fatalf("Failed to create test event: %v", err)
	}

	tests := []struct {
		name    string
		id      int64
		wantErr bool
	}{
		{
			name:    "Valid ID",
			id:      event.ID,
			wantErr: false,
		},
		{
			name:    "Invalid ID",
			id:      999,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			retrievedEvent, err := GetEventByID(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetEventByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if retrievedEvent.ID != event.ID {
					t.Errorf("GetEventByID() returned wrong event ID = %d, want %d", retrievedEvent.ID, event.ID)
				}
				if retrievedEvent.Name != event.Name {
					t.Errorf("GetEventByID() returned wrong event name = %s, want %s", retrievedEvent.Name, event.Name)
				}
			}
		})
	}
}
