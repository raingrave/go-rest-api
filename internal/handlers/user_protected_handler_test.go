package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/raingrave/apirest/internal/models"
	"github.com/raingrave/apirest/internal/repositories"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

// Helper function to create a user, log them in, and return a valid JWT.
func getAuthToken(t *testing.T) (string, models.User) {
	password := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := models.User{
		Name:     "Protected User",
		Email:    "protected@example.com",
		Password: string(hashedPassword),
	}
	id, err := userRepo.CreateUser(user)
	assert.NoError(t, err)
	user.ID = id

	loginCredentials := map[string]string{
		"email":    user.Email,
		"password": password,
	}
	body, _ := json.Marshal(loginCredentials)
	req, _ := http.NewRequest("POST", "/api/v1/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	return response["token"], user
}

func TestProtectedUserEndpoints(t *testing.T) {
	clearUserTable()
	token, user := getAuthToken(t)

	t.Run("List Users - Success", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/users", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var users []models.User
		json.Unmarshal(w.Body.Bytes(), &users)
		assert.NotEmpty(t, users)
	})

	t.Run("List Users - No Token", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/users", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Get User by ID - Success", func(t *testing.T) {
		url := fmt.Sprintf("/api/v1/users/%s", user.ID)
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var foundUser models.User
		json.Unmarshal(w.Body.Bytes(), &foundUser)
		assert.Equal(t, user.ID, foundUser.ID)
	})

	t.Run("Update User - Success", func(t *testing.T) {
		url := fmt.Sprintf("/api/v1/users/%s", user.ID)
		updateData := map[string]string{"name": "Updated Name"}
		body, _ := json.Marshal(updateData)
		req, _ := http.NewRequest("PUT", url, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Delete User - Success", func(t *testing.T) {
		url := fmt.Sprintf("/api/v1/users/%s", user.ID)
		req, _ := http.NewRequest("DELETE", url, nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// Verify user is deleted
		reqGet, _ := http.NewRequest("GET", url, nil)
		reqGet.Header.Set("Authorization", "Bearer "+token)
		wGet := httptest.NewRecorder()
		router.ServeHTTP(wGet, reqGet)
		assert.Equal(t, http.StatusNotFound, wGet.Code)
	})
}
