package main

import (
	"context"
	"embed"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/ramadhantriyant/buku/internal/database"
	"github.com/ramadhantriyant/buku/internal/models"
)

const (
	dataDir          = "data"
	dbConnectionPool = 1
	dbFileName       = "buku.db"
	shutdownTimeout  = 10 * time.Second
	port             = ":8080"
)

//go:embed sql/schema/*.sql
var embedMigrations embed.FS

func main() {
	dbPath := filepath.Join(dataDir, dbFileName)
	db, err := connectDatabase(dbPath)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Println("Running migrations...")
	if err := runMigrations(db); err != nil {
		log.Fatalf("Error running migrations: %v", err)
	}

	log.Println("Migrations completed")
	log.Printf("Database initialized on %s\n", dbPath)

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variables is required")
	}

	config := &models.Config{
		DB:        db,
		Queries:   database.New(db),
		JWTSecret: jwtSecret,
	}

	server := createServer(config, port)

	if err := runServer(context.Background(), server, shutdownTimeout); err != nil {
		log.Fatalf("Cannot bind to 0.0.0.0%s. Error: %v", port, err)
	}
}
