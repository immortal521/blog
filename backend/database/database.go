// Package database provides a GORM-based interface for interacting with PostgreSQL databases.
// It includes transaction management and automatic migration capabilities.
package database

import (
	"context"
)

type DB interface {
	isDB()
}

type Database interface {
	DB
	Trans(fn func(TxContext) error) error
	TransWithContext(ctx context.Context, fn func(TxContext) error) error

	Close() error
	Ping() error
}

type TxContext interface {
	DB
	// GetTx() Database
	// Ctx() context.Context
}
