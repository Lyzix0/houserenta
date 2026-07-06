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
