package repo

import (
	"context"

	entity "github.com/potom_pridumaem/internal/entity/users"
)

type (
	UserRepo interface {
		Store(ctx context.Context, user *entity.User) error
		GetByEmail(ctx context.Context, email string) (entity.User, error)
	}
)
