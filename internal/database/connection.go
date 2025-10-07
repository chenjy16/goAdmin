package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"admin/internal/config"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	*sql.DB
	*Queries
}

func NewConnection(cfg *config.Config) (*DB, error) {
	// 确保数据目录存在
	dsn := cfg.GetDatabaseDSN()
	dir := filepath.Dir(dsn)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	db, err := sql.Open(cfg.GetDatabaseDriver(), dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	queries := New(db)

	log.Println("Database connection established successfully")

	return &DB{
		DB:      db,
		Queries: queries,
	}, nil
}

func (db *DB) Close() error {
	return db.DB.Close()
}

func (db *DB) Ping(ctx context.Context) error {
	return db.DB.PingContext(ctx)
}
