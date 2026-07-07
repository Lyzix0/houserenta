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
		GetByID(ctx context.Context, id string) (entity.User, error)
		GetByEmail(ctx context.Context, email string) (entity.User, error)
		UpdateProfile(ctx context.Context, id string, name, document, phone, email *string, paymentCard *string) error
		ChangePassword(ctx context.Context, id, oldPassword, newPassword string) error
		SwitchRole(ctx context.Context, id, targetRole string) (entity.User, error)
	}

	Property interface {
		CreateProperty(ctx context.Context, body request.Property) (entity.Property, error)
		GetProperties(ctx context.Context, landlordID string) ([]entity.Property, error)
		GetProperty(ctx context.Context, id string) (entity.Property, error)
		UpdateProperty(ctx context.Context, id string, body request.Property) error
		DeleteProperty(ctx context.Context, id string) error
		GetVacantProperties(ctx context.Context) ([]entity.Property, error)
	}
)
