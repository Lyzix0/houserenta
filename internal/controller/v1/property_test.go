package v1_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v3"
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

func TestGetProperties(t *testing.T) {
	t.Run("unauthenticated", func(t *testing.T) {
		app := newTestApp(&userUseCaseMock{}, &propertyUseCaseMock{})

		resp := doRequest(t, app, http.MethodGet, "/v1/properties/", nil)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnauthorized {
			t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusUnauthorized)
		}
	})

	t.Run("authenticated", func(t *testing.T) {
		app := newTestApp(&userUseCaseMock{
			loginFn: func(_ context.Context, email, _ string) (entity.User, error) {
				return entity.User{ID: "landlord-1", Email: email, Role: entity.RoleLandlord}, nil
			},
		}, &propertyUseCaseMock{
			getPropertiesFn: func(_ context.Context, landlordID string) ([]entity.Property, error) {
				return []entity.Property{{ID: "prop-1", LandlordID: landlordID, Name: "Sunny Apartment"}}, nil
			},
		})

		cookie := loginAndGetCookie(t, app)

		resp := doRequest(t, app, http.MethodGet, "/v1/properties/", nil, cookie)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusOK)
		}

		var body []entity.Property
		decodeJSON(t, resp, &body)

		if len(body) != 1 || body[0].LandlordID != "landlord-1" {
			t.Fatalf("body = %+v, want a single property for landlord-1", body)
		}
	})

	t.Run("internal error", func(t *testing.T) {
		app := newTestApp(&userUseCaseMock{
			loginFn: func(_ context.Context, email, _ string) (entity.User, error) {
				return entity.User{ID: "landlord-1", Email: email, Role: entity.RoleLandlord}, nil
			},
		}, &propertyUseCaseMock{
			getPropertiesFn: func(context.Context, string) ([]entity.Property, error) {
				return nil, errors.New("boom")
			},
		})

		cookie := loginAndGetCookie(t, app)

		resp := doRequest(t, app, http.MethodGet, "/v1/properties/", nil, cookie)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusInternalServerError {
			t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusInternalServerError)
		}
	})
}

func TestGetProperty(t *testing.T) {
	t.Run("unauthenticated", func(t *testing.T) {
		app := newTestApp(&userUseCaseMock{}, &propertyUseCaseMock{})

		resp := doRequest(t, app, http.MethodGet, "/v1/properties/prop-1", nil)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnauthorized {
			t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusUnauthorized)
		}
	})

	t.Run("authenticated", func(t *testing.T) {
		app := newTestApp(&userUseCaseMock{
			loginFn: func(_ context.Context, email, _ string) (entity.User, error) {
				return entity.User{ID: "landlord-1", Email: email, Role: entity.RoleLandlord}, nil
			},
		}, &propertyUseCaseMock{
			getPropertyFn: func(_ context.Context, id, landlordID string) (entity.Property, error) {
				return entity.Property{ID: id, LandlordID: landlordID, Name: "Sunny Apartment"}, nil
			},
		})

		cookie := loginAndGetCookie(t, app)

		resp := doRequest(t, app, http.MethodGet, "/v1/properties/prop-1", nil, cookie)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusOK)
		}

		var body entity.Property
		decodeJSON(t, resp, &body)

		if body.ID != "prop-1" || body.LandlordID != "landlord-1" {
			t.Fatalf("body = %+v, want property prop-1 for landlord-1", body)
		}
	})

	t.Run("not found", func(t *testing.T) {
		app := newTestApp(&userUseCaseMock{
			loginFn: func(_ context.Context, email, _ string) (entity.User, error) {
				return entity.User{ID: "landlord-1", Email: email, Role: entity.RoleLandlord}, nil
			},
		}, &propertyUseCaseMock{
			getPropertyFn: func(context.Context, string, string) (entity.Property, error) {
				return entity.Property{}, repo.ErrPropertyNotFound
			},
		})

		cookie := loginAndGetCookie(t, app)

		resp := doRequest(t, app, http.MethodGet, "/v1/properties/prop-1", nil, cookie)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusNotFound)
		}
	})

	t.Run("internal error", func(t *testing.T) {
		app := newTestApp(&userUseCaseMock{
			loginFn: func(_ context.Context, email, _ string) (entity.User, error) {
				return entity.User{ID: "landlord-1", Email: email, Role: entity.RoleLandlord}, nil
			},
		}, &propertyUseCaseMock{
			getPropertyFn: func(context.Context, string, string) (entity.Property, error) {
				return entity.Property{}, errors.New("boom")
			},
		})

		cookie := loginAndGetCookie(t, app)

		resp := doRequest(t, app, http.MethodGet, "/v1/properties/prop-1", nil, cookie)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusInternalServerError {
			t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusInternalServerError)
		}
	})
}

func TestUpdateProperty(t *testing.T) {
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

	t.Run("unauthenticated", func(t *testing.T) {
		app := newTestApp(&userUseCaseMock{}, &propertyUseCaseMock{})

		resp := doRequest(t, app, http.MethodPut, "/v1/properties/prop-1", validBody)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnauthorized {
			t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusUnauthorized)
		}
	})

	t.Run("forbidden for non-landlord", func(t *testing.T) {
		app := newTestApp(&userUseCaseMock{
			loginFn: func(_ context.Context, email, _ string) (entity.User, error) {
				return entity.User{ID: "tenant-1", Email: email, Role: entity.RoleTenant}, nil
			},
		}, &propertyUseCaseMock{})

		cookie := loginAndGetCookie(t, app)

		resp := doRequest(t, app, http.MethodPut, "/v1/properties/prop-1", validBody, cookie)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusForbidden {
			t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusForbidden)
		}
	})

	t.Run("invalid body", func(t *testing.T) {
		app := newTestApp(&userUseCaseMock{
			loginFn: func(_ context.Context, email, _ string) (entity.User, error) {
				return entity.User{ID: "landlord-1", Email: email, Role: entity.RoleLandlord}, nil
			},
		}, &propertyUseCaseMock{})

		cookie := loginAndGetCookie(t, app)

		resp := doRequest(t, app, http.MethodPut, "/v1/properties/prop-1", map[string]string{"name": "missing required fields"}, cookie)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusBadRequest)
		}
	})

	t.Run("success", func(t *testing.T) {
		app := newTestApp(&userUseCaseMock{
			loginFn: func(_ context.Context, email, _ string) (entity.User, error) {
				return entity.User{ID: "landlord-1", Email: email, Role: entity.RoleLandlord}, nil
			},
		}, &propertyUseCaseMock{
			updatePropertyFn: func(context.Context, string, string, request.Property) error {
				return nil
			},
		})

		cookie := loginAndGetCookie(t, app)

		resp := doRequest(t, app, http.MethodPut, "/v1/properties/prop-1", validBody, cookie)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusOK)
		}

		var body map[string]bool
		decodeJSON(t, resp, &body)

		if !body["ok"] {
			t.Fatalf("body = %+v, want ok=true", body)
		}
	})

	t.Run("not found", func(t *testing.T) {
		app := newTestApp(&userUseCaseMock{
			loginFn: func(_ context.Context, email, _ string) (entity.User, error) {
				return entity.User{ID: "landlord-1", Email: email, Role: entity.RoleLandlord}, nil
			},
		}, &propertyUseCaseMock{
			updatePropertyFn: func(context.Context, string, string, request.Property) error {
				return repo.ErrPropertyNotFound
			},
		})

		cookie := loginAndGetCookie(t, app)

		resp := doRequest(t, app, http.MethodPut, "/v1/properties/prop-1", validBody, cookie)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusNotFound)
		}
	})

	t.Run("internal error", func(t *testing.T) {
		app := newTestApp(&userUseCaseMock{
			loginFn: func(_ context.Context, email, _ string) (entity.User, error) {
				return entity.User{ID: "landlord-1", Email: email, Role: entity.RoleLandlord}, nil
			},
		}, &propertyUseCaseMock{
			updatePropertyFn: func(context.Context, string, string, request.Property) error {
				return errors.New("boom")
			},
		})

		cookie := loginAndGetCookie(t, app)

		resp := doRequest(t, app, http.MethodPut, "/v1/properties/prop-1", validBody, cookie)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusInternalServerError {
			t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusInternalServerError)
		}
	})
}

func TestDeleteProperty(t *testing.T) {
	t.Run("unauthenticated", func(t *testing.T) {
		app := newTestApp(&userUseCaseMock{}, &propertyUseCaseMock{})

		resp := doRequest(t, app, http.MethodDelete, "/v1/properties/prop-1", nil)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnauthorized {
			t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusUnauthorized)
		}
	})

	t.Run("forbidden for non-landlord", func(t *testing.T) {
		app := newTestApp(&userUseCaseMock{
			loginFn: func(_ context.Context, email, _ string) (entity.User, error) {
				return entity.User{ID: "tenant-1", Email: email, Role: entity.RoleTenant}, nil
			},
		}, &propertyUseCaseMock{})

		cookie := loginAndGetCookie(t, app)

		resp := doRequest(t, app, http.MethodDelete, "/v1/properties/prop-1", nil, cookie)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusForbidden {
			t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusForbidden)
		}
	})

	t.Run("success", func(t *testing.T) {
		app := newTestApp(&userUseCaseMock{
			loginFn: func(_ context.Context, email, _ string) (entity.User, error) {
				return entity.User{ID: "landlord-1", Email: email, Role: entity.RoleLandlord}, nil
			},
		}, &propertyUseCaseMock{
			deletePropertyFn: func(context.Context, string, string) error {
				return nil
			},
		})

		cookie := loginAndGetCookie(t, app)

		resp := doRequest(t, app, http.MethodDelete, "/v1/properties/prop-1", nil, cookie)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusOK)
		}

		var body map[string]bool
		decodeJSON(t, resp, &body)

		if !body["ok"] {
			t.Fatalf("body = %+v, want ok=true", body)
		}
	})

	t.Run("not found", func(t *testing.T) {
		app := newTestApp(&userUseCaseMock{
			loginFn: func(_ context.Context, email, _ string) (entity.User, error) {
				return entity.User{ID: "landlord-1", Email: email, Role: entity.RoleLandlord}, nil
			},
		}, &propertyUseCaseMock{
			deletePropertyFn: func(context.Context, string, string) error {
				return repo.ErrPropertyNotFound
			},
		})

		cookie := loginAndGetCookie(t, app)

		resp := doRequest(t, app, http.MethodDelete, "/v1/properties/prop-1", nil, cookie)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusNotFound)
		}
	})

	t.Run("internal error", func(t *testing.T) {
		app := newTestApp(&userUseCaseMock{
			loginFn: func(_ context.Context, email, _ string) (entity.User, error) {
				return entity.User{ID: "landlord-1", Email: email, Role: entity.RoleLandlord}, nil
			},
		}, &propertyUseCaseMock{
			deletePropertyFn: func(context.Context, string, string) error {
				return errors.New("boom")
			},
		})

		cookie := loginAndGetCookie(t, app)

		resp := doRequest(t, app, http.MethodDelete, "/v1/properties/prop-1", nil, cookie)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusInternalServerError {
			t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusInternalServerError)
		}
	})
}

func loginAndGetCookie(t *testing.T, app *fiber.App) *http.Cookie {
	t.Helper()

	loginResp := doRequest(t, app, http.MethodPost, "/v1/auth/login", request.Login{Email: "john@example.com", Password: "password123"})
	defer loginResp.Body.Close()

	cookie := sessionCookie(loginResp)
	if cookie == nil {
		t.Fatal("expected session_id cookie from login")
	}

	return cookie
}
