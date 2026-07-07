package billing

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	entity "github.com/potom_pridumaem/internal/entity/users"
	"github.com/potom_pridumaem/internal/repo"
)

const (
	billingInterval = 30 * 24 * time.Hour
	billDueInDays   = 10
	billTypeRent    = "rent"
)

type UseCase struct {
	leases      repo.LeaseRepo
	properties  repo.PropertyRepo
	bills       repo.BillRepo
	readings    repo.ReadingRepo
	customItems repo.CustomNextItemRepo
}

func New(
	leases repo.LeaseRepo,
	properties repo.PropertyRepo,
	bills repo.BillRepo,
	readings repo.ReadingRepo,
	customItems repo.CustomNextItemRepo,
) *UseCase {
	return &UseCase{
		leases:      leases,
		properties:  properties,
		bills:       bills,
		readings:    readings,
		customItems: customItems,
	}
}

// Run is the lazy billing check: for every active lease, it issues a new rent bill
// if 30+ days have passed since the last one (or none has been issued since move-in).
func (uc *UseCase) Run(ctx context.Context) error {
	leases, err := uc.leases.GetAll(ctx)
	if err != nil {
		return fmt.Errorf("BillingUseCase - Run - uc.leases.GetAll: %w", err)
	}

	for _, lease := range leases {
		due, err := uc.isDue(ctx, lease.PropertyID)
		if err != nil {
			return fmt.Errorf("BillingUseCase - Run - uc.isDue: %w", err)
		}

		if !due {
			continue
		}

		if err := uc.generateBill(ctx, lease); err != nil {
			return fmt.Errorf("BillingUseCase - Run - uc.generateBill: %w", err)
		}
	}

	return nil
}

func (uc *UseCase) isDue(ctx context.Context, propertyID string) (bool, error) {
	lastBill, err := uc.bills.GetLastByPropertyIDAndType(ctx, propertyID, billTypeRent)
	if errors.Is(err, repo.ErrBillNotFound) {
		return true, nil
	}
	if err != nil {
		return false, fmt.Errorf("uc.bills.GetLastByPropertyIDAndType: %w", err)
	}

	lastDate, err := time.Parse(time.RFC3339, lastBill.Date)
	if err != nil {
		return false, fmt.Errorf("time.Parse last bill date: %w", err)
	}

	return time.Since(lastDate) >= billingInterval, nil
}

func (uc *UseCase) generateBill(ctx context.Context, lease entity.Lease) error {
	property, err := uc.properties.GetByID(ctx, lease.PropertyID)
	if err != nil {
		return fmt.Errorf("uc.properties.GetByID: %w", err)
	}

	now := time.Now().UTC()
	bill := entity.Bill{
		ID:         uuid.NewString(),
		PropertyID: property.ID,
		Date:       now.Format(time.RFC3339),
		DueDate:    now.AddDate(0, 0, billDueInDays).Format(time.RFC3339),
		Status:     "unpaid",
		Type:       billTypeRent,
	}

	total := lease.Price
	bill.Items = append(bill.Items, entity.BillItem{
		ID:          uuid.NewString(),
		BillID:      bill.ID,
		Description: "Аренда за период",
		Amount:      lease.Price,
	})

	customItems, err := uc.customItems.GetByPropertyID(ctx, property.ID)
	if err != nil {
		return fmt.Errorf("uc.customItems.GetByPropertyID: %w", err)
	}

	for _, item := range customItems {
		bill.Items = append(bill.Items, entity.BillItem{
			ID:          uuid.NewString(),
			BillID:      bill.ID,
			Description: item.Description,
			Amount:      item.Amount,
		})
		total += item.Amount
	}

	utilityItems, utilityTotal, accountedReadingID, err := uc.calculateUtilities(ctx, property, bill.ID)
	if err != nil {
		return fmt.Errorf("uc.calculateUtilities: %w", err)
	}
	bill.Items = append(bill.Items, utilityItems...)
	total += utilityTotal

	bill.Total = total

	if err := uc.bills.Store(ctx, bill); err != nil {
		return fmt.Errorf("uc.bills.Store: %w", err)
	}

	if accountedReadingID != "" {
		if err := uc.readings.MarkAccounted(ctx, accountedReadingID); err != nil {
			return fmt.Errorf("uc.readings.MarkAccounted: %w", err)
		}
	}

	if len(customItems) > 0 {
		if err := uc.customItems.DeleteByPropertyID(ctx, property.ID); err != nil {
			return fmt.Errorf("uc.customItems.DeleteByPropertyID: %w", err)
		}
	}

	if err := uc.properties.AddBalance(ctx, property.ID, -total); err != nil {
		return fmt.Errorf("uc.properties.AddBalance: %w", err)
	}

	return nil
}

// calculateUtilities prices the consumption between the last billed reading and the
// oldest not-yet-billed one. It returns the resulting bill items, their combined total,
// and the ID of the reading to mark accounted (empty if there was nothing to bill).
func (uc *UseCase) calculateUtilities(ctx context.Context, property entity.Property, billID string) ([]entity.BillItem, float64, string, error) {
	newest, err := uc.readings.GetOldestUnaccounted(ctx, property.ID)
	if errors.Is(err, repo.ErrReadingNotFound) {
		return nil, 0, "", nil
	}
	if err != nil {
		return nil, 0, "", fmt.Errorf("uc.readings.GetOldestUnaccounted: %w", err)
	}

	baseline, err := uc.readings.GetLatestAccounted(ctx, property.ID)
	switch {
	case err == nil:
		// baseline already populated below
	case errors.Is(err, repo.ErrReadingNotFound):
		baseline = entity.Reading{}
	default:
		return nil, 0, "", fmt.Errorf("uc.readings.GetLatestAccounted: %w", err)
	}

	var items []entity.BillItem
	var total float64

	addItem := func(diff, tariff float64, description string) {
		if diff <= 0 {
			return
		}
		amount := diff * tariff
		items = append(items, entity.BillItem{
			ID:          uuid.NewString(),
			BillID:      billID,
			Description: description,
			Amount:      amount,
		})
		total += amount
	}

	addItem(newest.Gvs-baseline.Gvs, property.GvsTariff, fmt.Sprintf("ГВС (расход: %.2f куб.м)", newest.Gvs-baseline.Gvs))
	addItem(newest.Hvs-baseline.Hvs, property.HvsTariff, fmt.Sprintf("ХВС (расход: %.2f куб.м)", newest.Hvs-baseline.Hvs))
	addItem(newest.El1-baseline.El1, property.El1Tariff, fmt.Sprintf("Электроэнергия Т1 (расход: %.2f кВт·ч)", newest.El1-baseline.El1))

	if newest.El2 != nil && property.El2Tariff != nil {
		var baseEl2 float64
		if baseline.El2 != nil {
			baseEl2 = *baseline.El2
		}
		addItem(*newest.El2-baseEl2, *property.El2Tariff, fmt.Sprintf("Электроэнергия Т2 (расход: %.2f кВт·ч)", *newest.El2-baseEl2))
	}

	return items, total, newest.ID, nil
}
