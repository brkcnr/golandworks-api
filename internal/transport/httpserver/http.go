package httpserver

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/brkcnr/golandworks-api/internal/apierror"
	"github.com/brkcnr/golandworks-api/internal/handler"
	"github.com/brkcnr/golandworks-api/internal/service"
)

const (
	readTimeout       = 15 * time.Second
	writeTimeout      = 15 * time.Second
	idleTimeout       = 60 * time.Second
	readHeaderTimeout = 5 * time.Second
)

// Server is a HTTP server.
type Server struct {
	mux *http.ServeMux
}

// New creates a new HTTP server.
func New(todoSvc *service.TodoService) *Server {
	mux := http.NewServeMux()

	todoHandler := handler.New(

		handler.WithTodoService(todoSvc),

		handler.WithLogger(log.New(os.Stdout, "TODO-API: ", log.LstdFlags)),
	)

	mux.HandleFunc("GET /todo", todoHandler.ListTodos)

	mux.HandleFunc("POST /todo", todoHandler.Add)

	mux.HandleFunc("GET /search", todoHandler.Search)

	return &Server{
		mux: mux,
	}
}

// Serve serves the HTTP server.
func (s *Server) Serve() error {
	srv := &http.Server{
		Addr: ":8080",

		Handler: s.mux,

		ReadTimeout: readTimeout,

		WriteTimeout: writeTimeout,

		IdleTimeout: idleTimeout,

		ReadHeaderTimeout: readHeaderTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		return apierror.Wrap(err, http.StatusInternalServerError, "failed to start HTTP server")
	}

	return nil
}
