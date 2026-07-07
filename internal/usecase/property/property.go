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
		LandlordID:  req.LandlordID,
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

func (uc *UseCase) GetProperty(ctx context.Context, id string) (entity.Property, error) {
	property, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return entity.Property{}, fmt.Errorf("PropertyUseCase - GetProperty - uc.repo.GetByID: %w", err)
	}
	return property, nil
}

func (uc *UseCase) UpdateProperty(ctx context.Context, id string, req request.Property) error {
	property := entity.Property{
		ID:          id,
		LandlordID:  req.LandlordID,
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
	}

	if err := uc.repo.Update(ctx, &property); err != nil {
		return fmt.Errorf("PropertyUseCase - UpdateProperty - uc.repo.Update: %w", err)
	}

	return nil
}

func (uc *UseCase) DeleteProperty(ctx context.Context, id string) error {
	if err := uc.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("PropertyUseCase - DeleteProperty - uc.repo.Delete: %w", err)
	}
	return nil
}

func (uc *UseCase) GetVacantProperties(ctx context.Context) ([]entity.Property, error) {
	properties, err := uc.repo.GetVacant(ctx)
	if err != nil {
		return nil, fmt.Errorf("PropertyUseCase - GetVacantProperties - uc.repo.GetVacant: %w", err)
	}
	return properties, nil
}
