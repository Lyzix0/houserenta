package user

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"github.com/google/uuid"
	"github.com/potom_pridumaem/internal/controller/v1/request"
	entity "github.com/potom_pridumaem/internal/entity/users"
	"github.com/potom_pridumaem/internal/repo"
	"github.com/potom_pridumaem/internal/usecase"
)

type UseCase struct {
	repo   repo.UserRepo
	leases repo.LeaseRepo
}

func New(r repo.UserRepo, leases repo.LeaseRepo) *UseCase {
	return &UseCase{
		repo:   r,
		leases: leases,
	}
}

func (uc *UseCase) Register(ctx context.Context, name, email, password, role, document, phone string, paymentCard *string) (entity.User, error) {
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
		PaymentCard:  paymentCard,
	}

	if err := user.Validate(); err != nil {
		return entity.User{}, fmt.Errorf("UserUseCase - Register - user.Validate: %w", err)
	}

	if err := uc.repo.Store(ctx, &user); err != nil {
		return entity.User{}, fmt.Errorf("UserUseCase - Register - uc.repo.Store: %w", err)
	}

	return user, nil
}

func (uc *UseCase) Login(ctx context.Context, identifier, password string) (entity.User, error) {
	user, err := uc.repo.GetByEmailOrPhone(ctx, identifier)
	if err != nil {
		if errors.Is(err, repo.ErrUserNotFound) {
			return entity.User{}, usecase.ErrInvalidCredentials
		}
		return entity.User{}, fmt.Errorf("UserUseCase - Login - uc.repo.GetByEmailOrPhone: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return entity.User{}, usecase.ErrInvalidCredentials
	}

	return user, nil
}

func (uc *UseCase) UpdateProfile(ctx context.Context, userID string, body request.Profile) error {
	user, err := uc.repo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("UserUseCase - UpdateProfile - uc.repo.GetByID: %w", err)
	}

	if body.Name != nil {
		user.Name = *body.Name
	}
	if body.Document != nil {
		user.Document = *body.Document
	}
	if body.Phone != nil {
		user.Phone = *body.Phone
	}
	if body.PaymentCard != nil {
		user.PaymentCard = body.PaymentCard
	}
	if body.Email != nil {
		user.Email = *body.Email
	}
	if body.Password != nil {
		hash, err := bcrypt.GenerateFromPassword([]byte(*body.Password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("UserUseCase - UpdateProfile - bcrypt.GenerateFromPassword: %w", err)
		}
		user.PasswordHash = string(hash)
	}

	if err := user.Validate(); err != nil {
		return fmt.Errorf("UserUseCase - UpdateProfile - user.Validate: %w", err)
	}

	if err := uc.repo.Update(ctx, &user); err != nil {
		return fmt.Errorf("UserUseCase - UpdateProfile - uc.repo.Update: %w", err)
	}

	return nil
}

func (uc *UseCase) Me(ctx context.Context, userID string) (usecase.UserProfile, error) {
	usr, err := uc.repo.GetByID(ctx, userID)
	if err != nil {
		return usecase.UserProfile{}, fmt.Errorf("UserUseCase - Me - uc.repo.GetByID: %w", err)
	}

	profile := usecase.UserProfile{User: usr}

	if usr.Role == entity.RoleTenant {
		lease, err := uc.leases.GetByTenantUserID(ctx, userID)
		switch {
		case err == nil:
			propertyID := lease.PropertyID
			profile.TenantPropertyID = &propertyID
		case errors.Is(err, repo.ErrLeaseNotFound):
			// tenant has no active lease yet: tenantPropertyId stays nil
		default:
			return usecase.UserProfile{}, fmt.Errorf("UserUseCase - Me - uc.leases.GetByTenantUserID: %w", err)
		}
	}

	return profile, nil
}
