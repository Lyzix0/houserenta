package usecase

import (
	"context"

	"github.com/potom_pridumaem/internal/controller/v1/request"
	entity "github.com/potom_pridumaem/internal/entity/users"
)

type (
	User interface {
		Register(ctx context.Context, name, email, password, role, document, phone string) (entity.User, error)
		Login(ctx context.Context, email, password string) (entity.User, error)
	}

	Property interface {
		CreateProperty(ctx context.Context, body request.Property) (entity.Property, error)
	}
)
