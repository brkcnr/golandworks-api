package main

import (
	"log"

	"github.com/brkcnr/golandworks-api/internal/config"
	"github.com/brkcnr/golandworks-api/internal/db"
	"github.com/brkcnr/golandworks-api/internal/service"
	"github.com/brkcnr/golandworks-api/internal/transport/httpserver"
)

// main is the entry point for the application.
func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	dbConn, err := db.New(cfg.DB)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer dbConn.Close()

	todoService := service.New(
		service.WithDB(dbConn),
	)

	server := httpserver.New(todoService)

	if err = server.Serve(); err != nil {
		log.Printf("Server error: %v", err)

		return
	}
}
