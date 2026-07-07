package persistent

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	entity "github.com/potom_pridumaem/internal/entity/users"
	"github.com/potom_pridumaem/internal/repo"
	"github.com/potom_pridumaem/pkg/postgres"
)

type UserRepo struct {
	*postgres.Postgres
}

func NewUserRepo(pg *postgres.Postgres) *UserRepo {
	return &UserRepo{pg}
}

func (r *UserRepo) Store(ctx context.Context, user *entity.User) error {
	if err := user.Validate(); err != nil {
		return fmt.Errorf("UserRepo - Store - user.Validate: %w", err)
	}

	sql, args, err := r.Builder.
		Insert("app.users").
		Columns("id, name, email, password_hash, role, document, phone, payment_card").
		Values(
			user.ID,
			user.Name,
			user.Email,
			user.PasswordHash,
			user.Role,
			user.Document,
			user.Phone,
			user.PaymentCard,
		).
		ToSql()
	if err != nil {
		return fmt.Errorf("UserRepo - Store - r.Builder: %w", err)
	}

	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			switch pgErr.ConstraintName {
			case "users_email_key":
				return repo.ErrEmailAlreadyTaken
			default:
				return repo.ErrUserAlreadyExists
			}
		}
		return fmt.Errorf("UserRepo - Store - r.Pool.Exec: %w", err)
	}

	return nil
}

func (r *UserRepo) GetByEmail(
	ctx context.Context,
	email string,
) (entity.User, error) {
	sql, args, err := r.Builder.
		Select("id, name, email, password_hash, role, document, phone, payment_card").
		From("app.users").
		Where("email = ?", email).
		ToSql()
	if err != nil {
		return entity.User{}, fmt.Errorf("UserRepo - GetByEmail - r.Builder: %w", err)
	}

	var u entity.User
	err = r.Pool.QueryRow(ctx, sql, args...).Scan(
		&u.ID,
		&u.Name,
		&u.Email,
		&u.PasswordHash,
		&u.Role,
		&u.Document,
		&u.Phone,
		&u.PaymentCard,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.User{}, repo.ErrUserNotFound
		}
		return entity.User{}, fmt.Errorf("UserRepo - GetByEmail - r.Pool.QueryRow: %w", err)
	}

	return u, nil
}

func (r *UserRepo) GetByID(
	ctx context.Context,
	id string,
) (entity.User, error) {
	sql, args, err := r.Builder.
		Select("id, name, email, password_hash, role, document, phone, payment_card").
		From("app.users").
		Where("id = ?", id).
		ToSql()
	if err != nil {
		return entity.User{}, fmt.Errorf("UserRepo - GetByID - r.Builder: %w", err)
	}

	var u entity.User
	err = r.Pool.QueryRow(ctx, sql, args...).Scan(
		&u.ID,
		&u.Name,
		&u.Email,
		&u.PasswordHash,
		&u.Role,
		&u.Document,
		&u.Phone,
		&u.PaymentCard,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.User{}, repo.ErrUserNotFound
		}
		return entity.User{}, fmt.Errorf("UserRepo - GetByID - r.Pool.QueryRow: %w", err)
	}

	return u, nil
}
