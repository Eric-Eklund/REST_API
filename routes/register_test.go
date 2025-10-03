package routes

import (
	"REST_API/db"
	"REST_API/models"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// createTestEventForRegistration creates a test event for registration tests
func createTestEventForRegistration(t *testing.T, userID int64) *models.Event {
	event := &models.Event{
		Name:        "Registration Test Event",
		Description: "Event for testing registration functionality",
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

func makeRegistrationRequest(router *gin.Engine, method, eventID, token string) *httptest.ResponseRecorder {
	url := "/events/" + eventID + "/register"
	req := httptest.NewRequest(method, url, nil)
	if token != "" {
		req.Header.Set("Authorization", token)
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	return w
}

// assertResponseAndMessage checks status code and extracts message/error from response
func assertResponseAndMessage(t *testing.T, w *httptest.ResponseRecorder, expectedStatus int, expectedText, messageKey string) {
	assert.Equal(t, expectedStatus, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response[messageKey], expectedText)
}

// verifyRegistrationCount checks the number of registrations for an event/user combination
func verifyRegistrationCount(t *testing.T, eventID, userID int64, expectedCount int) {
	var count int
	err := db.DB.QueryRow("SELECT COUNT(*) FROM registrations WHERE event_id = ? AND user_id = ?", eventID, userID).Scan(&count)
	assert.NoError(t, err)
	assert.Equal(t, expectedCount, count)
}

// Test POST /events/:id/register - Register for an event
func TestRegisterEvent(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Cleanup()

	router := SetupTestRouter()

	// Create a test event
	event := createTestEventForRegistration(t, 1) // Event owned by user 1

	// Generate JWT token for user 2 (different from an event owner)
	testUsers := GetTestUsers()
	user := testUsers["user2"]
	token := GenerateTestJWT(t, user.ID, user.Email)

	t.Run("Successful event registration", func(t *testing.T) {
		w := makeRegistrationRequest(router, http.MethodPost, strconv.FormatInt(event.ID, 10), token)
		assertResponseAndMessage(t, w, http.StatusCreated, "Event registered successfully", "message")
		verifyRegistrationCount(t, event.ID, user.ID, 1)
	})

	t.Run("Register for non-existent event", func(t *testing.T) {
		w := makeRegistrationRequest(router, http.MethodPost, "999", token)
		assertResponseAndMessage(t, w, http.StatusInternalServerError, "Event not found", "error")
	})

	t.Run("Register with invalid event ID", func(t *testing.T) {
		w := makeRegistrationRequest(router, http.MethodPost, "invalid", token)
		assertResponseAndMessage(t, w, http.StatusBadRequest, "Invalid event ID", "error")
	})

	t.Run("Register without authentication", func(t *testing.T) {
		w := makeRegistrationRequest(router, http.MethodPost, strconv.FormatInt(event.ID, 10), "")
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Duplicate registration should fail", func(t *testing.T) {
		testEvent := createTestEventForRegistration(t, 1)
		eventID := strconv.FormatInt(testEvent.ID, 10)

		w1 := makeRegistrationRequest(router, http.MethodPost, eventID, token)
		assert.Equal(t, http.StatusCreated, w1.Code)

		w2 := makeRegistrationRequest(router, http.MethodPost, eventID, token)
		assertResponseAndMessage(t, w2, http.StatusInternalServerError, "Event could not be registered", "error")
	})
}

// Test DELETE /events/:id/register - Unregister from an event
func TestUnregisterEvent(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Cleanup()

	router := SetupTestRouter()

	event := createTestEventForRegistration(t, 1)

	// Generate JWT token for user 2
	testUsers := GetTestUsers()
	user := testUsers["user2"]
	token := GenerateTestJWT(t, user.ID, user.Email)

	err := event.Register(user.ID)
	assert.NoError(t, err)

	t.Run("Successful event unregistration", func(t *testing.T) {
		w := makeRegistrationRequest(router, http.MethodDelete, strconv.FormatInt(event.ID, 10), token)
		assertResponseAndMessage(t, w, http.StatusOK, "Event unregistered successfully", "message")
		verifyRegistrationCount(t, event.ID, user.ID, 0)
	})

	t.Run("Unregister from non-existent event", func(t *testing.T) {
		w := makeRegistrationRequest(router, http.MethodDelete, "999", token)
		assertResponseAndMessage(t, w, http.StatusInternalServerError, "Event not found", "error")
	})

	t.Run("Unregister with invalid event ID", func(t *testing.T) {
		w := makeRegistrationRequest(router, http.MethodDelete, "invalid", token)
		assertResponseAndMessage(t, w, http.StatusBadRequest, "Invalid event ID", "error")
	})

	t.Run("Unregister without authentication", func(t *testing.T) {
		w := makeRegistrationRequest(router, http.MethodDelete, strconv.FormatInt(event.ID, 10), "")
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Unregister when not registered", func(t *testing.T) {
		newEvent := createTestEventForRegistration(t, 1)
		w := makeRegistrationRequest(router, http.MethodDelete, strconv.FormatInt(newEvent.ID, 10), token)
		assertResponseAndMessage(t, w, http.StatusOK, "Event unregistered successfully", "message")
	})
}

// Integration test: Register then Unregister
func TestRegisterThenUnregister(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Cleanup()

	router := SetupTestRouter()

	event := createTestEventForRegistration(t, 1)

	// Generate JWT token for user 2
	testUsers := GetTestUsers()
	user := testUsers["user2"]
	token := GenerateTestJWT(t, user.ID, user.Email)

	t.Run("Complete registration flow: register then unregister", func(t *testing.T) {
		eventID := strconv.FormatInt(event.ID, 10)

		// Register
		registerW := makeRegistrationRequest(router, http.MethodPost, eventID, token)
		assert.Equal(t, http.StatusCreated, registerW.Code)
		verifyRegistrationCount(t, event.ID, user.ID, 1)

		// Unregister
		unregisterW := makeRegistrationRequest(router, http.MethodDelete, eventID, token)
		assert.Equal(t, http.StatusOK, unregisterW.Code)
		verifyRegistrationCount(t, event.ID, user.ID, 0)
	})
}

// Test with different users
func TestRegistrationWithMultipleUsers(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Cleanup()

	router := SetupTestRouter()

	// Create a test event
	event := createTestEventForRegistration(t, 1)

	// Generate JWT tokens for different users
	testUsers := GetTestUsers()
	user1 := testUsers["user1"]
	user2 := testUsers["user2"]
	token1 := GenerateTestJWT(t, user1.ID, user1.Email)
	token2 := GenerateTestJWT(t, user2.ID, user2.Email)

	t.Run("Multiple users can register for same event", func(t *testing.T) {
		eventID := strconv.FormatInt(event.ID, 10)

		// User 1 registers
		w1 := makeRegistrationRequest(router, http.MethodPost, eventID, token1)
		assert.Equal(t, http.StatusCreated, w1.Code)

		// User 2 registers
		w2 := makeRegistrationRequest(router, http.MethodPost, eventID, token2)
		assert.Equal(t, http.StatusCreated, w2.Code)

		// Verify both registrations exist
		var count int
		err := db.DB.QueryRow("SELECT COUNT(*) FROM registrations WHERE event_id = ?", event.ID).Scan(&count)
		assert.NoError(t, err)
		assert.Equal(t, 2, count)
	})

	t.Run("User can only unregister their own registration", func(t *testing.T) {
		newEvent := createTestEventForRegistration(t, 1)
		err := newEvent.Register(user2.ID)
		assert.NoError(t, err)

		// User 2 unregisters (should work)
		w := makeRegistrationRequest(router, http.MethodDelete, strconv.FormatInt(newEvent.ID, 10), token2)
		assert.Equal(t, http.StatusOK, w.Code)
		verifyRegistrationCount(t, newEvent.ID, user2.ID, 0)
	})
}
