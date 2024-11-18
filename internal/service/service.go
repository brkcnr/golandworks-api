package service

import (
	"context"
	"net/http"
	"strings"

	"github.com/brkcnr/golandworks-api/internal/apierror"
	"github.com/brkcnr/golandworks-api/internal/db"
)

// TodoService handles todo business logic.
type TodoService struct {
	db db.Storer
}

// Option is a function that configures a TodoService.
type Option func(*TodoService)

// WithDB sets the database for the service.
func WithDB(store db.Storer) Option {
	return func(s *TodoService) {
		s.db = store
	}
}

// New creates a new TodoService with the given options.
func New(opts ...Option) *TodoService {
	svc := &TodoService{}
	for _, opt := range opts {
		opt(svc)
	}

	return svc
}

// Add creates a new todo item.
func (s *TodoService) Add(ctx context.Context, todo string) error {
	if todo == "" {
		return apierror.Wrap(
			apierror.ErrInvalidRequest,
			http.StatusBadRequest,
			"todo item cannot be empty",
		)
	}

	items, err := s.ListTodos(ctx)
	if err != nil {
		return apierror.Wrap(err, http.StatusInternalServerError, "failed to check for duplicates")
	}

	for _, t := range items {
		if t.Task == todo {
			return apierror.ErrDuplicateTodo
		}
	}

	if insertErr := s.db.InsertItem(ctx, db.Item{
		Task: todo,

		Status: "TO_BE_STARTED",
	}); insertErr != nil {
		return apierror.Wrap(insertErr, http.StatusInternalServerError, "failed to insert todo item")
	}

	return nil
}

// Search finds todos containing the query string.
func (s *TodoService) Search(ctx context.Context, query string) ([]string, error) {
	items, err := s.ListTodos(ctx)
	if err != nil {
		return nil, apierror.Wrap(err, http.StatusInternalServerError, "failed to list todos")
	}

	var results []string
	for _, todo := range items {
		if strings.Contains(strings.ToLower(todo.Task), strings.ToLower(query)) {
			results = append(results, todo.Task)
		}
	}

	return results, nil
}

// ListTodos lists all todo items.
func (s *TodoService) ListTodos(ctx context.Context) ([]db.Item, error) {
	items, err := s.db.GetAllItems(ctx)
	if err != nil {
		return nil, apierror.Wrap(err, http.StatusInternalServerError, "failed to get todos from database")
	}

	return items, nil
}
