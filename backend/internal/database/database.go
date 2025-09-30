// Package database provides a GORM-based interface for interacting with PostgresSQL databases.
// It includes transaction management and automatic migration capabilities.
package database

import (
	"database/sql"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DB interface {
	Trans(fn func(txc *TxContext) error) error
	Conn() *gorm.DB
	BeginTx(opts ...*sql.TxOptions) (*gorm.DB, error)
}

type db struct {
	db *gorm.DB
}

func NewDB(dsn string) (DB, error) {
	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:         logger.Default.LogMode(logger.Info),
		TranslateError: true,
	})
	if err != nil {
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

func (d *db) BeginTx(opts ...*sql.TxOptions) (*gorm.DB, error) {
	return d.db.Begin(opts...), nil
}

func (d *db) Conn() *gorm.DB {
	return d.db
}

type TxContext struct {
	tx *gorm.DB
}

func (t *TxContext) GetTx() *gorm.DB {
	return t.tx
}
