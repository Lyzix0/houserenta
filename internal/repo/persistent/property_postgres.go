package persistent

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	entity "github.com/potom_pridumaem/internal/entity/users"
	"github.com/potom_pridumaem/internal/repo"
	"github.com/potom_pridumaem/pkg/postgres"
)

var propertyColumns = []string{
	"id", "landlord_id", "name", "coordinates", "country",
	"region", "city", "street", "house", "apartment",
	"gvs_tariff", "hvs_tariff", "el1_tariff", "el2_tariff", "balance",
}

type PropertyRepo struct {
	*postgres.Postgres
}

func NewPropertyRepo(pg *postgres.Postgres) *PropertyRepo {
	return &PropertyRepo{pg}
}

func (r *PropertyRepo) Store(ctx context.Context, prop entity.Property) error {
	sql, args, err := r.Builder.
		Insert("app.properties").
		Columns(propertyColumns...).
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

func (r *PropertyRepo) GetByLandlordID(ctx context.Context, landlordID string) ([]entity.Property, error) {
	sql, args, err := r.Builder.
		Select(propertyColumns...).
		From("app.properties").
		Where("landlord_id = ?", landlordID).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("PropertyRepo - GetByLandlordID - r.Builder: %w", err)
	}

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("PropertyRepo - GetByLandlordID - r.Pool.Query: %w", err)
	}
	defer rows.Close()

	properties := make([]entity.Property, 0)
	for rows.Next() {
		var prop entity.Property
		if err := rows.Scan(
			&prop.ID, &prop.LandlordID, &prop.Name, &prop.Coordinates, &prop.Country,
			&prop.Region, &prop.City, &prop.Street, &prop.House, &prop.Apartment,
			&prop.GvsTariff, &prop.HvsTariff, &prop.El1Tariff, &prop.El2Tariff, &prop.Balance,
		); err != nil {
			return nil, fmt.Errorf("PropertyRepo - GetByLandlordID - rows.Scan: %w", err)
		}
		properties = append(properties, prop)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("PropertyRepo - GetByLandlordID - rows.Err: %w", err)
	}

	return properties, nil
}

func (r *PropertyRepo) GetByID(ctx context.Context, id string) (entity.Property, error) {
	sql, args, err := r.Builder.
		Select(propertyColumns...).
		From("app.properties").
		Where("id = ?", id).
		ToSql()
	if err != nil {
		return entity.Property{}, fmt.Errorf("PropertyRepo - GetByID - r.Builder: %w", err)
	}

	var prop entity.Property
	err = r.Pool.QueryRow(ctx, sql, args...).Scan(
		&prop.ID, &prop.LandlordID, &prop.Name, &prop.Coordinates, &prop.Country,
		&prop.Region, &prop.City, &prop.Street, &prop.House, &prop.Apartment,
		&prop.GvsTariff, &prop.HvsTariff, &prop.El1Tariff, &prop.El2Tariff, &prop.Balance,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Property{}, repo.ErrPropertyNotFound
		}
		return entity.Property{}, fmt.Errorf("PropertyRepo - GetByID - r.Pool.QueryRow: %w", err)
	}

	return prop, nil
}

func (r *PropertyRepo) Update(ctx context.Context, prop entity.Property) error {
	sql, args, err := r.Builder.
		Update("app.properties").
		Set("name", prop.Name).
		Set("coordinates", prop.Coordinates).
		Set("country", prop.Country).
		Set("region", prop.Region).
		Set("city", prop.City).
		Set("street", prop.Street).
		Set("house", prop.House).
		Set("apartment", prop.Apartment).
		Set("gvs_tariff", prop.GvsTariff).
		Set("hvs_tariff", prop.HvsTariff).
		Set("el1_tariff", prop.El1Tariff).
		Set("el2_tariff", prop.El2Tariff).
		Where("id = ?", prop.ID).
		ToSql()
	if err != nil {
		return fmt.Errorf("PropertyRepo - Update - r.Builder: %w", err)
	}

	tag, err := r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("PropertyRepo - Update - r.Pool.Exec: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return repo.ErrPropertyNotFound
	}

	return nil
}

func (r *PropertyRepo) Delete(ctx context.Context, id string) error {
	sql, args, err := r.Builder.
		Delete("app.properties").
		Where("id = ?", id).
		ToSql()
	if err != nil {
		return fmt.Errorf("PropertyRepo - Delete - r.Builder: %w", err)
	}

	tag, err := r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("PropertyRepo - Delete - r.Pool.Exec: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return repo.ErrPropertyNotFound
	}

	return nil
}
