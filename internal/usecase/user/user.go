package user

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"github.com/google/uuid"
	entity "github.com/potom_pridumaem/internal/entity/users"
	"github.com/potom_pridumaem/internal/repo"
	"github.com/potom_pridumaem/internal/usecase"
)

type UseCase struct {
	repo repo.UserRepo
}

func New(r repo.UserRepo) *UseCase {
	return &UseCase{
		repo: r,
	}
}

func (uc *UseCase) Register(ctx context.Context, name, email, password, role, document, phone string) (entity.User, error) {
	if entity.Role(role) == entity.RoleAdmin {
		return entity.User{}, fmt.Errorf("UserUseCase - Register: %w", entity.ErrInvalidRole)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return entity.User{}, fmt.Errorf("UserUseCase - Register - bcrypt.GenerateFromPassword: %w", err)
	}

	user := entity.User{
		ID:           uuid.New().String(),
		Name:         name,
		Email:        email,
		PasswordHash: string(hash),
		Role:         entity.Role(role),
		Document:     document,
		Phone:        phone,
	}

	if err := user.Validate(); err != nil {
		return entity.User{}, fmt.Errorf("UserUseCase - Register - user.Validate: %w", err)
	}

	if err := uc.repo.Store(ctx, &user); err != nil {
		return entity.User{}, fmt.Errorf("UserUseCase - Register - uc.repo.Store: %w", err)
	}

	return user, nil
}

func (uc *UseCase) Login(ctx context.Context, email, password string) (entity.User, error) {
	user, err := uc.repo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repo.ErrUserNotFound) {
			return entity.User{}, usecase.ErrInvalidCredentials
		}
		return entity.User{}, fmt.Errorf("UserUseCase - Login - uc.repo.GetByEmail: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return entity.User{}, usecase.ErrInvalidCredentials
	}

	return user, nil
}
