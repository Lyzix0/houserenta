package property

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/potom_pridumaem/internal/controller/v1/request"
	entity "github.com/potom_pridumaem/internal/entity/users"
	"github.com/potom_pridumaem/internal/repo"
)

type UseCase struct {
	repo repo.PropertyRepo
}

func New(r repo.PropertyRepo) *UseCase {
	return &UseCase{
		repo: r,
	}
}

func (uc *UseCase) CreateProperty(ctx context.Context, req request.Property) (entity.Property, error) {
	property := entity.Property{
		ID:          uuid.NewString(),
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

func (uc *UseCase) GetProperties(ctx context.Context, landlordID string) ([]entity.Property, error) {
	properties, err := uc.repo.GetByLandlordID(ctx, landlordID)
	if err != nil {
		return nil, fmt.Errorf("PropertyUseCase - GetProperties - uc.repo.GetByLandlordID: %w", err)
	}

	return properties, nil
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
