package v1_test

import (
	"context"

	"github.com/potom_pridumaem/internal/controller/v1/request"
	entity "github.com/potom_pridumaem/internal/entity/users"
	"github.com/potom_pridumaem/internal/usecase"
)

type userUseCaseMock struct {
	registerFn      func(ctx context.Context, name, email, password, role, document, phone string, paymentCard *string) (entity.User, error)
	loginFn         func(ctx context.Context, identifier, password string) (entity.User, error)
	meFn            func(ctx context.Context, userID string) (usecase.UserProfile, error)
	updateProfileFn func(ctx context.Context, userID string, body request.Profile) error
}

func (m *userUseCaseMock) Register(ctx context.Context, name, email, password, role, document, phone string, paymentCard *string) (entity.User, error) {
	return m.registerFn(ctx, name, email, password, role, document, phone, paymentCard)
}

func (m *userUseCaseMock) Login(ctx context.Context, identifier, password string) (entity.User, error) {
	return m.loginFn(ctx, identifier, password)
}

func (m *userUseCaseMock) Me(ctx context.Context, userID string) (usecase.UserProfile, error) {
	return m.meFn(ctx, userID)
}

func (m *userUseCaseMock) UpdateProfile(ctx context.Context, userID string, body request.Profile) error {
	return m.updateProfileFn(ctx, userID, body)
}

type propertyUseCaseMock struct {
	createPropertyFn   func(ctx context.Context, landlordID string, body request.Property) (entity.Property, error)
	getPropertiesFn    func(ctx context.Context, userID string, role entity.Role) ([]entity.PropertyDetail, error)
	getPropertyFn      func(ctx context.Context, id, landlordID string) (entity.Property, error)
	updatePropertyFn   func(ctx context.Context, id, landlordID string, body request.Property) error
	deletePropertyFn   func(ctx context.Context, id, landlordID string) error
	createLeaseFn      func(ctx context.Context, propertyID, landlordID string, body request.Lease) error
	deleteLeaseFn      func(ctx context.Context, propertyID, landlordID string) error
	createReadingFn    func(ctx context.Context, propertyID, userID string, role entity.Role, body request.Reading) error
	payFn              func(ctx context.Context, propertyID, userID string, role entity.Role, body request.Payment) error
	createCustomItemFn func(ctx context.Context, propertyID, landlordID string, body request.CustomItem) error
}

func (m *propertyUseCaseMock) CreateProperty(ctx context.Context, landlordID string, body request.Property) (entity.Property, error) {
	return m.createPropertyFn(ctx, landlordID, body)
}

func (m *propertyUseCaseMock) GetProperties(ctx context.Context, userID string, role entity.Role) ([]entity.PropertyDetail, error) {
	return m.getPropertiesFn(ctx, userID, role)
}

func (m *propertyUseCaseMock) GetProperty(ctx context.Context, id, landlordID string) (entity.Property, error) {
	return m.getPropertyFn(ctx, id, landlordID)
}

func (m *propertyUseCaseMock) UpdateProperty(ctx context.Context, id, landlordID string, body request.Property) error {
	return m.updatePropertyFn(ctx, id, landlordID, body)
}

func (m *propertyUseCaseMock) DeleteProperty(ctx context.Context, id, landlordID string) error {
	return m.deletePropertyFn(ctx, id, landlordID)
}

func (m *propertyUseCaseMock) CreateLease(ctx context.Context, propertyID, landlordID string, body request.Lease) error {
	return m.createLeaseFn(ctx, propertyID, landlordID, body)
}

func (m *propertyUseCaseMock) DeleteLease(ctx context.Context, propertyID, landlordID string) error {
	return m.deleteLeaseFn(ctx, propertyID, landlordID)
}

func (m *propertyUseCaseMock) CreateReading(ctx context.Context, propertyID, userID string, role entity.Role, body request.Reading) error {
	return m.createReadingFn(ctx, propertyID, userID, role, body)
}

func (m *propertyUseCaseMock) Pay(ctx context.Context, propertyID, userID string, role entity.Role, body request.Payment) error {
	return m.payFn(ctx, propertyID, userID, role, body)
}

func (m *propertyUseCaseMock) CreateCustomItem(ctx context.Context, propertyID, landlordID string, body request.CustomItem) error {
	return m.createCustomItemFn(ctx, propertyID, landlordID, body)
}
