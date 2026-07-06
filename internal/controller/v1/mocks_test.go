package v1_test

import (
	"context"

	"github.com/potom_pridumaem/internal/controller/v1/request"
	entity "github.com/potom_pridumaem/internal/entity/users"
)

type userUseCaseMock struct {
	registerFn func(ctx context.Context, name, email, password, role, document, phone string) (entity.User, error)
	loginFn    func(ctx context.Context, email, password string) (entity.User, error)
}

func (m *userUseCaseMock) Register(ctx context.Context, name, email, password, role, document, phone string) (entity.User, error) {
	return m.registerFn(ctx, name, email, password, role, document, phone)
}

func (m *userUseCaseMock) Login(ctx context.Context, email, password string) (entity.User, error) {
	return m.loginFn(ctx, email, password)
}

type propertyUseCaseMock struct {
	createPropertyFn func(ctx context.Context, body request.Property) (entity.Property, error)
}

func (m *propertyUseCaseMock) CreateProperty(ctx context.Context, body request.Property) (entity.Property, error) {
	return m.createPropertyFn(ctx, body)
}
