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
		GetUnlinkedTenants(ctx context.Context) ([]entity.User, error)
	}

	PropertyRepo interface {
		Store(ctx context.Context, property entity.Property) error
		GetByLandlordID(ctx context.Context, landlordID string) ([]entity.Property, error)
		GetByID(ctx context.Context, id string) (entity.Property, error)
		Update(ctx context.Context, property entity.Property) error
		Delete(ctx context.Context, id string) error
		AddBalance(ctx context.Context, propertyID string, amount float64) error
		GetVacant(ctx context.Context) ([]entity.Property, error)
	}

	LeaseRepo interface {
		GetByTenantUserID(ctx context.Context, tenantUserID string) (entity.Lease, error)
		GetByPropertyID(ctx context.Context, propertyID string) (entity.Lease, error)
		GetAll(ctx context.Context) ([]entity.Lease, error)
		Upsert(ctx context.Context, lease entity.Lease) error
		DeleteByPropertyID(ctx context.Context, propertyID string) error
	}

	ReadingRepo interface {
		GetByPropertyID(ctx context.Context, propertyID string) ([]entity.Reading, error)
		Store(ctx context.Context, reading entity.Reading) error
		GetOldestUnaccounted(ctx context.Context, propertyID string) (entity.Reading, error)
		GetLatestAccounted(ctx context.Context, propertyID string) (entity.Reading, error)
		MarkAccounted(ctx context.Context, id string) error
	}

	BillRepo interface {
		GetByPropertyID(ctx context.Context, propertyID string) ([]entity.Bill, error)
		UpdateStatus(ctx context.Context, billID, propertyID, status string) error
		GetLastByPropertyIDAndType(ctx context.Context, propertyID, billType string) (entity.Bill, error)
		Store(ctx context.Context, bill entity.Bill) error
	}

	CustomNextItemRepo interface {
		GetByPropertyID(ctx context.Context, propertyID string) ([]entity.CustomNextItem, error)
		Store(ctx context.Context, item entity.CustomNextItem) error
		DeleteByPropertyID(ctx context.Context, propertyID string) error
	}

	ApplicationRepo interface {
		Store(ctx context.Context, application entity.Application) error
		GetByPropertyID(ctx context.Context, propertyID string) ([]entity.Application, error)
		DeleteByPropertyID(ctx context.Context, propertyID string) error
	}
)
