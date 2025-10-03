// Package database provides a GORM-based interface for interacting with PostgreSQL databases.
// It includes transaction management and automatic migration capabilities.
package database

import (
	"blog-server/internal/config"
	"context"
	"database/sql"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DB interface {
	Trans(fn func(txc *TxContext) error) error
	TransWithContext(ctx context.Context, fn func(txc *TxContext) error) error
	Conn() *gorm.DB
	BeginTx(opts *sql.TxOptions) *gorm.DB
	BeginTxWithContext(ctx context.Context, opts *sql.TxOptions) *gorm.DB

	Close() error
	Ping() error
}

type db struct {
	db *gorm.DB
}

func NewDB(cfg *config.Config) (DB, error) {
	dsn := cfg.Database.GetDSN()
	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:         logger.Default.LogMode(logger.Info),
		TranslateError: true,
	})
	if err != nil {
		return nil, err
	}
	sqlDB, err := gormDB.DB()
	if err != nil {
		return nil, err
	}
	if err := sqlDB.Ping(); err != nil {
		return nil, err
	}
	err = AutoMigrate(gormDB)
	if err != nil {
		return nil, err
	}
	return &db{db: gormDB}, nil
}

func (d *db) Trans(fn func(txc *TxContext) error) error {
	return d.db.Transaction(func(tx *gorm.DB) error {
		return fn(&TxContext{tx: tx})
	})
}

func (d *db) TransWithContext(ctx context.Context, fn func(txc *TxContext) error) error {
	return d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(&TxContext{tx: tx, ctx: ctx})
	})
}

func (d *db) BeginTx(opts *sql.TxOptions) *gorm.DB {
	if opts != nil {
		return d.db.Begin(opts)
	}
	return d.db.Begin()
}

func (d *db) BeginTxWithContext(ctx context.Context, opts *sql.TxOptions) *gorm.DB {
	if opts != nil {
		return d.db.WithContext(ctx).Begin(opts)
	}
	return d.db.WithContext(ctx).Begin()
}

func (d *db) Conn() *gorm.DB {
	return d.db
}

func (d *db) Close() error {
	sqlDB, err := d.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func (d *db) Ping() error {
	sqlDB, err := d.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
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
