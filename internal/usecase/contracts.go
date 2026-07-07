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
		Me(ctx context.Context, userID string) (UserProfile, error)
	}

	Property interface {
		CreateProperty(ctx context.Context, body request.Property) (entity.Property, error)
		GetProperties(ctx context.Context, landlordID string) ([]entity.Property, error)
		GetProperty(ctx context.Context, id, landlordID string) (entity.Property, error)
		UpdateProperty(ctx context.Context, id, landlordID string, body request.Property) error
		DeleteProperty(ctx context.Context, id, landlordID string) error
	}

	// UserProfile is the full session profile returned by the auth/me endpoint,
	// enriched with the tenant's linked property when applicable.
	UserProfile struct {
		User             entity.User
		TenantPropertyID *string
	}
)
