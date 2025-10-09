package database

import (
	"context"
	"database/sql"
	"fmt"

	"goMcp/internal/database/generated/users"
	"goMcp/internal/logger"

	_ "github.com/mattn/go-sqlite3"
)

// DB wraps the database connection and provides access to generated queries
type DB struct {
	conn  *sql.DB
	Users *users.Queries
}

// NewConnection creates a new database connection
func NewConnection(driverName, dataSourceName string) (*DB, error) {
	logger.Info(logger.MsgDBConnecting,
		logger.Module(logger.ModuleDatabase),
		logger.Operation(logger.OpConnect),
		logger.String("driver", driverName))

	conn, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		logger.LogError(logger.MsgDBError,
			logger.Module(logger.ModuleDatabase),
			logger.Operation(logger.OpConnect),
			logger.String("driver", driverName),
			logger.ZapError(err))
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := conn.Ping(); err != nil {
		logger.LogError(logger.MsgDBError,
			logger.Module(logger.ModuleDatabase),
			logger.Operation("ping"),
			logger.String("driver", driverName),
			logger.ZapError(err))
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Info(logger.MsgDBConnected,
		logger.Module(logger.ModuleDatabase),
		logger.Operation(logger.OpConnect),
		logger.String("driver", driverName))

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
	logger.DebugCtx(ctx, logger.MsgDBTransaction,
		logger.Module(logger.ModuleDatabase),
		logger.Operation("begin"))

	tx, err := db.conn.BeginTx(ctx, nil)
	if err != nil {
		logger.ErrorCtx(ctx, logger.MsgDBError,
			logger.Module(logger.ModuleDatabase),
			logger.Operation("begin_tx"),
			logger.ZapError(err))
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	qtx := db.Users.WithTx(tx)
	if err := fn(qtx); err != nil {
		logger.ErrorCtx(ctx, logger.MsgDBError,
			logger.Module(logger.ModuleDatabase),
			logger.Operation("execute_tx"),
			logger.ZapError(err))
		return err
	}

	if err := tx.Commit(); err != nil {
		logger.ErrorCtx(ctx, logger.MsgDBError,
			logger.Module(logger.ModuleDatabase),
			logger.Operation("commit_tx"),
			logger.ZapError(err))
		return err
	}

	logger.DebugCtx(ctx, logger.MsgDBTransaction,
		logger.Module(logger.ModuleDatabase),
		logger.Operation("commit"))

	return nil
}
