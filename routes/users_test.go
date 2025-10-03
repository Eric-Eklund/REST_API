package routes

import (
	"REST_API/auth"
	"REST_API/db"
	"REST_API/models"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test POST /signup
func TestSignup(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Cleanup()

	router := SetupTestRouter()

	t.Run("Successful signup", func(t *testing.T) {
		userData := models.User{
			Email:    "newuser@example.com",
			Password: "password123",
		}

		jsonData, err := json.Marshal(userData)
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["message"], "User created successfully")

		var count int
		err = db.DB.QueryRow("SELECT COUNT(*) FROM users WHERE email = ?", userData.Email).Scan(&count)
		assert.NoError(t, err)
		assert.Equal(t, 1, count)
	})

	t.Run("Signup with invalid JSON", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["error"], "Invalid user data")
	})

	t.Run("Signup with missing email", func(t *testing.T) {
		userData := models.User{
			Password: "password123",
		}

		jsonData, _ := json.Marshal(userData)

		req := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["error"], "Invalid user data")
	})

	t.Run("Signup with missing password", func(t *testing.T) {
		userData := models.User{
			Email: "nopassword@example.com",
		}

		jsonData, _ := json.Marshal(userData)

		req := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Signup with duplicate email", func(t *testing.T) {
		userData1 := models.User{
			Email:    "duplicate@example.com",
			Password: "password123",
		}
		jsonData1, _ := json.Marshal(userData1)

		req1 := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(jsonData1))
		req1.Header.Set("Content-Type", "application/json")

		w1 := httptest.NewRecorder()
		router.ServeHTTP(w1, req1)

		assert.Equal(t, http.StatusCreated, w1.Code)

		userData2 := models.User{
			Email:    "duplicate@example.com",
			Password: "differentpassword",
		}
		jsonData2, _ := json.Marshal(userData2)

		req2 := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(jsonData2))
		req2.Header.Set("Content-Type", "application/json")

		w2 := httptest.NewRecorder()
		router.ServeHTTP(w2, req2)

		assert.Equal(t, http.StatusInternalServerError, w2.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w2.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["error"], "User could not be saved")
	})
}

// Test POST /login
func TestLogin(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Cleanup()

	router := SetupTestRouter()

	// Get test user credentials from common utilities
	testUsers := GetTestUsers()
	loginUser := testUsers["logintest"]

	t.Run("Successful login", func(t *testing.T) {
		loginData := models.User{
			Email:    loginUser.Email,
			Password: loginUser.Password,
		}

		jsonData, err := json.Marshal(loginData)
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["message"], "User logged in successfully")
		assert.Contains(t, response, "token")

		// Verify the token is valid
		token := response["token"].(string)
		assert.NotEmpty(t, token)

		// Validate the token
		userID, err := auth.ValidateToken(token)
		assert.NoError(t, err)
		assert.Equal(t, loginUser.ID, userID)
	})

	t.Run("Login with wrong password", func(t *testing.T) {
		loginData := models.User{
			Email:    loginUser.Email,
			Password: "wrongpassword",
		}

		jsonData, _ := json.Marshal(loginData)

		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["error"], "Invalid credentials")
	})

	t.Run("Login with non-existent email", func(t *testing.T) {
		loginData := models.User{
			Email:    "nonexistent@example.com",
			Password: "anypassword",
		}

		jsonData, _ := json.Marshal(loginData)

		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["error"], "Invalid credentials")
	})

	t.Run("Login with invalid JSON", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["error"], "Invalid user data")
	})

	t.Run("Login with missing credentials", func(t *testing.T) {
		loginData := models.User{}

		jsonData, _ := json.Marshal(loginData)

		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

// Integration test: Signup then Login
func TestSignupThenLogin(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Cleanup()

	router := SetupTestRouter()

	t.Run("Complete user flow: signup then login", func(t *testing.T) {
		userData := models.User{
			Email:    "integration@example.com",
			Password: "integrationtest123",
		}

		jsonData, err := json.Marshal(userData)
		assert.NoError(t, err)

		// Signup request
		signupReq := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(jsonData))
		signupReq.Header.Set("Content-Type", "application/json")

		signupW := httptest.NewRecorder()
		router.ServeHTTP(signupW, signupReq)

		assert.Equal(t, http.StatusCreated, signupW.Code)

		// Login with same credentials
		loginReq := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonData))
		loginReq.Header.Set("Content-Type", "application/json")

		loginW := httptest.NewRecorder()
		router.ServeHTTP(loginW, loginReq)

		assert.Equal(t, http.StatusOK, loginW.Code)

		var loginResponse map[string]interface{}
		err = json.Unmarshal(loginW.Body.Bytes(), &loginResponse)
		assert.NoError(t, err)
		assert.Contains(t, loginResponse["message"], "User logged in successfully")
		assert.Contains(t, loginResponse, "token")

		// Token should be valid
		token := loginResponse["token"].(string)
		userID, err := auth.ValidateToken(token)
		assert.NoError(t, err)
		assert.NotZero(t, userID)
	})
}
