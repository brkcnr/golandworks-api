package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/brkcnr/golandworks-api/internal/db"
	"github.com/brkcnr/golandworks-api/internal/handler"
	"github.com/brkcnr/golandworks-api/internal/service"
)

// MockDB implements db.Storer interface for testing
type MockDB struct {
	insertItemFunc  func(ctx context.Context, item db.Item) error
	getAllItemsFunc func(ctx context.Context) ([]db.Item, error)
}

func (m *MockDB) InsertItem(ctx context.Context, item db.Item) error {
	return m.insertItemFunc(ctx, item)
}

func (m *MockDB) GetAllItems(ctx context.Context) ([]db.Item, error) {
	return m.getAllItemsFunc(ctx)
}

func TestListTodos(t *testing.T) {
	mockDB := &MockDB{
		getAllItemsFunc: func(ctx context.Context) ([]db.Item, error) {
			return []db.Item{
				{Task: "todo1", Status: "TO_BE_STARTED"},
				{Task: "todo2", Status: "TO_BE_STARTED"},
			}, nil
		},
	}

	todoService := service.New(service.WithDB(mockDB))
	h := handler.New(
		handler.WithTodoService(todoService),
		handler.WithLogger(log.Default()),
	)

	req := httptest.NewRequest(http.MethodGet, "/todos", nil)
	w := httptest.NewRecorder()

	h.ListTodos(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var response []db.Item
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(response) != 2 {
		t.Errorf("expected 2 items, got %d", len(response))
	}
}

func TestAdd(t *testing.T) {
	mockDB := &MockDB{
		insertItemFunc: func(ctx context.Context, item db.Item) error {
			return nil
		},
		getAllItemsFunc: func(ctx context.Context) ([]db.Item, error) {
			return []db.Item{}, nil
		},
	}

	todoService := service.New(service.WithDB(mockDB))
	h := handler.New(
		handler.WithTodoService(todoService),
		handler.WithLogger(log.Default()),
	)

	body := bytes.NewBufferString(`{"item": "test todo"}`)
	req := httptest.NewRequest(http.MethodPost, "/todos", body)
	w := httptest.NewRecorder()

	h.Add(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status code %d, got %d", http.StatusCreated, w.Code)
	}
}

func TestSearch(t *testing.T) {
	mockDB := &MockDB{
		getAllItemsFunc: func(ctx context.Context) ([]db.Item, error) {
			return []db.Item{
				{Task: "test todo", Status: "TO_BE_STARTED"},
				{Task: "another task", Status: "TO_BE_STARTED"},
			}, nil
		},
	}

	todoService := service.New(service.WithDB(mockDB))
	h := handler.New(
		handler.WithTodoService(todoService),
		handler.WithLogger(log.Default()),
	)

	req := httptest.NewRequest(http.MethodGet, "/todos/search?q=test", nil)
	w := httptest.NewRecorder()

	h.Search(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var response []string
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(response) != 1 || response[0] != "test todo" {
		t.Errorf("unexpected response: %v", response)
	}
}

func TestSearch_EmptyQuery(t *testing.T) {
	mockDB := &MockDB{
		getAllItemsFunc: func(ctx context.Context) ([]db.Item, error) {
			return []db.Item{}, nil
		},
	}

	todoService := service.New(service.WithDB(mockDB))
	h := handler.New(
		handler.WithTodoService(todoService),
		handler.WithLogger(log.Default()),
	)

	req := httptest.NewRequest(http.MethodGet, "/todos/search", nil)
	w := httptest.NewRecorder()

	h.Search(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
}
