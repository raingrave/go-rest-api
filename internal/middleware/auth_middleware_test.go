package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/raingrave/apirest/configs"
	"github.com/stretchr/testify/assert"
)

func setupTestRouterWithMiddleware() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(AuthMiddleware())
	r.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	return r
}

func generateTestToken(userID uuid.UUID, secret string, expiration time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(expiration).Unix(),
	})
	return token.SignedString([]byte(secret))
}

func TestAuthMiddleware(t *testing.T) {
	// Load .env for JWT_SECRET_KEY
	// Note: This assumes tests are run from the root directory.
	// A more robust solution might be needed if run from different locations.
	configs.LoadEnvForTests()
	secret := configs.EnvJWTSecretKey()
	router := setupTestRouterWithMiddleware()
	userID := uuid.New()

	t.Run("Success - Valid Token", func(t *testing.T) {
		token, _ := generateTestToken(userID, secret, time.Hour)
		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Failure - No Authorization Header", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Failure - Malformed Header (No Bearer)", func(t *testing.T) {
		token, _ := generateTestToken(userID, secret, time.Hour)
		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Failure - Expired Token", func(t *testing.T) {
		token, _ := generateTestToken(userID, secret, -time.Hour) // Expired 1 hour ago
		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Failure - Invalid Signature", func(t *testing.T) {
		token, _ := generateTestToken(userID, "wrong-secret", time.Hour)
		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Failure - Invalid Token String", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer invalid-token-string")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}
