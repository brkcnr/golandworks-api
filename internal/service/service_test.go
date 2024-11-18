package service_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/brkcnr/golandworks-api/internal/db"
	"github.com/brkcnr/golandworks-api/internal/service"
)

// mockDB implements db.Storer interface for testing
type mockDB struct {
	items []db.Item
	err   error
}

func (m *mockDB) InsertItem(_ context.Context, item db.Item) error {
	if m.err != nil {
		return m.err
	}
	m.items = append(m.items, item)
	return nil
}

func (m *mockDB) GetAllItems(_ context.Context) ([]db.Item, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.items, nil
}

func TestNew(t *testing.T) {
	mock := &mockDB{}
	svc := service.New(service.WithDB(mock))
	if svc == nil {
		t.Error("New() returned nil service")
	}
}

func TestTodoService_Add(t *testing.T) {
	tests := []struct {
		name    string
		todo    string
		dbItems []db.Item
		dbErr   error
		wantErr bool
	}{
		{
			name:    "valid todo",
			todo:    "test todo",
			wantErr: false,
		},
		{
			name:    "empty todo",
			todo:    "",
			wantErr: true,
		},
		{
			name: "duplicate todo",
			todo: "existing todo",
			dbItems: []db.Item{
				{Task: "existing todo"},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockDB{items: tt.dbItems, err: tt.dbErr}
			svc := service.New(service.WithDB(mock))

			err := svc.Add(context.Background(), tt.todo)
			if (err != nil) != tt.wantErr {
				t.Errorf("Add() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTodoService_Search(t *testing.T) {
	mockItems := []db.Item{
		{Task: "Buy groceries"},
		{Task: "Do laundry"},
		{Task: "Buy new shoes"},
	}

	tests := []struct {
		name      string
		query     string
		dbItems   []db.Item
		dbErr     error
		want      []string
		wantErr   bool
	}{
		{
			name:    "find matching items",
			query:   "buy",
			dbItems: mockItems,
			want:    []string{"Buy groceries", "Buy new shoes"},
			wantErr: false,
		},
		{
			name:    "no matches",
			query:   "nonexistent",
			dbItems: mockItems,
			want:    []string{},
			wantErr: false,
		},
		{
			name:    "empty query",
			query:   "",
			dbItems: mockItems,
			want:    []string{"Buy groceries", "Do laundry", "Buy new shoes"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockDB{items: tt.dbItems, err: tt.dbErr}
			svc := service.New(service.WithDB(mock))

			got, err := svc.Search(context.Background(), tt.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("Search() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(got) != len(tt.want) {
				t.Errorf("Search() got = %v, want %v", got, tt.want)
				return
			}

			for i, v := range got {
				if v != tt.want[i] {
					t.Errorf("Search() got[%d] = %v, want[%d] = %v", i, v, i, tt.want[i])
				}
			}
		})
	}
}

func TestTodoService_ListTodos(t *testing.T) {
	mockItems := []db.Item{
		{Task: "Task 1"},
		{Task: "Task 2"},
	}

	tests := []struct {
		name    string
		dbItems []db.Item
		dbErr   error
		want    []db.Item
		wantErr bool
	}{
		{
			name:    "successful list",
			dbItems: mockItems,
			want:    mockItems,
			wantErr: false,
		},
		{
			name:    "empty list",
			dbItems: []db.Item{},
			want:    []db.Item{},
			wantErr: false,
		},
		{
			name:    "database error",
			dbErr:   fmt.Errorf("database error"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockDB{items: tt.dbItems, err: tt.dbErr}
			svc := service.New(service.WithDB(mock))

			got, err := svc.ListTodos(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("ListTodos() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if len(got) != len(tt.want) {
					t.Errorf("ListTodos() got = %v, want %v", got, tt.want)
				}
			}
		})
	}
} 