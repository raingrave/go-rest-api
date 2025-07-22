package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	clearUserTable()

	t.Run("Successful User Creation", func(t *testing.T) {
		userCredentials := map[string]string{
			"name":     "New User",
			"email":    "newuser@example.com",
			"password": "password123",
		}
		body, _ := json.Marshal(userCredentials)
		req, _ := http.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "New User", response["name"])
		assert.Equal(t, "newuser@example.com", response["email"])
		assert.NotContains(t, response, "password", "Password should not be in the response")
	})

	t.Run("Validation Error - Short Password", func(t *testing.T) {
		userCredentials := map[string]string{
			"name":     "Another User",
			"email":    "another@example.com",
			"password": "123",
		}
		body, _ := json.Marshal(userCredentials)
		req, _ := http.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
