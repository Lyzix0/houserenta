package v1_test

import (
	"context"

	"github.com/potom_pridumaem/internal/controller/v1/request"
	entity "github.com/potom_pridumaem/internal/entity/users"
	"github.com/potom_pridumaem/internal/usecase"
)

type userUseCaseMock struct {
	registerFn func(ctx context.Context, name, email, password, role, document, phone string) (entity.User, error)
	loginFn    func(ctx context.Context, email, password string) (entity.User, error)
	meFn       func(ctx context.Context, userID string) (usecase.UserProfile, error)
}

func (m *userUseCaseMock) Register(ctx context.Context, name, email, password, role, document, phone string) (entity.User, error) {
	return m.registerFn(ctx, name, email, password, role, document, phone)
}

func (m *userUseCaseMock) Login(ctx context.Context, email, password string) (entity.User, error) {
	return m.loginFn(ctx, email, password)
}

func (m *userUseCaseMock) Me(ctx context.Context, userID string) (usecase.UserProfile, error) {
	return m.meFn(ctx, userID)
}

type propertyUseCaseMock struct {
	createPropertyFn func(ctx context.Context, body request.Property) (entity.Property, error)
	getPropertiesFn  func(ctx context.Context, landlordID string) ([]entity.Property, error)
	getPropertyFn    func(ctx context.Context, id, landlordID string) (entity.Property, error)
	updatePropertyFn func(ctx context.Context, id, landlordID string, body request.Property) error
	deletePropertyFn func(ctx context.Context, id, landlordID string) error
}

func (m *propertyUseCaseMock) CreateProperty(ctx context.Context, body request.Property) (entity.Property, error) {
	return m.createPropertyFn(ctx, body)
}

func (m *propertyUseCaseMock) GetProperties(ctx context.Context, landlordID string) ([]entity.Property, error) {
	return m.getPropertiesFn(ctx, landlordID)
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
