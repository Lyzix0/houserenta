package v1_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/potom_pridumaem/internal/controller/v1/response"
	entity "github.com/potom_pridumaem/internal/entity/users"
)

func TestGetUnlinkedTenants(t *testing.T) {
	t.Run("unauthenticated", func(t *testing.T) {
		app := newTestApp(&userUseCaseMock{}, &propertyUseCaseMock{})

		resp := doRequest(t, app, http.MethodGet, "/v1/tenants/unlinked", nil)
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

		resp := doRequest(t, app, http.MethodGet, "/v1/tenants/unlinked", nil, cookie)
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
			getUnlinkedTenants: func(context.Context) ([]entity.User, error) {
				return []entity.User{
					{ID: "tenant-2", Name: "Иванов Иван", Email: "ivanov@example.com", Document: "1234", Phone: "+79991112233"},
				}, nil
			},
		}, &propertyUseCaseMock{})

		cookie := loginAndGetCookie(t, app)

		resp := doRequest(t, app, http.MethodGet, "/v1/tenants/unlinked", nil, cookie)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusOK)
		}

		var body []response.TenantSummary
		decodeJSON(t, resp, &body)

		if len(body) != 1 || body[0].ID != "tenant-2" {
			t.Fatalf("body = %+v, want a single tenant-2 summary", body)
		}
	})

	t.Run("internal error", func(t *testing.T) {
		app := newTestApp(&userUseCaseMock{
			loginFn: func(_ context.Context, email, _ string) (entity.User, error) {
				return entity.User{ID: "landlord-1", Email: email, Role: entity.RoleLandlord}, nil
			},
			getUnlinkedTenants: func(context.Context) ([]entity.User, error) {
				return nil, errors.New("boom")
			},
		}, &propertyUseCaseMock{})

		cookie := loginAndGetCookie(t, app)

		resp := doRequest(t, app, http.MethodGet, "/v1/tenants/unlinked", nil, cookie)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusInternalServerError {
			t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusInternalServerError)
		}
	})
}
