package data_test

import (
	"context"
	"crypto/rsa"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/yourproject/auth-service/data"
	"github.com/golang-jwt/jwt/v5"
)

var (
	testModels       data.Models
	testCtx          = context.Background()
	validPrivateKey  *rsa.PrivateKey
	invalidPrivateKey *rsa.PrivateKey
)

func TestMain(m *testing.M) {
	// Konfiguracja testowego Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	// Generowanie kluczy RSA dla testów
	km := data.NewKeyManager()
	if err := km.GenerateKeyPair("test_kid"); err != nil {
		panic(fmt.Sprintf("Failed to generate test keys: %v", err))
	}
	if err := km.GenerateKeyPair("invalid_kid"); err != nil {
		panic(fmt.Sprintf("Failed to generate invalid keys: %v", err))
	}

	testModels = data.New(redisClient, km)

	exitCode := m.Run()
	redisClient.Close()
	os.Exit(exitCode)
}

func TestTokenValidation(t *testing.T) {
	// Test poprawnego tokenu
	t.Run("Poprawny token z prawidłowym zakresem", func(t *testing.T) {
		token, err := testModels.Token.GenerateToken(testCtx, 123, data.RoleUser, time.Hour, data.ScopeAuthentication, "test_kid")
		if err != nil {
			t.Fatalf("Błąd generowania tokenu: %v", err)
		}

		userID, role, err := testModels.Token.GetUserIDForToken(testCtx, token, data.ScopeAuthentication)
		if err != nil {
			t.Errorf("Błąd walidacji poprawnego tokenu: %v", err)
		}

		if userID != 123 || role != data.RoleUser {
			t.Errorf("Nieprawidłowe dane użytkownika. Oczekiwano (123, user), otrzymano (%d, %s)", userID, role)
		}
	})

	// Test tokenu z nieprawidłowym zakresem
	t.Run("Token z nieprawidłowym zakresem", func(t *testing.T) {
		token, err := testModels.Token.GenerateToken(testCtx, 123, data.RoleUser, time.Hour, data.ScopeAuthentication, "test_kid")
		if err != nil {
			t.Fatalf("Błąd generowania tokenu: %v", err)
		}

		_, _, err = testModels.Token.GetUserIDForToken(testCtx, token, data.ScopeRefresh)
		if err == nil {
			t.Error("Oczekiwano błędu dla nieprawidłowego zakresu, ale go nie otrzymano")
		}
	})

	// Test przedawnionego tokenu
	t.Run("Przedawniony token", func(t *testing.T) {
		token, err := testModels.Token.GenerateToken(testCtx, 123, data.RoleUser, time.Millisecond, data.ScopeAuthentication, "test_kid")
		if err != nil {
			t.Fatalf("Błąd generowania tokenu: %v", err)
		}

		time.Sleep(2 * time.Millisecond)

		_, _, err = testModels.Token.GetUserIDForToken(testCtx, token, data.ScopeAuthentication)
		if err == nil {
			t.Error("Oczekiwano błędu dla przedawnionego tokenu, ale go nie otrzymano")
		}
	})

	// Test dezaktywowanego tokenu
	t.Run("Dezaktywowany token", func(t *testing.T) {
		token, err := testModels.Token.GenerateToken(testCtx, 123, data.RoleUser, time.Hour, data.ScopeAuthentication, "test_kid")
		if err != nil {
			t.Fatalf("Błąd generowania tokenu: %v", err)
		}

		if err := testModels.Token.InsertDeactivatedToken(testCtx, token, time.Hour); err != nil {
			t.Fatalf("Błąd dezaktywacji tokenu: %v", err)
		}
		defer testModels.Token.RedisClient.Del(testCtx, "deactivated_token:"+token)

		_, _, err = testModels.Token.GetUserIDForToken(testCtx, token, data.ScopeAuthentication)
		if err == nil {
			t.Error("Oczekiwano błędu dla dezaktywowanego tokenu, ale go nie otrzymano")
		}
	})

	// Test tokenu z nieprawidłowym podpisem
	t.Run("Token z nieprawidłowym podpisem", func(t *testing.T) {
		token, err := testModels.Token.GenerateToken(testCtx, 123, data.RoleUser, time.Hour, data.ScopeAuthentication, "invalid_kid")
		if err != nil {
			t.Fatalf("Błąd generowania tokenu: %v", err)
		}

		_, _, err = testModels.Token.GetUserIDForToken(testCtx, token, data.ScopeAuthentication)
		if err == nil {
			t.Error("Oczekiwano błędu dla tokenu z nieprawidłowym podpisem, ale go nie otrzymano")
		}
	})
}

func TestRefreshTokenFlow(t *testing.T) {
	refreshToken, err := testModels.Token.GenerateToken(testCtx, 456, data.RoleAdmin, time.Hour, data.ScopeRefresh, "test_kid")
	if err != nil {
		t.Fatalf("Błąd generowania refresh tokenu: %v", err)
	}

	newAccessToken, err := testModels.Token.RefreshAccessToken(testCtx, refreshToken, "test_kid")
	if err != nil {
		t.Fatalf("Błąd odświeżania tokenu: %v", err)
	}

	userID, role, err := testModels.Token.GetUserIDForToken(testCtx, newAccessToken, data.ScopeAuthentication)
	if err != nil {
		t.Errorf("Błąd walidacji nowego tokenu: %v", err)
	}

	if userID != 456 || role != data.RoleAdmin {
		t.Errorf("Nieprawidłowe dane użytkownika. Oczekiwano (456, admin), otrzymano (%d, %s)", userID, role)
	}
}