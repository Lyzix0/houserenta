package repo

import (
	"context"

	entity "github.com/potom_pridumaem/internal/entity/users"
)

type (
	UserRepo interface {
		Store(ctx context.Context, user *entity.User) error
		GetByEmailOrPhone(ctx context.Context, identifier string) (entity.User, error)
		GetByID(ctx context.Context, id string) (entity.User, error)
		Update(ctx context.Context, user *entity.User) error
	}

	PropertyRepo interface {
		Store(ctx context.Context, property entity.Property) error
		GetByLandlordID(ctx context.Context, landlordID string) ([]entity.Property, error)
		GetByID(ctx context.Context, id string) (entity.Property, error)
		Update(ctx context.Context, property entity.Property) error
		Delete(ctx context.Context, id string) error
		AddBalance(ctx context.Context, propertyID string, amount float64) error
	}

	LeaseRepo interface {
		GetByTenantUserID(ctx context.Context, tenantUserID string) (entity.Lease, error)
		GetByPropertyID(ctx context.Context, propertyID string) (entity.Lease, error)
		Upsert(ctx context.Context, lease entity.Lease) error
		DeleteByPropertyID(ctx context.Context, propertyID string) error
	}

	ReadingRepo interface {
		GetByPropertyID(ctx context.Context, propertyID string) ([]entity.Reading, error)
		Store(ctx context.Context, reading entity.Reading) error
	}

	BillRepo interface {
		GetByPropertyID(ctx context.Context, propertyID string) ([]entity.Bill, error)
		UpdateStatus(ctx context.Context, billID, propertyID, status string) error
	}

	CustomNextItemRepo interface {
		GetByPropertyID(ctx context.Context, propertyID string) ([]entity.CustomNextItem, error)
		Store(ctx context.Context, item entity.CustomNextItem) error
	}
)
