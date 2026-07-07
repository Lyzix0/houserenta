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
