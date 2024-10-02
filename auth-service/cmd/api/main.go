package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"auth-service/data"

	"github.com/go-redis/redis/v8"
)

const webPort = "80"

// Config is a structure of an application
type Config struct {
	RedisClient *redis.Client
	Models      data.Models
	KeyManager  *data.KeyManager 
}

func main() {
	log.Println("Starting auth service...")

	redisClient := connectToRedis()
	if redisClient == nil {
		log.Panic("Can't connect to Redis!")
	}

	vaultConfig := data.VaultConfig{
		Address: os.Getenv("VAULT_ADDR"),
		Token:   os.Getenv("VAULT_TOKEN"),
	}

	keyManager, err := data.NewKeyManager(vaultConfig, time.Hour)
	if err != nil {
		log.Fatalf("failed to initialize key manager: %v", err)
	}

	app := Config{
		RedisClient: redisClient,
		Models:      data.New(redisClient, keyManager),
		KeyManager:  keyManager, // Przypisz KeyManager
	}

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%s", webPort),
		Handler:           app.routes(),
		IdleTimeout:       30 * time.Second,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      5 * time.Second,
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func connectToRedis() *redis.Client {
	ctx := context.Background()
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	client := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	// Test connection
	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Printf("Error connecting to Redis: %v", err)
		return nil
	}

	log.Println("Successfully connected to Redis")
	return client
}
