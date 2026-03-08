package database

import (
	"context"

	"blog-server/config"
	"blog-server/errs"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type GormDatabase struct {
	db *gorm.DB
}

type GormTxContext struct {
	db *gorm.DB
}

func (d *GormDatabase) isDB()  {}
func (t *GormTxContext) isDB() {}

func (d *GormDatabase) Trans(fn func(TxContext) error) error {
	return d.db.Transaction(func(tx *gorm.DB) error {
		return fn(&GormTxContext{db: tx})
	})
}

func (d *GormDatabase) TransWithContext(ctx context.Context, fn func(txc TxContext) error) error {
	return d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(&GormTxContext{db: tx})
	})
}

func (d *GormDatabase) Close() error {
	sqlDB, err := d.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func (d *GormDatabase) Ping() error {
	sqlDB, err := d.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}

func NewGormDatabase(cfg *config.Config) (Database, error) {
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
	return &GormDatabase{db: gormDB}, nil
}

func NewGormTxContext(tx *gorm.DB) DB {
	return &GormTxContext{db: tx}
}

func ToGormDB(db DB) (*gorm.DB, error) {
	switch v := db.(type) {
	case *GormDatabase:
		return v.db, nil
	case *GormTxContext:
		return v.db, nil
	default:
		return nil, errs.New(errs.CodeInternalError, "Invalid database type, please use GormDatabase or GormTxContext", nil)
	}
}
