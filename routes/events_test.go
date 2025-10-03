package routes

import (
	"REST_API/models"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// createTestEvent helper function to create events in the database
func createTestEvent(t *testing.T, userID int64) *models.Event {
	event := &models.Event{
		Name:        "Test Event",
		Description: "Test Description",
		Location:    "Test Location",
		DateTime:    time.Now().Add(24 * time.Hour),
		UserID:      userID,
	}

	err := event.Save()
	if err != nil {
		t.Fatalf("Failed to create test event: %v", err)
	}

	return event
}

// Test GET /events - Public endpoint
func TestGetEvents(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Cleanup()

	router := SetupTestRouter()

	// Create some test events
	createTestEvent(t, 1)
	createTestEvent(t, 1)

	t.Run("Get all events successfully", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/events", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var events []models.Event
		err := json.Unmarshal(w.Body.Bytes(), &events)
		assert.NoError(t, err)
		assert.Len(t, events, 2)
	})
}

// Test GET /events/:id - Public endpoint
func TestGetEventByID(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Cleanup()

	router := SetupTestRouter()

	// Create a test event
	event := createTestEvent(t, 1)

	t.Run("Get event by valid ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/events/"+strconv.FormatInt(event.ID, 10), nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var retrievedEvent models.Event
		err := json.Unmarshal(w.Body.Bytes(), &retrievedEvent)
		assert.NoError(t, err)
		assert.Equal(t, event.ID, retrievedEvent.ID)
		assert.Equal(t, event.Name, retrievedEvent.Name)
	})

	t.Run("Get event by invalid ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/events/999", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["error"], "Event not found")
	})

	t.Run("Get event by malformed ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/events/invalid", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["error"], "Invalid event ID")
	})
}

// Test POST /events - Authenticated endpoint
func TestCreateEvent(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Cleanup()

	router := SetupTestRouter()

	// Generate JWT token for authentication
	testUsers := GetTestUsers()
	user := testUsers["testuser"]
	token := GenerateTestJWT(t, user.ID, user.Email)

	t.Run("Create event successfully", func(t *testing.T) {
		eventData := models.Event{
			Name:        "New Test Event",
			Description: "New Test Description",
			Location:    "New Test Location",
			DateTime:    time.Now().Add(48 * time.Hour),
		}

		jsonData, err := json.Marshal(eventData)
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/events", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", token)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var createdEvent models.Event
		err = json.Unmarshal(w.Body.Bytes(), &createdEvent)
		assert.NoError(t, err)
		assert.NotZero(t, createdEvent.ID)
		assert.Equal(t, eventData.Name, createdEvent.Name)
		assert.Equal(t, user.ID, createdEvent.UserID)
	})

	t.Run("Create event without authentication", func(t *testing.T) {
		eventData := models.Event{
			Name:        "Unauthorized Event",
			Description: "This should fail",
			Location:    "Nowhere",
			DateTime:    time.Now().Add(48 * time.Hour),
		}

		jsonData, err := json.Marshal(eventData)
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/events", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Create event with invalid JSON", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/events", bytes.NewBuffer([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", token)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["error"], "Invalid event data")
	})
}

// Test PUT /events/:id - Authenticated endpoint
func TestUpdateEvent(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Cleanup()

	router := SetupTestRouter()

	testUsers := GetTestUsers()
	user := testUsers["testuser"]
	token := GenerateTestJWT(t, user.ID, user.Email)

	// Create a test event owned by user 1
	event := createTestEvent(t, user.ID)

	t.Run("Update event successfully", func(t *testing.T) {
		updatedData := models.Event{
			Name:        "Updated Event Name",
			Description: "Updated Description",
			Location:    "Updated Location",
			DateTime:    time.Now().Add(72 * time.Hour),
			UserID:      user.ID,
		}

		jsonData, err := json.Marshal(updatedData)
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodPut, "/events/"+strconv.FormatInt(event.ID, 10), bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", token)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["message"], "Event updated successfully")
	})

	t.Run("Update non-existent event", func(t *testing.T) {
		updatedData := models.Event{Name: "Updated"}
		jsonData, _ := json.Marshal(updatedData)

		req := httptest.NewRequest(http.MethodPut, "/events/999", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", token)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

// Test DELETE /events/:id - Authenticated endpoint
func TestDeleteEvent(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Cleanup()

	router := SetupTestRouter()

	testUsers := GetTestUsers()
	user := testUsers["testuser"]
	token := GenerateTestJWT(t, user.ID, user.Email)

	// Create a test event
	event := createTestEvent(t, user.ID)

	t.Run("Delete event successfully", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/events/"+strconv.FormatInt(event.ID, 10), nil)
		req.Header.Set("Authorization", token)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["message"], "Event deleted successfully")

		// Verify the event is actually deleted
		_, err = models.GetEventByID(event.ID)
		assert.Error(t, err)
	})
}
