// Package repo
package repo

import (
	"context"
	"errors"

	"blog-server/pkg/errs"

	"gorm.io/gorm"
)

func ExistsBy[T any](ctx context.Context, db *gorm.DB, field string, value any) (bool, error) {
	_, err := gorm.G[*T](db).Where(field+" = ?", value).Take(ctx)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	if err != nil {
		return false, errs.New(errs.CodeDatabaseError, "database error", err)
	}
	return true, nil
}
