package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/brkcnr/golandworks-api/internal/apierror"
	"github.com/brkcnr/golandworks-api/internal/service"
)

// TodoItem is a todo item.
type TodoItem struct {
	Item string `json:"item"`
}

// Handler is a HTTP handler.
type Handler struct {
	todoSvc *service.TodoService
	logger  *log.Logger
}

// Option is a function that configures a Handler.
type Option func(*Handler)

// WithTodoService sets the todo service.
func WithTodoService(svc *service.TodoService) Option {
	return func(h *Handler) {
		h.todoSvc = svc
	}
}

// WithLogger sets the logger.
func WithLogger(logger *log.Logger) Option {
	return func(h *Handler) {
		h.logger = logger
	}
}

// New creates a new HTTP handler with the given options.
func New(opts ...Option) *Handler {
	handler := &Handler{
		logger: log.Default(), // Set default logger
	}

	// Apply all options
	for _, opt := range opts {
		opt(handler)
	}

	return handler
}

// ListTodos lists all todos.
func (h *Handler) ListTodos(resp http.ResponseWriter, req *http.Request) {
	todoItems, err := h.todoSvc.ListTodos(req.Context())
	if err != nil {
		h.logger.Println(err)
		resp.WriteHeader(http.StatusInternalServerError)

		return
	}
	jsonBytes, err := json.Marshal(todoItems)
	if err != nil {
		h.logger.Println(err)
		resp.WriteHeader(http.StatusInternalServerError)

		return
	}
	if _, err = resp.Write(jsonBytes); err != nil {
		h.logger.Println(err)
	}
}

// Add adds a todo.
func (h *Handler) Add(resp http.ResponseWriter, req *http.Request) {
	var todoItem TodoItem
	if err := json.NewDecoder(req.Body).Decode(&todoItem); err != nil {
		h.handleError(resp, apierror.Wrap(err, http.StatusBadRequest, "invalid JSON request"))

		return
	}

	if err := h.todoSvc.Add(req.Context(), todoItem.Item); err != nil {
		h.handleError(resp, err)

		return
	}

	resp.WriteHeader(http.StatusCreated)
}

// Search searches for todos that contain the query.
func (h *Handler) Search(resp http.ResponseWriter, req *http.Request) {
	query := req.URL.Query().Get("q")
	if query == "" {
		resp.WriteHeader(http.StatusBadRequest)

		return
	}
	results, err := h.todoSvc.Search(req.Context(), query)
	if err != nil {
		h.logger.Print(err.Error())
		resp.WriteHeader(http.StatusInternalServerError)

		return
	}
	jsonBytes, err := json.Marshal(results)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)

		return
	}
	if _, err = resp.Write(jsonBytes); err != nil {
		h.logger.Println(err)
	}
}

// handleError handles an error.
func (h *Handler) handleError(resp http.ResponseWriter, err error) {
	h.logger.Printf("Error: %v", err)

	var apiErr *apierror.APIError
	if errors.As(err, &apiErr) {
		resp.WriteHeader(apiErr.Code)
		if encodeErr := json.NewEncoder(resp).Encode(apiErr); encodeErr != nil {
			h.logger.Printf("Failed to encode error response: %v", encodeErr)
		}

		return
	}

	// Default to internal server error
	resp.WriteHeader(http.StatusInternalServerError)
	if encodeErr := json.NewEncoder(resp).Encode(apierror.ErrInternalServer); encodeErr != nil {
		h.logger.Printf("Failed to encode default error response: %v", encodeErr)
	}
}
