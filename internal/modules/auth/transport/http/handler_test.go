package http

import (
	"bytes"
	"encoding/json"
	authDomain "english-learning/internal/modules/auth/domain"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// setupRouter creates a gin engine with the handler registered for testing.
func setupRouter(h *AuthHandler) *gin.Engine {
	r := gin.New()
	r.POST("/auth/register", h.Register)
	r.POST("/auth/login", h.Login)
	r.POST("/auth/refresh-token", h.RefreshToken)
	r.POST("/auth/logout", h.Logout)
	return r
}

func performRequest(r *gin.Engine, method, path string, body interface{}) *httptest.ResponseRecorder {
	jsonBytes, _ := json.Marshal(body)
	req, _ := http.NewRequest(method, path, bytes.NewBuffer(jsonBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

// --- Register Handler Tests ---

func TestRegisterHandler_Success(t *testing.T) {
	t.Parallel()
	mockService := new(MockAuthService)
	handler := NewAuthHandler(mockService)
	router := setupRouter(handler)

	mockService.On("Register", mock.AnythingOfType("*domain.RegisterRequest")).Return(nil)

	body := RegisterRequestDTO{
		Email:    "test@example.com",
		Password: "password123",
	}

	w := performRequest(router, "POST", "/auth/register", body)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockService.AssertExpectations(t)
}

func TestRegisterHandler_InvalidJSON_MissingEmail(t *testing.T) {
	t.Parallel()
	mockService := new(MockAuthService)
	handler := NewAuthHandler(mockService)
	router := setupRouter(handler)

	body := map[string]string{
		"password": "password123",
	}

	w := performRequest(router, "POST", "/auth/register", body)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockService.AssertNotCalled(t, "Register", mock.Anything)
}

func TestRegisterHandler_InvalidJSON_ShortPassword(t *testing.T) {
	t.Parallel()
	mockService := new(MockAuthService)
	handler := NewAuthHandler(mockService)
	router := setupRouter(handler)

	body := RegisterRequestDTO{
		Email:    "test@example.com",
		Password: "short",
	}

	w := performRequest(router, "POST", "/auth/register", body)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockService.AssertNotCalled(t, "Register", mock.Anything)
}

func TestRegisterHandler_InvalidEmail(t *testing.T) {
	t.Parallel()
	mockService := new(MockAuthService)
	handler := NewAuthHandler(mockService)
	router := setupRouter(handler)

	body := RegisterRequestDTO{
		Email:    "not-an-email",
		Password: "password123",
	}

	w := performRequest(router, "POST", "/auth/register", body)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockService.AssertNotCalled(t, "Register", mock.Anything)
}

func TestRegisterHandler_ServiceError(t *testing.T) {
	t.Parallel()
	mockService := new(MockAuthService)
	handler := NewAuthHandler(mockService)
	router := setupRouter(handler)

	mockService.On("Register", mock.AnythingOfType("*domain.RegisterRequest")).Return(errors.New("email already registered"))

	body := RegisterRequestDTO{
		Email:    "test@example.com",
		Password: "password123",
	}

	w := performRequest(router, "POST", "/auth/register", body)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "email already registered", resp["message"])
}

// --- Login Handler Tests ---

func TestLoginHandler_Success(t *testing.T) {
	t.Parallel()
	mockService := new(MockAuthService)
	handler := NewAuthHandler(mockService)
	router := setupRouter(handler)

	tokenPair := &authDomain.TokenPair{
		AccessToken:  "access-token",
		RefreshToken: "refresh-token",
	}

	mockService.On("Login", mock.AnythingOfType("*domain.LoginRequest"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(tokenPair, nil)

	body := LoginRequestDTO{
		Email:    "test@example.com",
		Password: "password123",
	}

	w := performRequest(router, "POST", "/auth/login", body)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)

	data := resp["data"].(map[string]interface{})
	assert.Equal(t, "access-token", data["accessToken"])
	assert.Equal(t, "refresh-token", data["refreshToken"])
}

func TestLoginHandler_InvalidJSON(t *testing.T) {
	t.Parallel()
	mockService := new(MockAuthService)
	handler := NewAuthHandler(mockService)
	router := setupRouter(handler)

	body := map[string]string{
		"email": "test@example.com",
		// missing password
	}

	w := performRequest(router, "POST", "/auth/login", body)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockService.AssertNotCalled(t, "Login", mock.Anything, mock.Anything, mock.Anything)
}

func TestLoginHandler_Unauthorized(t *testing.T) {
	t.Parallel()
	mockService := new(MockAuthService)
	handler := NewAuthHandler(mockService)
	router := setupRouter(handler)

	mockService.On("Login", mock.AnythingOfType("*domain.LoginRequest"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil, errors.New("invalid credentials"))

	body := LoginRequestDTO{
		Email:    "test@example.com",
		Password: "wrong-password",
	}

	w := performRequest(router, "POST", "/auth/login", body)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// --- RefreshToken Handler Tests ---

func TestRefreshTokenHandler_Success(t *testing.T) {
	t.Parallel()
	mockService := new(MockAuthService)
	handler := NewAuthHandler(mockService)
	router := setupRouter(handler)

	tokenPair := &authDomain.TokenPair{
		AccessToken:  "new-access-token",
		RefreshToken: "new-refresh-token",
	}

	mockService.On("RefreshToken", "valid-refresh-token").Return(tokenPair, nil)

	body := RefreshTokenRequestDTO{
		RefreshToken: "valid-refresh-token",
	}

	w := performRequest(router, "POST", "/auth/refresh-token", body)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)

	data := resp["data"].(map[string]interface{})
	assert.Equal(t, "new-access-token", data["accessToken"])
	assert.Equal(t, "new-refresh-token", data["refreshToken"])
}

func TestRefreshTokenHandler_InvalidJSON(t *testing.T) {
	t.Parallel()
	mockService := new(MockAuthService)
	handler := NewAuthHandler(mockService)
	router := setupRouter(handler)

	body := map[string]string{}

	w := performRequest(router, "POST", "/auth/refresh-token", body)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockService.AssertNotCalled(t, "RefreshToken", mock.Anything)
}

func TestRefreshTokenHandler_Unauthorized(t *testing.T) {
	t.Parallel()
	mockService := new(MockAuthService)
	handler := NewAuthHandler(mockService)
	router := setupRouter(handler)

	mockService.On("RefreshToken", "invalid-token").Return(nil, errors.New("invalid refresh token"))

	body := RefreshTokenRequestDTO{
		RefreshToken: "invalid-token",
	}

	w := performRequest(router, "POST", "/auth/refresh-token", body)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// --- Logout Handler Tests ---

func TestLogoutHandler_Success(t *testing.T) {
	t.Parallel()
	mockService := new(MockAuthService)
	handler := NewAuthHandler(mockService)
	router := setupRouter(handler)

	mockService.On("Logout", "valid-refresh-token").Return(nil)

	body := RefreshTokenRequestDTO{
		RefreshToken: "valid-refresh-token",
	}

	w := performRequest(router, "POST", "/auth/logout", body)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestLogoutHandler_InvalidJSON(t *testing.T) {
	t.Parallel()
	mockService := new(MockAuthService)
	handler := NewAuthHandler(mockService)
	router := setupRouter(handler)

	body := map[string]string{}

	w := performRequest(router, "POST", "/auth/logout", body)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockService.AssertNotCalled(t, "Logout", mock.Anything)
}

func TestLogoutHandler_ServiceError(t *testing.T) {
	t.Parallel()
	mockService := new(MockAuthService)
	handler := NewAuthHandler(mockService)
	router := setupRouter(handler)

	mockService.On("Logout", "some-token").Return(errors.New("revoke failed"))

	body := RefreshTokenRequestDTO{
		RefreshToken: "some-token",
	}

	w := performRequest(router, "POST", "/auth/logout", body)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
