package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	"user-service/data"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const (
	webPort  = "80"
	gRPCPort = "50001"
)

var counts int64

// Config i a structure of a application
type Config struct {
	DB     *sql.DB
	Models data.Models
	KeyManager  *data.KeyManager 
}

func main() {
	log.Println("Starting user service")

	vaultConfig := data.VaultConfig{
		Address: os.Getenv("VAULT_ADDR"),
		Token:   os.Getenv("VAULT_TOKEN"),
	}

	keyManager, err := data.NewKeyManager(vaultConfig, time.Hour)
	if err != nil {
		log.Fatalf("failed to initialize key manager: %v", err)
	}

	// connect to DB
	conn := connectToDB()
	if conn == nil {
		log.Panic("Can't connect to Postgres!")
	}

	// set up config
	app := Config{
		DB:     conn,
		Models: data.New(conn),
		KeyManager:  keyManager,
	}

	// Start HTTP server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			log.Panic(err)
		}
	}()

	// Start gRPC server
	go app.gRPCListen()

	// Block main thread to keep the servers running
	select {}
}

func openDB(dsn string) (*sql.DB, error) {
	log.Println("Trying to connect to the database...")
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("error pinging database: %w", err)
	}

	log.Println("Successfully connected to the database")
	return db, nil
}

func connectToDB() *sql.DB {
	dsn := os.Getenv("DSN")
	log.Println("DSN:", dsn)

	for {
		connection, err := openDB(dsn)
		if err != nil {
			fmt.Errorf("Error opening database: %w", err)

			counts++
		} else {
			log.Printf("Connected to Postgres!")
			return connection
		}

		if counts > 10 {
			log.Println(err)
			return nil
		}

		log.Println("Backing off for two seconds....")
		time.Sleep(2 * time.Second)
		continue
	}
}
