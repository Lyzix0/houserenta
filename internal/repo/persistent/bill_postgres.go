package persistent

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	entity "github.com/potom_pridumaem/internal/entity/users"
	"github.com/potom_pridumaem/pkg/postgres"
)

type BillRepo struct {
	*postgres.Postgres
}

func NewBillRepo(pg *postgres.Postgres) *BillRepo {
	return &BillRepo{pg}
}

func (r *BillRepo) GetByPropertyID(ctx context.Context, propertyID string) ([]entity.Bill, error) {
	sql, args, err := r.Builder.
		Select("id", "property_id", "date", "due_date", "status", "total").
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
		if err := rows.Scan(&bill.ID, &bill.PropertyID, &bill.Date, &bill.DueDate, &bill.Status, &bill.Total); err != nil {
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
