package v1_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/potom_pridumaem/internal/controller/v1/request"
	entity "github.com/potom_pridumaem/internal/entity/users"
	"github.com/potom_pridumaem/internal/repo"
)

func TestCreateProperty(t *testing.T) {
	validBody := request.Property{
		LandlordID:  "landlord-1",
		Name:        "Sunny Apartment",
		Coordinates: "55.75,37.61",
		Region:      "Moscow",
		City:        "Moscow",
		Street:      "Tverskaya",
		House:       "1",
		Apartment:   "42",
		GvsTariff:   1.5,
		HvsTariff:   1.2,
		El1Tariff:   4.5,
	}

	tests := []struct {
		name             string
		body             any
		createPropertyFn func(ctx context.Context, body request.Property) (entity.Property, error)
		wantStatus       int
	}{
		{
			name: "success",
			body: validBody,
			createPropertyFn: func(_ context.Context, body request.Property) (entity.Property, error) {
				return entity.Property{ID: "prop-1", LandlordID: body.LandlordID, Name: body.Name}, nil
			},
			wantStatus: http.StatusCreated,
		},
		{
			name:             "invalid body",
			body:             map[string]string{"name": "missing required fields"},
			createPropertyFn: func(context.Context, request.Property) (entity.Property, error) { return entity.Property{}, nil },
			wantStatus:       http.StatusBadRequest,
		},
		{
			name: "landlord not found",
			body: validBody,
			createPropertyFn: func(context.Context, request.Property) (entity.Property, error) {
				return entity.Property{}, repo.ErrLandlordNotFound
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name: "property already exists",
			body: validBody,
			createPropertyFn: func(context.Context, request.Property) (entity.Property, error) {
				return entity.Property{}, repo.ErrPropertyAlreadyExists
			},
			wantStatus: http.StatusConflict,
		},
		{
			name: "internal error",
			body: validBody,
			createPropertyFn: func(context.Context, request.Property) (entity.Property, error) {
				return entity.Property{}, errors.New("boom")
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := newTestApp(&userUseCaseMock{}, &propertyUseCaseMock{createPropertyFn: tt.createPropertyFn})

			resp := doRequest(t, app, http.MethodPost, "/v1/properties/property", tt.body)
			defer resp.Body.Close()

			if resp.StatusCode != tt.wantStatus {
				t.Fatalf("status = %d, want %d", resp.StatusCode, tt.wantStatus)
			}
		})
	}
}
