package repo

import (
	"context"

	entity "github.com/potom_pridumaem/internal/entity/users"
)

type (
	UserRepo interface {
		Store(ctx context.Context, user *entity.User) error
		GetByEmail(ctx context.Context, email string) (entity.User, error)
		GetByID(ctx context.Context, id string) (entity.User, error)
		Update(ctx context.Context, user *entity.User) error
		UpdatePassword(ctx context.Context, userID, hash string) error
	}

	PropertyRepo interface {
		Store(ctx context.Context, property entity.Property) error
		GetByLandlordID(ctx context.Context, landlordID string) ([]entity.Property, error)
		GetByID(ctx context.Context, id string) (entity.Property, error)
		Update(ctx context.Context, prop *entity.Property) error
		Delete(ctx context.Context, id string) error
		GetVacant(ctx context.Context) ([]entity.Property, error)
	}
)
