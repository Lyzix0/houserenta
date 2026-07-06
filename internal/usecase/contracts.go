package usecase

import (
	"context"

	entity "github.com/potom_pridumaem/internal/entity/users"
)

type (
	User interface {
		Register(ctx context.Context, name, email, password, role, document, phone string) (entity.User, error)
		Login(ctx context.Context, email, password string) (entity.User, error)
	}
)
