package billing_test

import (
	"context"

	entity "github.com/potom_pridumaem/internal/entity/users"
)

type leaseRepoMock struct {
	getAllFn func(ctx context.Context) ([]entity.Lease, error)
}

func (m *leaseRepoMock) GetByTenantUserID(context.Context, string) (entity.Lease, error) {
	return entity.Lease{}, nil
}
func (m *leaseRepoMock) GetByPropertyID(context.Context, string) (entity.Lease, error) {
	return entity.Lease{}, nil
}
func (m *leaseRepoMock) GetAll(ctx context.Context) ([]entity.Lease, error) { return m.getAllFn(ctx) }
func (m *leaseRepoMock) Upsert(context.Context, entity.Lease) error         { return nil }
func (m *leaseRepoMock) DeleteByPropertyID(context.Context, string) error   { return nil }

type propertyRepoMock struct {
	getByIDFn    func(ctx context.Context, id string) (entity.Property, error)
	addBalanceFn func(ctx context.Context, propertyID string, amount float64) error
}

func (m *propertyRepoMock) Store(context.Context, entity.Property) error { return nil }
func (m *propertyRepoMock) GetByLandlordID(context.Context, string) ([]entity.Property, error) {
	return nil, nil
}
func (m *propertyRepoMock) GetByID(ctx context.Context, id string) (entity.Property, error) {
	return m.getByIDFn(ctx, id)
}
func (m *propertyRepoMock) Update(context.Context, entity.Property) error { return nil }
func (m *propertyRepoMock) Delete(context.Context, string) error          { return nil }
func (m *propertyRepoMock) AddBalance(ctx context.Context, propertyID string, amount float64) error {
	return m.addBalanceFn(ctx, propertyID, amount)
}
func (m *propertyRepoMock) GetVacant(context.Context) ([]entity.Property, error) { return nil, nil }

type billRepoMock struct {
	getLastFn func(ctx context.Context, propertyID, billType string) (entity.Bill, error)
	storeFn   func(ctx context.Context, bill entity.Bill) error
}

func (m *billRepoMock) GetByPropertyID(context.Context, string) ([]entity.Bill, error) {
	return nil, nil
}
func (m *billRepoMock) UpdateStatus(context.Context, string, string, string) error { return nil }
func (m *billRepoMock) GetLastByPropertyIDAndType(ctx context.Context, propertyID, billType string) (entity.Bill, error) {
	return m.getLastFn(ctx, propertyID, billType)
}
func (m *billRepoMock) Store(ctx context.Context, bill entity.Bill) error {
	return m.storeFn(ctx, bill)
}

type readingRepoMock struct {
	getOldestUnaccountedFn func(ctx context.Context, propertyID string) (entity.Reading, error)
	getLatestAccountedFn   func(ctx context.Context, propertyID string) (entity.Reading, error)
	markAccountedFn        func(ctx context.Context, id string) error
}

func (m *readingRepoMock) GetByPropertyID(context.Context, string) ([]entity.Reading, error) {
	return nil, nil
}
func (m *readingRepoMock) Store(context.Context, entity.Reading) error { return nil }
func (m *readingRepoMock) GetOldestUnaccounted(ctx context.Context, propertyID string) (entity.Reading, error) {
	return m.getOldestUnaccountedFn(ctx, propertyID)
}
func (m *readingRepoMock) GetLatestAccounted(ctx context.Context, propertyID string) (entity.Reading, error) {
	return m.getLatestAccountedFn(ctx, propertyID)
}
func (m *readingRepoMock) MarkAccounted(ctx context.Context, id string) error {
	return m.markAccountedFn(ctx, id)
}

type customNextItemRepoMock struct {
	getByPropertyIDFn  func(ctx context.Context, propertyID string) ([]entity.CustomNextItem, error)
	deleteByPropertyID func(ctx context.Context, propertyID string) error
}

func (m *customNextItemRepoMock) GetByPropertyID(ctx context.Context, propertyID string) ([]entity.CustomNextItem, error) {
	return m.getByPropertyIDFn(ctx, propertyID)
}
func (m *customNextItemRepoMock) Store(context.Context, entity.CustomNextItem) error { return nil }
func (m *customNextItemRepoMock) DeleteByPropertyID(ctx context.Context, propertyID string) error {
	return m.deleteByPropertyID(ctx, propertyID)
}
