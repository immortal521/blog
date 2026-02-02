// Package database provides a GORM-based interface for interacting with PostgreSQL databases.
// It includes transaction management and automatic migration capabilities.
package database

import (
	"context"

	"gorm.io/gorm"
)

type DB interface {
	Trans(fn func(txc *TxContext) error) error
	TransWithContext(ctx context.Context, fn func(txc *TxContext) error) error
	Conn() *gorm.DB

	Close() error
	Ping() error
}

type TxContext struct {
	tx  *gorm.DB
	ctx context.Context
}

func (t *TxContext) GetTx() *gorm.DB {
	return t.tx
}

func (t *TxContext) Ctx() context.Context {
	return t.ctx
}
