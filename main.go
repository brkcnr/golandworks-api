package main

import (
	"log"

	"github.com/brkcnr/golandworks-api/internal/db"
	"github.com/brkcnr/golandworks-api/internal/todo"
	"github.com/brkcnr/golandworks-api/internal/transport"
)

func main() {
	d, err := db.New("postgres", "example", "postgres", "localhost", 5432)
	if err != nil {
        log.Fatal(err)
    }

	svc := todo.NewService(d)
	server := transport.NewServer(svc)

	if err := server.Serve(); err != nil {
		log.Fatal(err)
	}
}

