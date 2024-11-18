package db

import (
	"context"
	"net/http"

	"github.com/brkcnr/golandworks-api/internal/apierror"
	"github.com/brkcnr/golandworks-api/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Item is a todo item.
type Item struct {
	Task   string `json:"task"`
	Status string `json:"status"`
}

// DB is a database.
type DB struct {
	pool *pgxpool.Pool
}

// Storer is a database storer.
type Storer interface {
	InsertItem(ctx context.Context, item Item) error
	GetAllItems(ctx context.Context) ([]Item, error)
}

// Compile time proof.
var _ Storer = (*DB)(nil)

// New creates a new database.
func New(cfg config.DBConfig) (*DB, error) {
	pool, err := pgxpool.New(context.Background(), cfg.ConnectionString())
	if err != nil {
		return nil, apierror.Wrap(err, http.StatusServiceUnavailable, "failed to connect to database")
	}

	if pingErr := pool.Ping(context.Background()); pingErr != nil {
		return nil, apierror.Wrap(pingErr, http.StatusServiceUnavailable, "failed to ping database")
	}

	return &DB{pool: pool}, nil
}

// InsertItem inserts a new item into the database.
func (db *DB) InsertItem(ctx context.Context, item Item) error {
	query := `INSERT INTO todo_items (task, status) VALUES ($1, $2)`
	_, err := db.pool.Exec(ctx, query, item.Task, item.Status)
	if err != nil {
		return apierror.Wrap(err, http.StatusInternalServerError, "failed to insert item into database")
	}

	return nil
}

// GetAllItems gets all items from the database.
func (db *DB) GetAllItems(ctx context.Context) ([]Item, error) {
	query := `SELECT task, status FROM todo_items`
	rows, err := db.pool.Query(ctx, query)
	if err != nil {
		return nil, apierror.Wrap(err, http.StatusInternalServerError, "failed to query database")
	}
	defer rows.Close()

	var items []Item
	for rows.Next() {
		var item Item
		if scanErr := rows.Scan(&item.Task, &item.Status); scanErr != nil {
			return nil, apierror.Wrap(scanErr, http.StatusInternalServerError, "failed to scan database row")
		}
		items = append(items, item)
	}

	if rowsErr := rows.Err(); rowsErr != nil {
		return nil, apierror.Wrap(rowsErr, http.StatusInternalServerError, "error iterating database rows")
	}

	return items, nil
}

// Close closes the database.
func (db *DB) Close() {
	db.pool.Close()
}
