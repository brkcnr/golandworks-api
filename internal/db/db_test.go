package db_test

import (
	"context"
	"os"
	"strconv"
	"testing"

	"github.com/brkcnr/golandworks-api/internal/config"
	"github.com/brkcnr/golandworks-api/internal/db"
)

func setupTestDB(t *testing.T) (*db.DB, func()) {
	// Read test database configuration from environment variables
	port := 5432 // default port
	if portStr := os.Getenv("TEST_DB_PORT"); portStr != "" {
		if p, err := strconv.Atoi(portStr); err == nil {
			port = p
		}
	}

	cfg := config.DBConfig{
		Host:     getEnvOrDefault("TEST_DB_HOST", "localhost"),
		Port:     port,
		User:     getEnvOrDefault("TEST_DB_USER", "postgres"),
		Password: getEnvOrDefault("TEST_DB_PASSWORD", "postgres"),
		DBName:   getEnvOrDefault("TEST_DB_NAME", "golandworks_test"),
	}

	database, err := db.New(cfg)
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Return cleanup function
	cleanup := func() {
		database.Close()
	}

	return database, cleanup
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func TestInsertAndGetItems(t *testing.T) {
	database, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	// Test cases
	testCases := []struct {
		name     string
		item     db.Item
		wantErr  bool
	}{
		{
			name: "Valid item",
			item: db.Item{
				Task:   "Test task",
				Status: "pending",
			},
			wantErr: false,
		},
		{
			name: "Another valid item",
			item: db.Item{
				Task:   "Another test task",
				Status: "completed",
			},
			wantErr: false,
		},
	}

	// Insert items
	for _, tc := range testCases {
		t.Run("Insert "+tc.name, func(t *testing.T) {
			err := database.InsertItem(ctx, tc.item)
			if (err != nil) != tc.wantErr {
				t.Errorf("InsertItem() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}

	// Test GetAllItems
	t.Run("GetAllItems", func(t *testing.T) {
		items, err := database.GetAllItems(ctx)
		if err != nil {
			t.Fatalf("GetAllItems() error = %v", err)
		}

		// Check if we got at least the number of items we inserted
		if len(items) < len(testCases) {
			t.Errorf("GetAllItems() got %d items, want at least %d", len(items), len(testCases))
		}

		// Verify that our test items are in the results
		itemFound := make(map[string]bool)
		for _, item := range items {
			for _, tc := range testCases {
				if item.Task == tc.item.Task && item.Status == tc.item.Status {
					itemFound[tc.name] = true
				}
			}
		}

		for _, tc := range testCases {
			if !itemFound[tc.name] {
				t.Errorf("GetAllItems() did not return expected item: %v", tc.item)
			}
		}
	})
} 