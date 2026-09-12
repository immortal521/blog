package contextx

import (
	"context"

	"blog-server/entity"
)

type User struct {
	ID   uint
	Role entity.UserRole
}

type contextKey string

const ContextKeyUser contextKey = "user"

func SetUser(ctx context.Context, user User) context.Context {
	return context.WithValue(ctx, ContextKeyUser, user)
}

func GetUser(ctx context.Context) (User, bool) {
	user, ok := ctx.Value(ContextKeyUser).(User)
	return user, ok
}
