package property

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/potom_pridumaem/internal/controller/v1/request"
	entity "github.com/potom_pridumaem/internal/entity/users"
	"github.com/potom_pridumaem/internal/repo"
)

type UseCase struct {
	repo        repo.PropertyRepo
	leases      repo.LeaseRepo
	readings    repo.ReadingRepo
	bills       repo.BillRepo
	customItems repo.CustomNextItemRepo
	users       repo.UserRepo
}

func New(
	r repo.PropertyRepo,
	leases repo.LeaseRepo,
	readings repo.ReadingRepo,
	bills repo.BillRepo,
	customItems repo.CustomNextItemRepo,
	users repo.UserRepo,
) *UseCase {
	return &UseCase{
		repo:        r,
		leases:      leases,
		readings:    readings,
		bills:       bills,
		customItems: customItems,
		users:       users,
	}
}

func (uc *UseCase) CreateProperty(ctx context.Context, landlordID string, req request.Property) (entity.Property, error) {
	property := entity.Property{
		ID:          uuid.NewString(),
		LandlordID:  landlordID,
		Name:        req.Name,
		Coordinates: req.Coordinates,
		Country:     req.Country,
		Region:      req.Region,
		City:        req.City,
		Street:      req.Street,
		House:       req.House,
		Apartment:   req.Apartment,
		GvsTariff:   req.GvsTariff,
		HvsTariff:   req.HvsTariff,
		El1Tariff:   req.El1Tariff,
		El2Tariff:   req.El2Tariff,
		Balance:     0,
	}

	if err := uc.repo.Store(ctx, property); err != nil {
		return entity.Property{}, fmt.Errorf("PropertyUseCase - CreateProperty - uc.repo.Store: %w", err)
	}

	return property, nil
}

// GetProperties returns the caller's properties enriched with their related entities
// (readings, bills, custom charges, lease/tenant info, landlord contact).
// Landlords see every property they own; tenants see the single property they lease, if any.
func (uc *UseCase) GetProperties(ctx context.Context, userID string, role entity.Role) ([]entity.PropertyDetail, error) {
	var properties []entity.Property

	switch role {
	case entity.RoleLandlord:
		props, err := uc.repo.GetByLandlordID(ctx, userID)
		if err != nil {
			return nil, fmt.Errorf("PropertyUseCase - GetProperties - uc.repo.GetByLandlordID: %w", err)
		}
		properties = props
	case entity.RoleTenant:
		lease, err := uc.leases.GetByTenantUserID(ctx, userID)
		switch {
		case err == nil:
			prop, err := uc.repo.GetByID(ctx, lease.PropertyID)
			if err != nil {
				return nil, fmt.Errorf("PropertyUseCase - GetProperties - uc.repo.GetByID: %w", err)
			}
			properties = []entity.Property{prop}
		case errors.Is(err, repo.ErrLeaseNotFound):
			properties = nil
		default:
			return nil, fmt.Errorf("PropertyUseCase - GetProperties - uc.leases.GetByTenantUserID: %w", err)
		}
	default:
		properties = nil
	}

	details := make([]entity.PropertyDetail, 0, len(properties))
	for _, prop := range properties {
		detail, err := uc.buildPropertyDetail(ctx, prop)
		if err != nil {
			return nil, err
		}
		details = append(details, detail)
	}

	return details, nil
}

func (uc *UseCase) buildPropertyDetail(ctx context.Context, prop entity.Property) (entity.PropertyDetail, error) {
	readings, err := uc.readings.GetByPropertyID(ctx, prop.ID)
	if err != nil {
		return entity.PropertyDetail{}, fmt.Errorf("PropertyUseCase - buildPropertyDetail - uc.readings.GetByPropertyID: %w", err)
	}

	bills, err := uc.bills.GetByPropertyID(ctx, prop.ID)
	if err != nil {
		return entity.PropertyDetail{}, fmt.Errorf("PropertyUseCase - buildPropertyDetail - uc.bills.GetByPropertyID: %w", err)
	}

	customItems, err := uc.customItems.GetByPropertyID(ctx, prop.ID)
	if err != nil {
		return entity.PropertyDetail{}, fmt.Errorf("PropertyUseCase - buildPropertyDetail - uc.customItems.GetByPropertyID: %w", err)
	}

	var tenant *entity.Lease
	lease, err := uc.leases.GetByPropertyID(ctx, prop.ID)
	switch {
	case err == nil:
		tenant = &lease
	case errors.Is(err, repo.ErrLeaseNotFound):
		tenant = nil
	default:
		return entity.PropertyDetail{}, fmt.Errorf("PropertyUseCase - buildPropertyDetail - uc.leases.GetByPropertyID: %w", err)
	}

	landlord, err := uc.users.GetByID(ctx, prop.LandlordID)
	if err != nil {
		return entity.PropertyDetail{}, fmt.Errorf("PropertyUseCase - buildPropertyDetail - uc.users.GetByID: %w", err)
	}

	return entity.PropertyDetail{
		Property:        prop,
		Readings:        readings,
		Bills:           bills,
		CustomNextItems: customItems,
		Tenant:          tenant,
		LandlordName:    landlord.Name,
		LandlordPhone:   landlord.Phone,
		Applications:    []any{},
	}, nil
}

func (uc *UseCase) GetProperty(ctx context.Context, id, landlordID string) (entity.Property, error) {
	property, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return entity.Property{}, fmt.Errorf("PropertyUseCase - GetProperty - uc.repo.GetByID: %w", err)
	}

	if property.LandlordID != landlordID {
		return entity.Property{}, repo.ErrPropertyNotFound
	}

	return property, nil
}

func (uc *UseCase) UpdateProperty(ctx context.Context, id, landlordID string, body request.Property) error {
	existing, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("PropertyUseCase - UpdateProperty - uc.repo.GetByID: %w", err)
	}

	if existing.LandlordID != landlordID {
		return repo.ErrPropertyNotFound
	}

	property := entity.Property{
		ID:          id,
		LandlordID:  existing.LandlordID,
		Name:        body.Name,
		Coordinates: body.Coordinates,
		Country:     body.Country,
		Region:      body.Region,
		City:        body.City,
		Street:      body.Street,
		House:       body.House,
		Apartment:   body.Apartment,
		GvsTariff:   body.GvsTariff,
		HvsTariff:   body.HvsTariff,
		El1Tariff:   body.El1Tariff,
		El2Tariff:   body.El2Tariff,
		Balance:     existing.Balance,
	}

	if err := uc.repo.Update(ctx, property); err != nil {
		return fmt.Errorf("PropertyUseCase - UpdateProperty - uc.repo.Update: %w", err)
	}

	return nil
}

func (uc *UseCase) DeleteProperty(ctx context.Context, id, landlordID string) error {
	existing, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("PropertyUseCase - DeleteProperty - uc.repo.GetByID: %w", err)
	}

	if existing.LandlordID != landlordID {
		return repo.ErrPropertyNotFound
	}

	if err := uc.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("PropertyUseCase - DeleteProperty - uc.repo.Delete: %w", err)
	}

	return nil
}

// CreateLease moves a registered tenant into the property, superseding whatever lease
// (if any) previously occupied it, since a property can only have one active tenant.
func (uc *UseCase) CreateLease(ctx context.Context, propertyID, landlordID string, body request.Lease) error {
	prop, err := uc.repo.GetByID(ctx, propertyID)
	if err != nil {
		return fmt.Errorf("PropertyUseCase - CreateLease - uc.repo.GetByID: %w", err)
	}

	if prop.LandlordID != landlordID {
		return repo.ErrPropertyNotFound
	}

	tenant, err := uc.users.GetByID(ctx, body.TenantUserID)
	if err != nil {
		if errors.Is(err, repo.ErrUserNotFound) {
			return repo.ErrTenantNotFound
		}
		return fmt.Errorf("PropertyUseCase - CreateLease - uc.users.GetByID: %w", err)
	}

	if tenant.Role != entity.RoleTenant {
		return repo.ErrTenantNotFound
	}

	startDate := time.Now().UTC()
	endDate := startDate.AddDate(0, body.MonthsOfRent, 0)

	lease := entity.Lease{
		ID:           uuid.NewString(),
		PropertyID:   propertyID,
		TenantUserID: tenant.ID,
		Name:         tenant.Name,
		Document:     tenant.Document,
		Phone:        tenant.Phone,
		MonthsOfRent: body.MonthsOfRent,
		Price:        body.Price,
		PaymentDay:   body.PaymentDay,
		ReadingDay:   body.ReadingDay,
		StartDate:    startDate.Format(time.RFC3339),
		EndDate:      endDate.Format(time.RFC3339),
	}

	if err := uc.leases.Upsert(ctx, lease); err != nil {
		return fmt.Errorf("PropertyUseCase - CreateLease - uc.leases.Upsert: %w", err)
	}

	return nil
}

func (uc *UseCase) DeleteLease(ctx context.Context, propertyID, landlordID string) error {
	prop, err := uc.repo.GetByID(ctx, propertyID)
	if err != nil {
		return fmt.Errorf("PropertyUseCase - DeleteLease - uc.repo.GetByID: %w", err)
	}

	if prop.LandlordID != landlordID {
		return repo.ErrPropertyNotFound
	}

	if err := uc.leases.DeleteByPropertyID(ctx, propertyID); err != nil {
		return fmt.Errorf("PropertyUseCase - DeleteLease - uc.leases.DeleteByPropertyID: %w", err)
	}

	return nil
}
