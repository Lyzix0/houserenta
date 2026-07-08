package v1_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	entity "github.com/potom_pridumaem/internal/entity/users"
	"github.com/potom_pridumaem/internal/repo"
)

func TestGetVacantProperties(t *testing.T) {
	t.Run("unauthenticated", func(t *testing.T) {
		app := newTestApp(&userUseCaseMock{}, &propertyUseCaseMock{})

		resp := doRequest(t, app, http.MethodGet, "/v1/properties/vacant", nil)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnauthorized {
			t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusUnauthorized)
		}
	})

	t.Run("success", func(t *testing.T) {
		app := newTestApp(&userUseCaseMock{
			loginFn: func(_ context.Context, email, _ string) (entity.User, error) {
				return entity.User{ID: "tenant-1", Email: email, Role: entity.RoleTenant}, nil
			},
		}, &propertyUseCaseMock{
			getVacantPropertiesFn: func(context.Context) ([]entity.PropertyDetail, error) {
				return []entity.PropertyDetail{
					{
						Property:     entity.Property{ID: "prop-1", Name: "Vacant Studio"},
						Applications: []entity.Application{{ID: "app-1", PropertyID: "prop-1", TenantUserID: "tenant-2"}},
					},
				}, nil
			},
		})

		cookie := loginAndGetCookie(t, app)

		resp := doRequest(t, app, http.MethodGet, "/v1/properties/vacant", nil, cookie)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusOK)
		}

		var body []entity.PropertyDetail
		decodeJSON(t, resp, &body)

		if len(body) != 1 || body[0].ID != "prop-1" {
			t.Fatalf("body = %+v, want a single vacant property prop-1", body)
		}
		if len(body[0].Applications) != 1 {
			t.Fatalf("applications = %+v, want 1 entry", body[0].Applications)
		}
	})

	t.Run("internal error", func(t *testing.T) {
		app := newTestApp(&userUseCaseMock{
			loginFn: func(_ context.Context, email, _ string) (entity.User, error) {
				return entity.User{ID: "tenant-1", Email: email, Role: entity.RoleTenant}, nil
			},
		}, &propertyUseCaseMock{
			getVacantPropertiesFn: func(context.Context) ([]entity.PropertyDetail, error) {
				return nil, errors.New("boom")
			},
		})

		cookie := loginAndGetCookie(t, app)

		resp := doRequest(t, app, http.MethodGet, "/v1/properties/vacant", nil, cookie)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusInternalServerError {
			t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusInternalServerError)
		}
	})
}

func TestApplyToProperty(t *testing.T) {
	t.Run("unauthenticated", func(t *testing.T) {
		app := newTestApp(&userUseCaseMock{}, &propertyUseCaseMock{})

		resp := doRequest(t, app, http.MethodPost, "/v1/properties/prop-1/apply", nil)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnauthorized {
			t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusUnauthorized)
		}
	})

	t.Run("forbidden for non-tenant", func(t *testing.T) {
		app := newTestApp(&userUseCaseMock{
			loginFn: func(_ context.Context, email, _ string) (entity.User, error) {
				return entity.User{ID: "landlord-1", Email: email, Role: entity.RoleLandlord}, nil
			},
		}, &propertyUseCaseMock{})

		cookie := loginAndGetCookie(t, app)

		resp := doRequest(t, app, http.MethodPost, "/v1/properties/prop-1/apply", nil, cookie)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusForbidden {
			t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusForbidden)
		}
	})

	t.Run("success", func(t *testing.T) {
		app := newTestApp(&userUseCaseMock{
			loginFn: func(_ context.Context, email, _ string) (entity.User, error) {
				return entity.User{ID: "tenant-1", Email: email, Role: entity.RoleTenant}, nil
			},
		}, &propertyUseCaseMock{
			applyFn: func(_ context.Context, propertyID, tenantUserID string) error {
				if propertyID != "prop-1" || tenantUserID != "tenant-1" {
					t.Fatalf("Apply(%q, %q), want (prop-1, tenant-1)", propertyID, tenantUserID)
				}
				return nil
			},
		})

		cookie := loginAndGetCookie(t, app)

		resp := doRequest(t, app, http.MethodPost, "/v1/properties/prop-1/apply", nil, cookie)
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

	t.Run("already applied", func(t *testing.T) {
		app := newTestApp(&userUseCaseMock{
			loginFn: func(_ context.Context, email, _ string) (entity.User, error) {
				return entity.User{ID: "tenant-1", Email: email, Role: entity.RoleTenant}, nil
			},
		}, &propertyUseCaseMock{
			applyFn: func(context.Context, string, string) error {
				return repo.ErrApplicationAlreadyExists
			},
		})

		cookie := loginAndGetCookie(t, app)

		resp := doRequest(t, app, http.MethodPost, "/v1/properties/prop-1/apply", nil, cookie)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusBadRequest)
		}
	})

	t.Run("property not found", func(t *testing.T) {
		app := newTestApp(&userUseCaseMock{
			loginFn: func(_ context.Context, email, _ string) (entity.User, error) {
				return entity.User{ID: "tenant-1", Email: email, Role: entity.RoleTenant}, nil
			},
		}, &propertyUseCaseMock{
			applyFn: func(context.Context, string, string) error {
				return repo.ErrPropertyNotFound
			},
		})

		cookie := loginAndGetCookie(t, app)

		resp := doRequest(t, app, http.MethodPost, "/v1/properties/prop-1/apply", nil, cookie)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusNotFound)
		}
	})

	t.Run("internal error", func(t *testing.T) {
		app := newTestApp(&userUseCaseMock{
			loginFn: func(_ context.Context, email, _ string) (entity.User, error) {
				return entity.User{ID: "tenant-1", Email: email, Role: entity.RoleTenant}, nil
			},
		}, &propertyUseCaseMock{
			applyFn: func(context.Context, string, string) error {
				return errors.New("boom")
			},
		})

		cookie := loginAndGetCookie(t, app)

		resp := doRequest(t, app, http.MethodPost, "/v1/properties/prop-1/apply", nil, cookie)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusInternalServerError {
			t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusInternalServerError)
		}
	})
}
