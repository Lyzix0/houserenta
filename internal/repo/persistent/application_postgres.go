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

type ApplicationRepo struct {
	*postgres.Postgres
}

func NewApplicationRepo(pg *postgres.Postgres) *ApplicationRepo {
	return &ApplicationRepo{pg}
}

func (r *ApplicationRepo) Store(ctx context.Context, app entity.Application) error {
	sql, args, err := r.Builder.
		Insert("app.applications").
		Columns("id", "property_id", "tenant_user_id", "date").
		Values(app.ID, app.PropertyID, app.TenantUserID, app.Date).
		ToSql()
	if err != nil {
		return fmt.Errorf("ApplicationRepo - Store - r.Builder: %w", err)
	}

	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505":
				return repo.ErrApplicationAlreadyExists
			case "23503":
				return repo.ErrPropertyNotFound
			}
		}
		return fmt.Errorf("ApplicationRepo - Store - r.Pool.Exec: %w", err)
	}

	return nil
}

func (r *ApplicationRepo) GetByPropertyID(ctx context.Context, propertyID string) ([]entity.Application, error) {
	sql, args, err := r.Builder.
		Select("id", "property_id", "tenant_user_id", "date").
		From("app.applications").
		Where("property_id = ?", propertyID).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("ApplicationRepo - GetByPropertyID - r.Builder: %w", err)
	}

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("ApplicationRepo - GetByPropertyID - r.Pool.Query: %w", err)
	}
	defer rows.Close()

	applications := make([]entity.Application, 0)
	for rows.Next() {
		var app entity.Application
		if err := rows.Scan(&app.ID, &app.PropertyID, &app.TenantUserID, &app.Date); err != nil {
			return nil, fmt.Errorf("ApplicationRepo - GetByPropertyID - rows.Scan: %w", err)
		}
		applications = append(applications, app)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ApplicationRepo - GetByPropertyID - rows.Err: %w", err)
	}

	return applications, nil
}

func (r *ApplicationRepo) DeleteByPropertyID(ctx context.Context, propertyID string) error {
	sql, args, err := r.Builder.
		Delete("app.applications").
		Where("property_id = ?", propertyID).
		ToSql()
	if err != nil {
		return fmt.Errorf("ApplicationRepo - DeleteByPropertyID - r.Builder: %w", err)
	}

	if _, err := r.Pool.Exec(ctx, sql, args...); err != nil {
		return fmt.Errorf("ApplicationRepo - DeleteByPropertyID - r.Pool.Exec: %w", err)
	}

	return nil
}
