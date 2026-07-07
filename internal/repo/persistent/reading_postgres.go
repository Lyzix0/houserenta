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

var readingColumns = []string{"id", "property_id", "date", "gvs", "hvs", "el1", "el2", "is_accounted"}

type ReadingRepo struct {
	*postgres.Postgres
}

func NewReadingRepo(pg *postgres.Postgres) *ReadingRepo {
	return &ReadingRepo{pg}
}

func (r *ReadingRepo) GetByPropertyID(ctx context.Context, propertyID string) ([]entity.Reading, error) {
	sql, args, err := r.Builder.
		Select(readingColumns...).
		From("app.readings").
		Where("property_id = ?", propertyID).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("ReadingRepo - GetByPropertyID - r.Builder: %w", err)
	}

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("ReadingRepo - GetByPropertyID - r.Pool.Query: %w", err)
	}
	defer rows.Close()

	readings := make([]entity.Reading, 0)
	for rows.Next() {
		var reading entity.Reading
		if err := rows.Scan(
			&reading.ID, &reading.PropertyID, &reading.Date, &reading.Gvs, &reading.Hvs,
			&reading.El1, &reading.El2, &reading.IsAccounted,
		); err != nil {
			return nil, fmt.Errorf("ReadingRepo - GetByPropertyID - rows.Scan: %w", err)
		}
		readings = append(readings, reading)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ReadingRepo - GetByPropertyID - rows.Err: %w", err)
	}

	return readings, nil
}

func (r *ReadingRepo) Store(ctx context.Context, reading entity.Reading) error {
	sql, args, err := r.Builder.
		Insert("app.readings").
		Columns(readingColumns...).
		Values(
			reading.ID, reading.PropertyID, reading.Date, reading.Gvs, reading.Hvs,
			reading.El1, reading.El2, reading.IsAccounted,
		).
		ToSql()
	if err != nil {
		return fmt.Errorf("ReadingRepo - Store - r.Builder: %w", err)
	}

	if _, err := r.Pool.Exec(ctx, sql, args...); err != nil {
		return fmt.Errorf("ReadingRepo - Store - r.Pool.Exec: %w", err)
	}

	return nil
}

func (r *ReadingRepo) scanOne(ctx context.Context, propertyID string, isAccounted int, order string) (entity.Reading, error) {
	sql, args, err := r.Builder.
		Select(readingColumns...).
		From("app.readings").
		Where("property_id = ? AND is_accounted = ?", propertyID, isAccounted).
		OrderBy("date " + order).
		Limit(1).
		ToSql()
	if err != nil {
		return entity.Reading{}, fmt.Errorf("ReadingRepo - scanOne - r.Builder: %w", err)
	}

	var reading entity.Reading
	err = r.Pool.QueryRow(ctx, sql, args...).Scan(
		&reading.ID, &reading.PropertyID, &reading.Date, &reading.Gvs, &reading.Hvs,
		&reading.El1, &reading.El2, &reading.IsAccounted,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Reading{}, repo.ErrReadingNotFound
		}
		return entity.Reading{}, fmt.Errorf("ReadingRepo - scanOne - r.Pool.QueryRow: %w", err)
	}

	return reading, nil
}

// GetOldestUnaccounted returns the oldest reading not yet folded into a bill (is_accounted = 0).
func (r *ReadingRepo) GetOldestUnaccounted(ctx context.Context, propertyID string) (entity.Reading, error) {
	return r.scanOne(ctx, propertyID, 0, "ASC")
}

// GetLatestAccounted returns the most recent reading already billed (is_accounted = 1),
// used as the baseline to compute consumption since then.
func (r *ReadingRepo) GetLatestAccounted(ctx context.Context, propertyID string) (entity.Reading, error) {
	return r.scanOne(ctx, propertyID, 1, "DESC")
}

func (r *ReadingRepo) MarkAccounted(ctx context.Context, id string) error {
	sql, args, err := r.Builder.
		Update("app.readings").
		Set("is_accounted", 1).
		Where("id = ?", id).
		ToSql()
	if err != nil {
		return fmt.Errorf("ReadingRepo - MarkAccounted - r.Builder: %w", err)
	}

	if _, err := r.Pool.Exec(ctx, sql, args...); err != nil {
		return fmt.Errorf("ReadingRepo - MarkAccounted - r.Pool.Exec: %w", err)
	}

	return nil
}
