package persistent

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	entity "github.com/potom_pridumaem/internal/entity/users"
	"github.com/potom_pridumaem/internal/repo"
	"github.com/potom_pridumaem/pkg/postgres"
)

type LeaseRepo struct {
	*postgres.Postgres
}

func NewLeaseRepo(pg *postgres.Postgres) *LeaseRepo {
	return &LeaseRepo{pg}
}

var leaseColumns = []string{
	"id", "property_id", "tenant_user_id", "name", "document", "phone",
	"months_of_rent", "price", "payment_day", "reading_day", "start_date", "end_date",
}

func (r *LeaseRepo) GetByTenantUserID(
	ctx context.Context,
	tenantUserID string,
) (entity.Lease, error) {
	sql, args, err := r.Builder.
		Select(leaseColumns...).
		From("app.leases").
		Where("tenant_user_id = ?", tenantUserID).
		ToSql()
	if err != nil {
		return entity.Lease{}, fmt.Errorf("LeaseRepo - GetByTenantUserID - r.Builder: %w", err)
	}

	var lease entity.Lease
	err = r.Pool.QueryRow(ctx, sql, args...).Scan(
		&lease.ID, &lease.PropertyID, &lease.TenantUserID, &lease.Name, &lease.Document, &lease.Phone,
		&lease.MonthsOfRent, &lease.Price, &lease.PaymentDay, &lease.ReadingDay, &lease.StartDate, &lease.EndDate,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Lease{}, repo.ErrLeaseNotFound
		}
		return entity.Lease{}, fmt.Errorf("LeaseRepo - GetByTenantUserID - r.Pool.QueryRow: %w", err)
	}

	return lease, nil
}

func (r *LeaseRepo) GetByPropertyID(
	ctx context.Context,
	propertyID string,
) (entity.Lease, error) {
	sql, args, err := r.Builder.
		Select(leaseColumns...).
		From("app.leases").
		Where("property_id = ?", propertyID).
		ToSql()
	if err != nil {
		return entity.Lease{}, fmt.Errorf("LeaseRepo - GetByPropertyID - r.Builder: %w", err)
	}

	var lease entity.Lease
	err = r.Pool.QueryRow(ctx, sql, args...).Scan(
		&lease.ID, &lease.PropertyID, &lease.TenantUserID, &lease.Name, &lease.Document, &lease.Phone,
		&lease.MonthsOfRent, &lease.Price, &lease.PaymentDay, &lease.ReadingDay, &lease.StartDate, &lease.EndDate,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Lease{}, repo.ErrLeaseNotFound
		}
		return entity.Lease{}, fmt.Errorf("LeaseRepo - GetByPropertyID - r.Pool.QueryRow: %w", err)
	}

	return lease, nil
}

// Upsert creates a lease for the property, or replaces the existing one (keeping its ID)
// if the property already has an active lease, since app.leases.property_id is unique.
func (r *LeaseRepo) Upsert(ctx context.Context, lease entity.Lease) error {
	sql, args, err := r.Builder.
		Insert("app.leases").
		Columns(leaseColumns...).
		Values(
			lease.ID, lease.PropertyID, lease.TenantUserID, lease.Name, lease.Document, lease.Phone,
			lease.MonthsOfRent, lease.Price, lease.PaymentDay, lease.ReadingDay, lease.StartDate, lease.EndDate,
		).
		Suffix(`ON CONFLICT (property_id) DO UPDATE SET
			tenant_user_id = EXCLUDED.tenant_user_id,
			name = EXCLUDED.name,
			document = EXCLUDED.document,
			phone = EXCLUDED.phone,
			months_of_rent = EXCLUDED.months_of_rent,
			price = EXCLUDED.price,
			payment_day = EXCLUDED.payment_day,
			reading_day = EXCLUDED.reading_day,
			start_date = EXCLUDED.start_date,
			end_date = EXCLUDED.end_date`).
		ToSql()
	if err != nil {
		return fmt.Errorf("LeaseRepo - Upsert - r.Builder: %w", err)
	}

	if _, err := r.Pool.Exec(ctx, sql, args...); err != nil {
		return fmt.Errorf("LeaseRepo - Upsert - r.Pool.Exec: %w", err)
	}

	return nil
}

func (r *LeaseRepo) DeleteByPropertyID(ctx context.Context, propertyID string) error {
	sql, args, err := r.Builder.
		Delete("app.leases").
		Where("property_id = ?", propertyID).
		ToSql()
	if err != nil {
		return fmt.Errorf("LeaseRepo - DeleteByPropertyID - r.Builder: %w", err)
	}

	if _, err := r.Pool.Exec(ctx, sql, args...); err != nil {
		return fmt.Errorf("LeaseRepo - DeleteByPropertyID - r.Pool.Exec: %w", err)
	}

	return nil
}
