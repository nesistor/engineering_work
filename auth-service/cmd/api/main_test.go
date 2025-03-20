package main

import (
	"context"
	"os"
	"testing"
	"time"

	"auth/data"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)

func TestConnectToRedisSuccess(t *testing.T) {
	os.Setenv("REDIS_URL", "localhost:6379")
	defer os.Unsetenv("REDIS_URL")

	client := connectToRedis()
	assert.NotNil(t, client, "Powinien zwrócić klienta Redis")
}

func TestConnectToRedisFailure(t *testing.T) {
	os.Setenv("REDIS_URL", "invalid:1234")
	defer os.Unsetenv("REDIS_URL")

	client := connectToRedis()
	assert.Nil(t, client, "Powinien zwrócić nil przy błędnym połączeniu")
}

func TestKeyManagerInitialization(t *testing.T) {
	os.Setenv("VAULT_ADDR", "http://localhost:8200")
	os.Setenv("VAULT_TOKEN", "test-token")
	defer func() {
		os.Unsetenv("VAULT_ADDR")
		os.Unsetenv("VAULT_TOKEN")
	}()

	km, err := data.NewKeyManager(data.VaultConfig{
		Address: os.Getenv("VAULT_ADDR"),
		Token:   os.Getenv("VAULT_TOKEN"),
	}, time.Hour)
	
	assert.NoError(t, err, "Inicjalizacja KeyManager nie powinna zwracać błędów")
	assert.NotNil(t, km, "KeyManager powinien być prawidłowo zainicjalizowany")
}

func TestAppConfigInitialization(t *testing.T) {
	mockRedis := redis.NewClient(&redis.Options{})
	km := &data.KeyManager{}

	app := Config{
		RedisClient: mockRedis,
		Models:      data.New(mockRedis, km),
		KeyManager:  km,
	}

	assert.NotNil(t, app.RedisClient, "RedisClient powinien być ustawiony")
	assert.NotNil(t app.Models, "Models powinny być zainicjalizowane")
	assert.NotNil(t app.KeyManager, "KeyManager powinien być ustawiony")
}

func TestServerConfiguration(t *testing.T) {
	app := &Config{
		RedisClient: redis.NewClient(&redis.Options{}),
		KeyManager:  &data.KeyManager{},
	}

	srv := app.createServer(":8080")
	
	assert.Equal(t, 30*time.Second, srv.IdleTimeout, "IdleTimeout powinien być ustawiony")
	assert.Equal(t, 10*time.Second, srv.ReadTimeout, "ReadTimeout powinien być ustawiony")
	assert.NotNil(t, srv.Handler, "Handler powinien być zainicjalizowany")
}

// Pomocnicza funkcja do testowania konfiguracji serwera
func (app *Config) createServer(addr string) *http.Server {
	return &http.Server{
		Addr:              addr,
		Handler:           app.routes(),
		IdleTimeout:       30 * time.Second,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      5 * time.Second,
	}
}
s