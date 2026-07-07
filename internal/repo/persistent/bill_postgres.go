package persistent

import (
	"context"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	entity "github.com/potom_pridumaem/internal/entity/users"
	"github.com/potom_pridumaem/internal/repo"
	"github.com/potom_pridumaem/pkg/postgres"
)

var billColumns = []string{"id", "property_id", "date", "due_date", "status", "type", "total"}

type BillRepo struct {
	*postgres.Postgres
}

func NewBillRepo(pg *postgres.Postgres) *BillRepo {
	return &BillRepo{pg}
}

func (r *BillRepo) GetByPropertyID(ctx context.Context, propertyID string) ([]entity.Bill, error) {
	sql, args, err := r.Builder.
		Select(billColumns...).
		From("app.bills").
		Where("property_id = ?", propertyID).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("BillRepo - GetByPropertyID - r.Builder: %w", err)
	}

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("BillRepo - GetByPropertyID - r.Pool.Query: %w", err)
	}
	defer rows.Close()

	bills := make([]entity.Bill, 0)
	billIndex := make(map[string]int)
	billIDs := make([]string, 0)
	for rows.Next() {
		var bill entity.Bill
		if err := rows.Scan(&bill.ID, &bill.PropertyID, &bill.Date, &bill.DueDate, &bill.Status, &bill.Type, &bill.Total); err != nil {
			return nil, fmt.Errorf("BillRepo - GetByPropertyID - rows.Scan: %w", err)
		}
		bill.Items = make([]entity.BillItem, 0)
		billIndex[bill.ID] = len(bills)
		billIDs = append(billIDs, bill.ID)
		bills = append(bills, bill)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("BillRepo - GetByPropertyID - rows.Err: %w", err)
	}

	if len(billIDs) == 0 {
		return bills, nil
	}

	itemsSQL, itemsArgs, err := r.Builder.
		Select("id", "bill_id", "description", "amount").
		From("app.bill_items").
		Where(squirrel.Eq{"bill_id": billIDs}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("BillRepo - GetByPropertyID - r.Builder items: %w", err)
	}

	itemRows, err := r.Pool.Query(ctx, itemsSQL, itemsArgs...)
	if err != nil {
		return nil, fmt.Errorf("BillRepo - GetByPropertyID - r.Pool.Query items: %w", err)
	}
	defer itemRows.Close()

	for itemRows.Next() {
		var item entity.BillItem
		if err := itemRows.Scan(&item.ID, &item.BillID, &item.Description, &item.Amount); err != nil {
			return nil, fmt.Errorf("BillRepo - GetByPropertyID - itemRows.Scan: %w", err)
		}
		if idx, ok := billIndex[item.BillID]; ok {
			bills[idx].Items = append(bills[idx].Items, item)
		}
	}

	if err := itemRows.Err(); err != nil {
		return nil, fmt.Errorf("BillRepo - GetByPropertyID - itemRows.Err: %w", err)
	}

	return bills, nil
}

func (r *BillRepo) UpdateStatus(ctx context.Context, billID, propertyID, status string) error {
	sql, args, err := r.Builder.
		Update("app.bills").
		Set("status", status).
		Where(squirrel.Eq{"id": billID, "property_id": propertyID}).
		ToSql()
	if err != nil {
		return fmt.Errorf("BillRepo - UpdateStatus - r.Builder: %w", err)
	}

	tag, err := r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("BillRepo - UpdateStatus - r.Pool.Exec: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return repo.ErrBillNotFound
	}

	return nil
}

// GetLastByPropertyIDAndType returns the most recently dated bill of the given type
// for the property, or repo.ErrBillNotFound if none has ever been issued.
func (r *BillRepo) GetLastByPropertyIDAndType(ctx context.Context, propertyID, billType string) (entity.Bill, error) {
	sql, args, err := r.Builder.
		Select(billColumns...).
		From("app.bills").
		Where(squirrel.Eq{"property_id": propertyID, "type": billType}).
		OrderBy("date DESC").
		Limit(1).
		ToSql()
	if err != nil {
		return entity.Bill{}, fmt.Errorf("BillRepo - GetLastByPropertyIDAndType - r.Builder: %w", err)
	}

	var bill entity.Bill
	err = r.Pool.QueryRow(ctx, sql, args...).Scan(
		&bill.ID, &bill.PropertyID, &bill.Date, &bill.DueDate, &bill.Status, &bill.Type, &bill.Total,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Bill{}, repo.ErrBillNotFound
		}
		return entity.Bill{}, fmt.Errorf("BillRepo - GetLastByPropertyIDAndType - r.Pool.QueryRow: %w", err)
	}

	return bill, nil
}

// Store inserts the bill and all of its line items atomically.
func (r *BillRepo) Store(ctx context.Context, bill entity.Bill) error {
	tx, err := r.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("BillRepo - Store - r.Pool.Begin: %w", err)
	}
	defer tx.Rollback(ctx) //nolint:errcheck

	billSQL, billArgs, err := r.Builder.
		Insert("app.bills").
		Columns(billColumns...).
		Values(bill.ID, bill.PropertyID, bill.Date, bill.DueDate, bill.Status, bill.Type, bill.Total).
		ToSql()
	if err != nil {
		return fmt.Errorf("BillRepo - Store - r.Builder bill: %w", err)
	}

	if _, err := tx.Exec(ctx, billSQL, billArgs...); err != nil {
		return fmt.Errorf("BillRepo - Store - tx.Exec bill: %w", err)
	}

	for _, item := range bill.Items {
		itemSQL, itemArgs, err := r.Builder.
			Insert("app.bill_items").
			Columns("id", "bill_id", "description", "amount").
			Values(item.ID, bill.ID, item.Description, item.Amount).
			ToSql()
		if err != nil {
			return fmt.Errorf("BillRepo - Store - r.Builder item: %w", err)
		}

		if _, err := tx.Exec(ctx, itemSQL, itemArgs...); err != nil {
			return fmt.Errorf("BillRepo - Store - tx.Exec item: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("BillRepo - Store - tx.Commit: %w", err)
	}

	return nil
}
