package database

import (
	"context"
	"database/sql"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DB interface {
	Trans(ctx context.Context, fn func(tx *gorm.DB) error) error
	Conn(ctx context.Context) *gorm.DB
	BeginTx(ctx context.Context, opts ...*sql.TxOptions) (*gorm.DB, error)
}

type db struct {
	db *gorm.DB
}

func NewDB(dsn string) (DB, error) {
	gormDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
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

func (d *db) Trans(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return d.db.WithContext(ctx).Transaction(fn)
}

func (d *db) BeginTx(ctx context.Context, opts ...*sql.TxOptions) (*gorm.DB, error) {
	return d.db.WithContext(ctx).Begin(opts...), nil
}

// Conn 普通连接（非事务场景）
func (d *db) Conn(ctx context.Context) *gorm.DB {
	return d.db.WithContext(ctx)
}
