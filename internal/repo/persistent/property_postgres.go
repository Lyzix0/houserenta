package persistent

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	entity "github.com/potom_pridumaem/internal/entity/users"
	"github.com/potom_pridumaem/internal/repo"
	"github.com/potom_pridumaem/pkg/postgres"
)

type PropertyRepo struct {
	*postgres.Postgres
}

func NewPropertyRepo(pg *postgres.Postgres) *PropertyRepo {
	return &PropertyRepo{pg}
}

func (r *PropertyRepo) Store(ctx context.Context, prop entity.Property) error {
	sql, args, err := r.Builder.
		Insert("app.properties").
		Columns(
			"id", "landlord_id", "name", "coordinates", "country",
			"region", "city", "street", "house", "apartment",
			"gvs_tariff", "hvs_tariff", "el1_tariff", "el2_tariff", "balance",
		).
		Values(
			prop.ID, prop.LandlordID, prop.Name, prop.Coordinates, prop.Country,
			prop.Region, prop.City, prop.Street, prop.House, prop.Apartment,
			prop.GvsTariff, prop.HvsTariff, prop.El1Tariff, prop.El2Tariff, prop.Balance,
		).
		ToSql()
	if err != nil {
		return fmt.Errorf("PropertyRepo - Store - r.Builder: %w", err)
	}

	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505":
				return repo.ErrPropertyAlreadyExists
			case "23503":
				return repo.ErrLandlordNotFound
			}
		}
		return fmt.Errorf("PropertyRepo - Store - r.Pool.Exec: %w", err)
	}

	return nil
}
