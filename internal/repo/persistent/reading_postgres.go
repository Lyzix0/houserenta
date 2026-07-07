package persistent

import (
	"context"
	"fmt"

	entity "github.com/potom_pridumaem/internal/entity/users"
	"github.com/potom_pridumaem/pkg/postgres"
)

type ReadingRepo struct {
	*postgres.Postgres
}

func NewReadingRepo(pg *postgres.Postgres) *ReadingRepo {
	return &ReadingRepo{pg}
}

func (r *ReadingRepo) GetByPropertyID(ctx context.Context, propertyID string) ([]entity.Reading, error) {
	sql, args, err := r.Builder.
		Select("id", "property_id", "date", "gvs", "hvs", "el1", "el2", "is_accounted").
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
