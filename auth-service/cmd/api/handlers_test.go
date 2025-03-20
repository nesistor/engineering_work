package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"auth/data"
	"auth/users"
	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

// Mock dla UserServiceClient
type MockUserServiceClient struct {
	mock.Mock
}

func (m *MockUserServiceClient) ValidateUser(ctx context.Context, in *users.ValidateUserRequest, opts ...grpc.CallOption) (*users.ValidateUserResponse, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*users.ValidateUserResponse), args.Error(1)
}

// Mock dla AdminServiceClient
type MockAdminServiceClient struct {
	mock.Mock
}

func (m *MockAdminServiceClient) ValidateAdmin(ctx context.Context, in *admin.ValidateAdminRequest, opts ...grpc.CallOption) (*admin.ValidateAdminResponse, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*admin.ValidateAdminResponse), args.Error(1)
}

// Dodaj mock dla TokenManager
type MockTokenManager struct {
	mock.Mock
}

func (m *MockTokenManager) GenerateToken(ctx context.Context, userID int, role data.Role, duration time.Duration, scope data.Scope, keyID string) (string, error) {
	args := m.Called(ctx, userID, role, duration, scope, keyID)
	return args.String(0), args.Error(1)
}

// Ulepszony test z pełnym mockowaniem
func TestAuthenticateUserSuccess(t *testing.T) {
	t.Parallel()
	
	// Konfiguracja mocków
	redisClient, mockRedis := redismock.NewClientMock()
	mockUserClient := new(MockUserServiceClient)
	mockTokenManager := new(MockTokenManager)
	
	// Konfiguracja handlera
	handler := &authHandler{
		grpcClient: &mockGRPCProvider{userClient: mockUserClient},
		tokenManager: mockTokenManager,
		logger: &mockLogger{},
	}
	
	// Oczekiwania
	mockUserClient.On("ValidateUser", mock.Anything, mock.Anything).
		Return(&users.ValidateUserResponse{IsValid: true, UserId: 1}, nil)
		
	mockTokenManager.On("GenerateToken", mock.Anything, 1, data.RoleUser, 
		15*time.Minute, data.ScopeAuthentication, "user-key").
		Return("access-token", nil)
		
	mockTokenManager.On("GenerateToken", mock.Anything, 1, data.RoleUser, 
		7*24*time.Hour, data.ScopeRefresh, "user-key").
		Return("refresh-token", nil)
	
	// Test
	req := httptest.NewRequest("POST", "/api/auth/login", 
		bytes.NewBufferString(`{"email":"test@example.com","password":"password"}`))
	rr := httptest.NewRecorder()
	
	handler.AuthenticateUser(rr, req)
	
	// Asercje
	assert.Equal(t, http.StatusOK, rr.Code)
	
	var response jsonResponse
	json.Unmarshal(rr.Body.Bytes(), &response)
	
	assert.False(t, response.Error)
	assert.Equal(t, "Authentication successful", response.Message)
	assert.NotEmpty(t, response.Data.(map[string]interface{})["access_token"])
	
	mockUserClient.AssertExpectations(t)
	mockTokenManager.AssertExpectations(t)
	mockRedis.ExpectationsWereMet()
}

// TestAuthenticateUserInvalidCredentials testuje nieprawidłowe dane logowania
func TestAuthenticateUserInvalidCredentials(t *testing.T) {
	mockUserClient := new(MockUserServiceClient)
	app := &Config{}

	mockUserClient.On("ValidateUser", mock.Anything, mock.Anything).Return(
		&users.ValidateUserResponse{IsValid: false}, nil,
	)

	body := bytes.NewBufferString(`{"email":"wrong@example.com","password":"wrong"}=`)
	req := httptest.NewRequest("POST", "/api/auth/login", body)
	rr := httptest.NewRecorder()

	app.AuthenticateUser(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

// TestRefreshTokenSuccess testuje pomyślne odświeżenie tokenu
func TestRefreshTokenSuccess(t *testing.T) {
	redisClient, mockRedis := redismock.NewClientMock()
	km := &data.KeyManager{}
	app := &Config{
		RedisClient: redisClient,
		Models:      data.New(redisClient, km),
	}

	// Mockowanie Redis
	mockRedis.ExpectGet("token:refresh:1").SetVal("valid-refresh-token")
	mockRedis.ExpectSetEX("token:access:1", mock.Anything, 15*time.Minute).SetVal("OK")

	body := bytes.NewBufferString(`{"refresh_token":"valid-refresh-token"}`)
	req := httptest.NewRequest("POST", "/api/auth/refresh", body)
	rr := httptest.NewRecorder()

	app.RefreshToken(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockRedis.ExpectationsWereMet()
}

// TestRevokeTokenSuccess testuje pomyślne unieważnienie tokenu
func TestRevokeTokenSuccess(t *testing.T) {
	redisClient, mockRedis := redismock.NewClientMock()
	app := &Config{
		RedisClient: redisClient,
		Models:      data.New(redisClient, &data.KeyManager{}),
	}

	mockRedis.ExpectSetEX("token:deactivated:test-token", "revoked", 7*24*time.Hour).SetVal("OK")

	body := bytes.NewBufferString(`{"token":"test-token"}`)
	req := httptest.NewRequest("POST", "/api/auth/revoke", body)
	rr := httptest.NewRecorder()

	app.RevokeToken(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockRedis.ExpectationsWereMet()
}

// TestLogoutSuccess testuje pomyślne wylogowanie
func TestLogoutSuccess(t *testing.T) {
	redisClient, mockRedis := redismock.NewClientMock()
	app := &Config{
		RedisClient: redisClient,
		Models:      data.New(redisClient, &data.KeyManager{}),
	}

	mockRedis.ExpectSetEX("token:deactivated:test-token", "revoked", 15*time.Minute).SetVal("OK")

	req := httptest.NewRequest("POST", "/auth/logout", nil)
	req.Header.Set("Authorization", "Bearer test-token")
	rr := httptest.NewRecorder()

	app.Logout(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockRedis.ExpectationsWereMet()
}

// TestValidateUserTokenSuccess testuje pomyślną walidację tokenu użytkownika
func TestValidateUserTokenSuccess(t *testing.T) {
	redisClient, mockRedis := redismock.NewClientMock()
	app := &Config{
		RedisClient: redisClient,
		Models:      data.New(redisClient, &data.KeyManager{}),
	}

	mockRedis.ExpectGet("token:access:1").SetVal("valid-token")

	req := httptest.NewRequest("GET", "/validate", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	rr := httptest.NewRecorder()

	app.ValidateUserToken(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockRedis.ExpectationsWereMet()
}

// TestValidateAdminTokenUnauthorized testuje brak autoryzacji admina
func TestValidateAdminTokenUnauthorized(t *testing.T) {
	redisClient, mockRedis := redismock.NewClientMock()
	app := &Config{
		RedisClient: redisClient,
		Models:      data.New(redisClient, &data.KeyManager{}),
	}

	mockRedis.ExpectGet("token:access:1").SetVal("user-token")

	req := httptest.NewRequest("GET", "/validate-admin", nil)
	req.Header.Set("Authorization", "Bearer user-token")
	rr := httptest.NewRecorder()

	app.ValidateAdminToken(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

// TestLogRequestErrorHandling testuje obsługę błędów logowania
func TestLogRequestErrorHandling(t *testing.T) {
	app := &Config{}
	req := httptest.NewRequest("POST", "/error", nil)
	rr := httptest.NewRecorder()

	app.logRequest("TestError", "error data")

	// W rzeczywistości należy mockować wywołanie HTTP do logger-service
	// Tutaj przykład prostego testu
	err := app.logRequest("TestError", "error data")
	assert.Error(t, err, "Powinien zwrócić błąd przy nieudanym logowaniu")
}

func TestAuthenticateUserTokenGenerationFailure(t *testing.T) {
	// Konfiguracja mocków zwracająca błąd
	mockTokenManager.On("GenerateToken", mock.Anything, mock.Anything, mock.Anything, 
		mock.Anything, mock.Anything, mock.Anything).
		Return("", errors.New("generation error"))
	
	// Wywołanie testu
	// ...
	
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}