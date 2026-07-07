package persistent

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	entity "github.com/potom_pridumaem/internal/entity/users"
	"github.com/potom_pridumaem/internal/repo"
	"github.com/potom_pridumaem/pkg/postgres"
)

type PropertyRepo struct {
	*postgres.Postgres
}

func NewPropertyRepo(pg *postgres.Postgres) *PropertyRepo {
	return &PropertyRepo{pg}
}

func (r *PropertyRepo) Store(ctx context.Context, prop entity.Property) error {
	sql, args, err := r.Builder.
		Insert("app.properties").
		Columns(
			"id", "landlord_id", "name", "coordinates", "country",
			"region", "city", "street", "house", "apartment",
			"gvs_tariff", "hvs_tariff", "el1_tariff", "el2_tariff", "balance",
		).
		Values(
			prop.ID, prop.LandlordID, prop.Name, prop.Coordinates, prop.Country,
			prop.Region, prop.City, prop.Street, prop.House, prop.Apartment,
			prop.GvsTariff, prop.HvsTariff, prop.El1Tariff, prop.El2Tariff, prop.Balance,
		).
		ToSql()
	if err != nil {
		return fmt.Errorf("PropertyRepo - Store - r.Builder: %w", err)
	}

	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505":
				return repo.ErrPropertyAlreadyExists
			case "23503":
				return repo.ErrLandlordNotFound
			}
		}
		return fmt.Errorf("PropertyRepo - Store - r.Pool.Exec: %w", err)
	}

	return nil
}

func (r *PropertyRepo) GetByLandlordID(ctx context.Context, landlordID string) ([]entity.Property, error) {
	sql, args, err := r.Builder.
		Select(
			"id", "landlord_id", "name", "coordinates", "country",
			"region", "city", "street", "house", "apartment",
			"gvs_tariff", "hvs_tariff", "el1_tariff", "el2_tariff", "balance",
		).
		From("app.properties").
		Where("landlord_id = ?", landlordID).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("PropertyRepo - GetByLandlordID - r.Builder: %w", err)
	}

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("PropertyRepo - GetByLandlordID - r.Pool.Query: %w", err)
	}
	defer rows.Close()

	var properties []entity.Property
	for rows.Next() {
		var prop entity.Property
		if err := rows.Scan(
			&prop.ID, &prop.LandlordID, &prop.Name, &prop.Coordinates, &prop.Country,
			&prop.Region, &prop.City, &prop.Street, &prop.House, &prop.Apartment,
			&prop.GvsTariff, &prop.HvsTariff, &prop.El1Tariff, &prop.El2Tariff, &prop.Balance,
		); err != nil {
			return nil, fmt.Errorf("PropertyRepo - GetByLandlordID - rows.Scan: %w", err)
		}
		properties = append(properties, prop)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("PropertyRepo - GetByLandlordID - rows.Err: %w", err)
	}

	return properties, nil
}

func (r *PropertyRepo) GetByID(ctx context.Context, id string) (entity.Property, error) {
	sql, args, err := r.Builder.
		Select(
			"id", "landlord_id", "name", "coordinates", "country",
			"region", "city", "street", "house", "apartment",
			"gvs_tariff", "hvs_tariff", "el1_tariff", "el2_tariff", "balance",
		).
		From("app.properties").
		Where("id = ?", id).
		ToSql()
	if err != nil {
		return entity.Property{}, fmt.Errorf("PropertyRepo - GetByID - r.Builder: %w", err)
	}

	var prop entity.Property
	err = r.Pool.QueryRow(ctx, sql, args...).Scan(
		&prop.ID, &prop.LandlordID, &prop.Name, &prop.Coordinates, &prop.Country,
		&prop.Region, &prop.City, &prop.Street, &prop.House, &prop.Apartment,
		&prop.GvsTariff, &prop.HvsTariff, &prop.El1Tariff, &prop.El2Tariff, &prop.Balance,
	)
	if err != nil {
		return entity.Property{}, fmt.Errorf("PropertyRepo - GetByID - r.Pool.QueryRow: %w", err)
	}

	return prop, nil
}

func (r *PropertyRepo) Update(ctx context.Context, prop *entity.Property) error {
	sql, args, err := r.Builder.
		Update("app.properties").
		Set("landlord_id", prop.LandlordID).
		Set("name", prop.Name).
		Set("coordinates", prop.Coordinates).
		Set("country", prop.Country).
		Set("region", prop.Region).
		Set("city", prop.City).
		Set("street", prop.Street).
		Set("house", prop.House).
		Set("apartment", prop.Apartment).
		Set("gvs_tariff", prop.GvsTariff).
		Set("hvs_tariff", prop.HvsTariff).
		Set("el1_tariff", prop.El1Tariff).
		Set("el2_tariff", prop.El2Tariff).
		Set("balance", prop.Balance).
		Where("id = ?", prop.ID).
		ToSql()
	if err != nil {
		return fmt.Errorf("PropertyRepo - Update - r.Builder: %w", err)
	}

	result, err := r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("PropertyRepo - Update - r.Pool.Exec: %w", err)
	}

	if result.RowsAffected() == 0 {
		return repo.ErrPropertyNotFound
	}

	return nil
}

func (r *PropertyRepo) Delete(ctx context.Context, id string) error {
	sql, args, err := r.Builder.
		Delete("app.properties").
		Where("id = ?", id).
		ToSql()
	if err != nil {
		return fmt.Errorf("PropertyRepo - Delete - r.Builder: %w", err)
	}

	result, err := r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("PropertyRepo - Delete - r.Pool.Exec: %w", err)
	}

	if result.RowsAffected() == 0 {
		return repo.ErrPropertyNotFound
	}

	return nil
}

func (r *PropertyRepo) GetVacant(ctx context.Context) ([]entity.Property, error) {
	sql, args, err := r.Builder.
		Select(
			"id", "landlord_id", "name", "coordinates", "country",
			"region", "city", "street", "house", "apartment",
			"gvs_tariff", "hvs_tariff", "el1_tariff", "el2_tariff", "balance",
		).
		From("app.properties").
		Where("id NOT IN (SELECT property_id FROM app.leases)").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("PropertyRepo - GetVacant - r.Builder: %w", err)
	}

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("PropertyRepo - GetVacant - r.Pool.Query: %w", err)
	}
	defer rows.Close()

	var properties []entity.Property
	for rows.Next() {
		var prop entity.Property
		if err := rows.Scan(
			&prop.ID, &prop.LandlordID, &prop.Name, &prop.Coordinates, &prop.Country,
			&prop.Region, &prop.City, &prop.Street, &prop.House, &prop.Apartment,
			&prop.GvsTariff, &prop.HvsTariff, &prop.El1Tariff, &prop.El2Tariff, &prop.Balance,
		); err != nil {
			return nil, fmt.Errorf("PropertyRepo - GetVacant - rows.Scan: %w", err)
		}
		properties = append(properties, prop)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("PropertyRepo - GetVacant - rows.Err: %w", err)
	}

	return properties, nil
}

type UnlinkedTenant struct {
	ID       string
	Name     string
	Email    string
	Document string
	Phone    string
}

func (r *PropertyRepo) GetUnlinkedTenants(ctx context.Context) ([]UnlinkedTenant, error) {
	rows, err := r.Pool.Query(ctx,
		"SELECT id, name, email, document, phone FROM app.users WHERE role = 'tenant'")
	if err != nil {
		return nil, fmt.Errorf("PropertyRepo - GetUnlinkedTenants: %w", err)
	}
	defer rows.Close()

	var tenants []UnlinkedTenant
	for rows.Next() {
		var t UnlinkedTenant
		if err := rows.Scan(&t.ID, &t.Name, &t.Email, &t.Document, &t.Phone); err != nil {
			return nil, fmt.Errorf("PropertyRepo - GetUnlinkedTenants - scan: %w", err)
		}
		tenants = append(tenants, t)
	}
	return tenants, nil
}

type ChatMessage struct {
	ID         string `json:"id"`
	SenderID   string `json:"sender_id"`
	ReceiverID string `json:"receiver_id"`
	PropertyID string `json:"property_id"`
	Text       string `json:"text"`
	Timestamp  string `json:"timestamp"`
}

func (r *PropertyRepo) GetChatHistory(ctx context.Context, propertyID string) ([]ChatMessage, error) {
	rows, err := r.Pool.Query(ctx,
		"SELECT id, sender_id, receiver_id, property_id, text, timestamp FROM app.chat_messages WHERE property_id = $1 ORDER BY timestamp ASC", propertyID)
	if err != nil {
		return nil, fmt.Errorf("PropertyRepo - GetChatHistory: %w", err)
	}
	defer rows.Close()

	var msgs []ChatMessage
	for rows.Next() {
		var m ChatMessage
		if err := rows.Scan(&m.ID, &m.SenderID, &m.ReceiverID, &m.PropertyID, &m.Text, &m.Timestamp); err != nil {
			return nil, fmt.Errorf("PropertyRepo - GetChatHistory - scan: %w", err)
		}
		msgs = append(msgs, m)
	}
	return msgs, nil
}

func (r *PropertyRepo) SaveChatMessage(ctx context.Context, msg *ChatMessage) error {
	_, err := r.Pool.Exec(ctx,
		"INSERT INTO app.chat_messages (id, sender_id, receiver_id, property_id, text, timestamp) VALUES ($1,$2,$3,$4,$5,$6)",
		msg.ID, msg.SenderID, msg.ReceiverID, msg.PropertyID, msg.Text, msg.Timestamp)
	return err
}
