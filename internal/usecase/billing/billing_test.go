package billing_test

import (
	"context"
	"errors"
	"testing"
	"time"

	entity "github.com/potom_pridumaem/internal/entity/users"
	"github.com/potom_pridumaem/internal/repo"
	"github.com/potom_pridumaem/internal/usecase/billing"
)

func TestRun_NoLeases_NoOp(t *testing.T) {
	uc := billing.New(
		&leaseRepoMock{getAllFn: func(context.Context) ([]entity.Lease, error) { return nil, nil }},
		&propertyRepoMock{},
		&billRepoMock{},
		&readingRepoMock{},
		&customNextItemRepoMock{},
	)

	if err := uc.Run(context.Background()); err != nil {
		t.Fatalf("Run() error = %v, want nil", err)
	}
}

func TestRun_SkipsWhenNotDue(t *testing.T) {
	lease := entity.Lease{PropertyID: "prop-1", Price: 30000}
	recentBill := entity.Bill{Date: time.Now().UTC().Format(time.RFC3339)}

	var stored bool
	uc := billing.New(
		&leaseRepoMock{getAllFn: func(context.Context) ([]entity.Lease, error) { return []entity.Lease{lease}, nil }},
		&propertyRepoMock{},
		&billRepoMock{
			getLastFn: func(context.Context, string, string) (entity.Bill, error) { return recentBill, nil },
			storeFn:   func(context.Context, entity.Bill) error { stored = true; return nil },
		},
		&readingRepoMock{},
		&customNextItemRepoMock{},
	)

	if err := uc.Run(context.Background()); err != nil {
		t.Fatalf("Run() error = %v, want nil", err)
	}
	if stored {
		t.Fatal("expected no bill to be stored: last bill is recent")
	}
}

func TestRun_GeneratesFirstBill_NoPriorBillNoReadings(t *testing.T) {
	lease := entity.Lease{PropertyID: "prop-1", Price: 30000}
	property := entity.Property{ID: "prop-1", LandlordID: "landlord-1", GvsTariff: 200, HvsTariff: 50, El1Tariff: 6}

	var storedBill entity.Bill
	var balanceDelta float64
	var deletedCustomItems bool

	uc := billing.New(
		&leaseRepoMock{getAllFn: func(context.Context) ([]entity.Lease, error) { return []entity.Lease{lease}, nil }},
		&propertyRepoMock{
			getByIDFn: func(_ context.Context, id string) (entity.Property, error) {
				if id != "prop-1" {
					t.Fatalf("GetByID id = %q, want prop-1", id)
				}
				return property, nil
			},
			addBalanceFn: func(_ context.Context, propertyID string, amount float64) error {
				balanceDelta = amount
				return nil
			},
		},
		&billRepoMock{
			getLastFn: func(context.Context, string, string) (entity.Bill, error) {
				return entity.Bill{}, repo.ErrBillNotFound
			},
			storeFn: func(_ context.Context, bill entity.Bill) error {
				storedBill = bill
				return nil
			},
		},
		&readingRepoMock{
			getOldestUnaccountedFn: func(context.Context, string) (entity.Reading, error) {
				return entity.Reading{}, repo.ErrReadingNotFound
			},
		},
		&customNextItemRepoMock{
			getByPropertyIDFn: func(context.Context, string) ([]entity.CustomNextItem, error) { return nil, nil },
			deleteByPropertyID: func(context.Context, string) error {
				deletedCustomItems = true
				return nil
			},
		},
	)

	if err := uc.Run(context.Background()); err != nil {
		t.Fatalf("Run() error = %v, want nil", err)
	}

	if storedBill.ID == "" {
		t.Fatal("expected a bill to be stored")
	}
	if storedBill.Type != "rent" {
		t.Fatalf("bill.Type = %q, want %q", storedBill.Type, "rent")
	}
	if storedBill.Status != "unpaid" {
		t.Fatalf("bill.Status = %q, want %q", storedBill.Status, "unpaid")
	}
	if storedBill.Total != lease.Price {
		t.Fatalf("bill.Total = %v, want %v (rent only, no readings/custom items)", storedBill.Total, lease.Price)
	}
	if len(storedBill.Items) != 1 {
		t.Fatalf("len(bill.Items) = %d, want 1 (just the base rent line)", len(storedBill.Items))
	}
	if balanceDelta != -lease.Price {
		t.Fatalf("balance delta = %v, want %v", balanceDelta, -lease.Price)
	}
	if deletedCustomItems {
		t.Fatal("did not expect custom items deletion: there were none")
	}
}

func TestRun_AttachesAndClearsCustomItems(t *testing.T) {
	lease := entity.Lease{PropertyID: "prop-1", Price: 30000}
	property := entity.Property{ID: "prop-1"}
	customItems := []entity.CustomNextItem{
		{ID: "ci-1", PropertyID: "prop-1", Description: "Замена замка", Amount: 1500},
	}

	var storedBill entity.Bill
	var deletedCustomItems bool

	uc := billing.New(
		&leaseRepoMock{getAllFn: func(context.Context) ([]entity.Lease, error) { return []entity.Lease{lease}, nil }},
		&propertyRepoMock{
			getByIDFn:    func(context.Context, string) (entity.Property, error) { return property, nil },
			addBalanceFn: func(context.Context, string, float64) error { return nil },
		},
		&billRepoMock{
			getLastFn: func(context.Context, string, string) (entity.Bill, error) {
				return entity.Bill{}, repo.ErrBillNotFound
			},
			storeFn: func(_ context.Context, bill entity.Bill) error { storedBill = bill; return nil },
		},
		&readingRepoMock{
			getOldestUnaccountedFn: func(context.Context, string) (entity.Reading, error) {
				return entity.Reading{}, repo.ErrReadingNotFound
			},
		},
		&customNextItemRepoMock{
			getByPropertyIDFn: func(context.Context, string) ([]entity.CustomNextItem, error) { return customItems, nil },
			deleteByPropertyID: func(context.Context, string) error {
				deletedCustomItems = true
				return nil
			},
		},
	)

	if err := uc.Run(context.Background()); err != nil {
		t.Fatalf("Run() error = %v, want nil", err)
	}

	wantTotal := lease.Price + customItems[0].Amount
	if storedBill.Total != wantTotal {
		t.Fatalf("bill.Total = %v, want %v", storedBill.Total, wantTotal)
	}
	if len(storedBill.Items) != 2 {
		t.Fatalf("len(bill.Items) = %d, want 2 (rent + custom item)", len(storedBill.Items))
	}
	if !deletedCustomItems {
		t.Fatal("expected custom items to be cleared after being attached to the bill")
	}
}

func TestRun_ComputesUtilityDiffAgainstBaselineAndMarksAccounted(t *testing.T) {
	lease := entity.Lease{PropertyID: "prop-1", Price: 30000}
	property := entity.Property{ID: "prop-1", GvsTariff: 200, HvsTariff: 50, El1Tariff: 6}

	baseline := entity.Reading{ID: "read-old", Gvs: 10, Hvs: 20, El1: 300}
	newest := entity.Reading{ID: "read-new", Gvs: 12.5, Hvs: 24.1, El1: 340}

	var storedBill entity.Bill
	var markedAccountedID string

	uc := billing.New(
		&leaseRepoMock{getAllFn: func(context.Context) ([]entity.Lease, error) { return []entity.Lease{lease}, nil }},
		&propertyRepoMock{
			getByIDFn:    func(context.Context, string) (entity.Property, error) { return property, nil },
			addBalanceFn: func(context.Context, string, float64) error { return nil },
		},
		&billRepoMock{
			getLastFn: func(context.Context, string, string) (entity.Bill, error) {
				return entity.Bill{}, repo.ErrBillNotFound
			},
			storeFn: func(_ context.Context, bill entity.Bill) error { storedBill = bill; return nil },
		},
		&readingRepoMock{
			getOldestUnaccountedFn: func(context.Context, string) (entity.Reading, error) { return newest, nil },
			getLatestAccountedFn:   func(context.Context, string) (entity.Reading, error) { return baseline, nil },
			markAccountedFn: func(_ context.Context, id string) error {
				markedAccountedID = id
				return nil
			},
		},
		&customNextItemRepoMock{
			getByPropertyIDFn: func(context.Context, string) ([]entity.CustomNextItem, error) { return nil, nil },
		},
	)

	if err := uc.Run(context.Background()); err != nil {
		t.Fatalf("Run() error = %v, want nil", err)
	}

	gvsAmount := (newest.Gvs - baseline.Gvs) * property.GvsTariff
	hvsAmount := (newest.Hvs - baseline.Hvs) * property.HvsTariff
	el1Amount := (newest.El1 - baseline.El1) * property.El1Tariff
	wantTotal := lease.Price + gvsAmount + hvsAmount + el1Amount

	if storedBill.Total != wantTotal {
		t.Fatalf("bill.Total = %v, want %v", storedBill.Total, wantTotal)
	}
	if len(storedBill.Items) != 4 {
		t.Fatalf("len(bill.Items) = %d, want 4 (rent + gvs + hvs + el1)", len(storedBill.Items))
	}
	if markedAccountedID != newest.ID {
		t.Fatalf("marked accounted id = %q, want %q", markedAccountedID, newest.ID)
	}
}

func TestRun_PropagatesRepoErrors(t *testing.T) {
	uc := billing.New(
		&leaseRepoMock{getAllFn: func(context.Context) ([]entity.Lease, error) { return nil, errors.New("db down") }},
		&propertyRepoMock{},
		&billRepoMock{},
		&readingRepoMock{},
		&customNextItemRepoMock{},
	)

	if err := uc.Run(context.Background()); err == nil {
		t.Fatal("expected Run() to propagate the repo error")
	}
}
