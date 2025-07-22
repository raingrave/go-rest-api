package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/raingrave/apirest/internal/models"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestLogin(t *testing.T) {
	clearUserTable()

	password := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := models.User{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: string(hashedPassword),
	}
	_, err := userRepo.CreateUser(user)
	assert.NoError(t, err)

	t.Run("Successful Login", func(t *testing.T) {
		loginCredentials := map[string]string{
			"email":    "test@example.com",
			"password": "password123",
		}
		body, _ := json.Marshal(loginCredentials)
		req, _ := http.NewRequest("POST", "/api/v1/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response, "token", "Response should contain a token")
	})

	t.Run("Invalid Password", func(t *testing.T) {
		loginCredentials := map[string]string{
			"email":    "test@example.com",
			"password": "wrongpassword",
		}
		body, _ := json.Marshal(loginCredentials)
		req, _ := http.NewRequest("POST", "/api/v1/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Validation Error", func(t *testing.T) {
		loginCredentials := map[string]string{
			"email": "not-an-email",
		}
		body, _ := json.Marshal(loginCredentials)
		req, _ := http.NewRequest("POST", "/api/v1/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
