// Package repository
package repository

import (
	"context"
	"errors"

	"blog-server/database"
	"blog-server/errs"

	"gorm.io/gorm"
)

func existsBy[T any](ctx context.Context, db database.DB, field string, value any) (bool, error) {
	gdb := unwrapDB(db)
	_, err := gorm.G[*T](gdb).Where(field+" = ?", value).Take(ctx)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	if err != nil {
		return false, errs.New(errs.CodeDatabaseError, "database error", err)
	}
	return true, nil
}

func unwrapDB(db database.DB) *gorm.DB {
	gdb, err := database.ToGormDB(db)
	if err != nil {
		panic(err)
	}
	return gdb
}
