package database

import (
	"context"

	"blog-server/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type db struct {
	db *gorm.DB
}

func New(cfg *config.Config) (DB, error) {
	dsn := cfg.Database.GetDSN()
	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:                                   logger.Default.LogMode(logger.Info),
		TranslateError:                           true,
		DisableForeignKeyConstraintWhenMigrating: true,
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
