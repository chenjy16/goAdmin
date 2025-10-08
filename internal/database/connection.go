package database

import (
	"context"
	"database/sql"
	"fmt"

	"admin/internal/database/generated/users"

	_ "github.com/mattn/go-sqlite3"
)

// DB wraps the database connection and provides access to generated queries
type DB struct {
	conn  *sql.DB
	Users *users.Queries
}

// NewConnection creates a new database connection
func NewConnection(driverName, dataSourceName string) (*DB, error) {
	conn, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &DB{
		conn:  conn,
		Users: users.New(conn),
	}, nil
}

// Close closes the database connection
func (db *DB) Close() error {
	return db.conn.Close()
}

// Ping verifies the database connection is still alive
func (db *DB) Ping() error {
	return db.conn.Ping()
}

// GetConnection returns the underlying sql.DB connection
func (db *DB) GetConnection() *sql.DB {
	return db.conn
}

// WithTx executes a function within a database transaction
func (db *DB) WithTx(ctx context.Context, fn func(*users.Queries) error) error {
	tx, err := db.conn.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	qtx := db.Users.WithTx(tx)
	if err := fn(qtx); err != nil {
		return err
	}

	return tx.Commit()
}
