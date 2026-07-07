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

func (uc *UseCase) GetByID(ctx context.Context, id string) (entity.User, error) {
	return uc.repo.GetByID(ctx, id)
}

func (uc *UseCase) GetByEmail(ctx context.Context, email string) (entity.User, error) {
	return uc.repo.GetByEmail(ctx, email)
}

func (uc *UseCase) UpdateProfile(ctx context.Context, id string, name, document, phone, email *string, paymentCard *string) error {
	user, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if name != nil && *name != "" {
		user.Name = *name
	}
	if document != nil && *document != "" {
		user.Document = *document
	}
	if phone != nil && *phone != "" {
		user.Phone = *phone
	}
	if email != nil && *email != "" {
		user.Email = *email
	}
	if paymentCard != nil {
		if *paymentCard == "" {
			user.PaymentCard = nil
		} else {
			user.PaymentCard = paymentCard
		}
	}
	return uc.repo.Update(ctx, &user)
}

func (uc *UseCase) ChangePassword(ctx context.Context, id, oldPassword, newPassword string) error {
	user, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(oldPassword)); err != nil {
		return usecase.ErrInvalidCredentials
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return uc.repo.UpdatePassword(ctx, id, string(hash))
}

func (uc *UseCase) SwitchRole(ctx context.Context, id, targetRole string) (entity.User, error) {
	r := entity.Role(targetRole)
	if r != entity.RoleLandlord && r != entity.RoleTenant && r != entity.RoleAdmin {
		return entity.User{}, entity.ErrInvalidRole
	}
	user, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return entity.User{}, err
	}
	user.Role = r
	if err := uc.repo.Update(ctx, &user); err != nil {
		return entity.User{}, err
	}
	return user, nil
}
