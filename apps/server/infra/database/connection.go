package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/Dufyz/scd-server/internal/env"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func NewDBConnection() (*sql.DB, error) {
	url := env.GetString("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable")
	return newDBConnection(url)
}

func NewDBConnectionWithRetries(maxRetries int) (*sql.DB, error) {
	url := env.GetString("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable")
	return newDBConnectionWithRetries(url, maxRetries)
}

func NewReplicaDBConnectionWithRetries(maxRetries int) (*sql.DB, error) {
	primaryURL := env.GetString("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable")
	url := env.GetString("DATABASE_URL_REPLICA", primaryURL)
	return newDBConnectionWithRetries(url, maxRetries)
}

func newDBConnection(url string) (*sql.DB, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(10)                 // Maximum number of open connections (Supabase limit: 15)
	db.SetMaxIdleConns(2)                  // Maximum number of idle connections
	db.SetConnMaxLifetime(5 * time.Minute) // Maximum lifetime of a connection
	db.SetConnMaxIdleTime(1 * time.Minute) // Maximum idle time before closing

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	log.Println("Database connection pool configured successfully")
	return db, nil
}

func newDBConnectionWithRetries(url string, maxRetries int) (*sql.DB, error) {
	var db *sql.DB
	var err error

	for i := range maxRetries {
		fmt.Println("DATABASE_URL", url)

		db, err = sql.Open("postgres", url)
		if err != nil {
			log.Printf("attempt %d: failed to open connection: %v", i+1, err)
			time.Sleep(2 * time.Second)
			continue
		}

		db.SetMaxOpenConns(10)                  // Maximum number of open connections (Supabase limit: 15)
		db.SetMaxIdleConns(10)                  // Keep all connections warm
		db.SetConnMaxLifetime(30 * time.Minute) // Maximum lifetime of a connection
		db.SetConnMaxIdleTime(10 * time.Minute) // Maximum idle time before closing

		log.Printf("attempt %d: pinging database...", i+1)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		err = db.PingContext(ctx)
		cancel()
		if err == nil {
			log.Println("connected to DB with connection pool configured")
			return db, nil
		}

		log.Printf("attempt %d: failed to ping DB: %v", i+1, err)
		db.Close()
		time.Sleep(2 * time.Second)
	}

	log.Fatalf("could not connect to database after %d attempts: %v", maxRetries, err)
	return nil, err
}
