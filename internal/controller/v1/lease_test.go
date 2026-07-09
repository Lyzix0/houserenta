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

func TestCreateLease(t *testing.T) {
	validBody := request.Lease{
		TenantUserID: "tenant-1",
		Price:        32000,
		MonthsOfRent: 11,
		PaymentDay:   10,
		ReadingDay:   28,
	}

	t.Run("unauthenticated", func(t *testing.T) {
		app := newTestApp(&userUseCaseMock{}, &propertyUseCaseMock{})

		resp := doRequest(t, app, http.MethodPost, "/v1/properties/prop-1/lease", validBody)
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

		resp := doRequest(t, app, http.MethodPost, "/v1/properties/prop-1/lease", validBody, cookie)
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

		resp := doRequest(t, app, http.MethodPost, "/v1/properties/prop-1/lease", map[string]string{"tenantUserId": ""}, cookie)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusBadRequest)
		}
	})

	t.Run("price too large is rejected before reaching the usecase", func(t *testing.T) {
		app := newTestApp(&userUseCaseMock{
			loginFn: func(_ context.Context, email, _ string) (entity.User, error) {
				return entity.User{ID: "landlord-1", Email: email, Role: entity.RoleLandlord}, nil
			},
		}, &propertyUseCaseMock{
			createLeaseFn: func(context.Context, string, string, request.Lease) error {
				t.Fatal("usecase should not be called: oversized price must be rejected by validation")
				return nil
			},
		})

		cookie := loginAndGetCookie(t, app)

		oversized := validBody
		oversized.Price = 1e9

		resp := doRequest(t, app, http.MethodPost, "/v1/properties/prop-1/lease", oversized, cookie)
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
			createLeaseFn: func(context.Context, string, string, request.Lease) error {
				return nil
			},
		})

		cookie := loginAndGetCookie(t, app)

		resp := doRequest(t, app, http.MethodPost, "/v1/properties/prop-1/lease", validBody, cookie)
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

	t.Run("property not found", func(t *testing.T) {
		app := newTestApp(&userUseCaseMock{
			loginFn: func(_ context.Context, email, _ string) (entity.User, error) {
				return entity.User{ID: "landlord-1", Email: email, Role: entity.RoleLandlord}, nil
			},
		}, &propertyUseCaseMock{
			createLeaseFn: func(context.Context, string, string, request.Lease) error {
				return repo.ErrPropertyNotFound
			},
		})

		cookie := loginAndGetCookie(t, app)

		resp := doRequest(t, app, http.MethodPost, "/v1/properties/prop-1/lease", validBody, cookie)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusNotFound)
		}
	})

	t.Run("tenant not found", func(t *testing.T) {
		app := newTestApp(&userUseCaseMock{
			loginFn: func(_ context.Context, email, _ string) (entity.User, error) {
				return entity.User{ID: "landlord-1", Email: email, Role: entity.RoleLandlord}, nil
			},
		}, &propertyUseCaseMock{
			createLeaseFn: func(context.Context, string, string, request.Lease) error {
				return repo.ErrTenantNotFound
			},
		})

		cookie := loginAndGetCookie(t, app)

		resp := doRequest(t, app, http.MethodPost, "/v1/properties/prop-1/lease", validBody, cookie)
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
			createLeaseFn: func(context.Context, string, string, request.Lease) error {
				return errors.New("boom")
			},
		})

		cookie := loginAndGetCookie(t, app)

		resp := doRequest(t, app, http.MethodPost, "/v1/properties/prop-1/lease", validBody, cookie)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusInternalServerError {
			t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusInternalServerError)
		}
	})
}

func TestDeleteLease(t *testing.T) {
	t.Run("unauthenticated", func(t *testing.T) {
		app := newTestApp(&userUseCaseMock{}, &propertyUseCaseMock{})

		resp := doRequest(t, app, http.MethodDelete, "/v1/properties/prop-1/lease", nil)
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

		resp := doRequest(t, app, http.MethodDelete, "/v1/properties/prop-1/lease", nil, cookie)
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
			deleteLeaseFn: func(context.Context, string, string) error {
				return nil
			},
		})

		cookie := loginAndGetCookie(t, app)

		resp := doRequest(t, app, http.MethodDelete, "/v1/properties/prop-1/lease", nil, cookie)
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

	t.Run("property not found", func(t *testing.T) {
		app := newTestApp(&userUseCaseMock{
			loginFn: func(_ context.Context, email, _ string) (entity.User, error) {
				return entity.User{ID: "landlord-1", Email: email, Role: entity.RoleLandlord}, nil
			},
		}, &propertyUseCaseMock{
			deleteLeaseFn: func(context.Context, string, string) error {
				return repo.ErrPropertyNotFound
			},
		})

		cookie := loginAndGetCookie(t, app)

		resp := doRequest(t, app, http.MethodDelete, "/v1/properties/prop-1/lease", nil, cookie)
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
			deleteLeaseFn: func(context.Context, string, string) error {
				return errors.New("boom")
			},
		})

		cookie := loginAndGetCookie(t, app)

		resp := doRequest(t, app, http.MethodDelete, "/v1/properties/prop-1/lease", nil, cookie)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusInternalServerError {
			t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusInternalServerError)
		}
	})
}
