package persistent

import (
	"context"
	"fmt"

	entity "github.com/potom_pridumaem/internal/entity/users"
	"github.com/potom_pridumaem/pkg/postgres"
)

type CustomNextItemRepo struct {
	*postgres.Postgres
}

func NewCustomNextItemRepo(pg *postgres.Postgres) *CustomNextItemRepo {
	return &CustomNextItemRepo{pg}
}

func (r *CustomNextItemRepo) GetByPropertyID(ctx context.Context, propertyID string) ([]entity.CustomNextItem, error) {
	sql, args, err := r.Builder.
		Select("id", "property_id", "description", "amount").
		From("app.custom_next_items").
		Where("property_id = ?", propertyID).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("CustomNextItemRepo - GetByPropertyID - r.Builder: %w", err)
	}

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("CustomNextItemRepo - GetByPropertyID - r.Pool.Query: %w", err)
	}
	defer rows.Close()

	items := make([]entity.CustomNextItem, 0)
	for rows.Next() {
		var item entity.CustomNextItem
		if err := rows.Scan(&item.ID, &item.PropertyID, &item.Description, &item.Amount); err != nil {
			return nil, fmt.Errorf("CustomNextItemRepo - GetByPropertyID - rows.Scan: %w", err)
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("CustomNextItemRepo - GetByPropertyID - rows.Err: %w", err)
	}

	return items, nil
}

func (r *CustomNextItemRepo) Store(ctx context.Context, item entity.CustomNextItem) error {
	sql, args, err := r.Builder.
		Insert("app.custom_next_items").
		Columns("id", "property_id", "description", "amount").
		Values(item.ID, item.PropertyID, item.Description, item.Amount).
		ToSql()
	if err != nil {
		return fmt.Errorf("CustomNextItemRepo - Store - r.Builder: %w", err)
	}

	if _, err := r.Pool.Exec(ctx, sql, args...); err != nil {
		return fmt.Errorf("CustomNextItemRepo - Store - r.Pool.Exec: %w", err)
	}

	return nil
}

func (r *CustomNextItemRepo) DeleteByPropertyID(ctx context.Context, propertyID string) error {
	sql, args, err := r.Builder.
		Delete("app.custom_next_items").
		Where("property_id = ?", propertyID).
		ToSql()
	if err != nil {
		return fmt.Errorf("CustomNextItemRepo - DeleteByPropertyID - r.Builder: %w", err)
	}

	if _, err := r.Pool.Exec(ctx, sql, args...); err != nil {
		return fmt.Errorf("CustomNextItemRepo - DeleteByPropertyID - r.Pool.Exec: %w", err)
	}

	return nil
}
