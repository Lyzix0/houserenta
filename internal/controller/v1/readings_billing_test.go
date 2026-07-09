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

func TestCreateReading(t *testing.T) {
	validBody := request.Reading{Gvs: 15.2, Hvs: 29.8, El1: 395.5}

	t.Run("unauthenticated", func(t *testing.T) {
		app := newTestApp(&userUseCaseMock{}, &propertyUseCaseMock{})

		resp := doRequest(t, app, http.MethodPost, "/v1/properties/prop-1/readings", validBody)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnauthorized {
			t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusUnauthorized)
		}
	})

	t.Run("invalid body", func(t *testing.T) {
		app := newTestApp(&userUseCaseMock{
			loginFn: func(_ context.Context, email, _ string) (entity.User, error) {
				return entity.User{ID: "tenant-1", Email: email, Role: entity.RoleTenant}, nil
			},
		}, &propertyUseCaseMock{})

		cookie := loginAndGetCookie(t, app)

		resp := doRequest(t, app, http.MethodPost, "/v1/properties/prop-1/readings", map[string]string{"gvs": ""}, cookie)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusBadRequest)
		}
	})

	t.Run("value too large is rejected before reaching the usecase", func(t *testing.T) {
		app := newTestApp(&userUseCaseMock{
			loginFn: func(_ context.Context, email, _ string) (entity.User, error) {
				return entity.User{ID: "tenant-1", Email: email, Role: entity.RoleTenant}, nil
			},
		}, &propertyUseCaseMock{
			createReadingFn: func(context.Context, string, string, entity.Role, request.Reading) error {
				t.Fatal("usecase should not be called: oversized reading must be rejected by validation")
				return nil
			},
		})

		cookie := loginAndGetCookie(t, app)

		oversized := validBody
		oversized.Gvs = 1e9

		resp := doRequest(t, app, http.MethodPost, "/v1/properties/prop-1/readings", oversized, cookie)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusBadRequest)
		}
	})

	t.Run("success as tenant", func(t *testing.T) {
		app := newTestApp(&userUseCaseMock{
			loginFn: func(_ context.Context, email, _ string) (entity.User, error) {
				return entity.User{ID: "tenant-1", Email: email, Role: entity.RoleTenant}, nil
			},
		}, &propertyUseCaseMock{
			createReadingFn: func(_ context.Context, propertyID, userID string, role entity.Role, _ request.Reading) error {
				if role != entity.RoleTenant {
					t.Fatalf("role = %q, want %q", role, entity.RoleTenant)
				}
				return nil
			},
		})

		cookie := loginAndGetCookie(t, app)

		resp := doRequest(t, app, http.MethodPost, "/v1/properties/prop-1/readings", validBody, cookie)
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
			createReadingFn: func(context.Context, string, string, entity.Role, request.Reading) error {
				return repo.ErrPropertyNotFound
			},
		})

		cookie := loginAndGetCookie(t, app)

		resp := doRequest(t, app, http.MethodPost, "/v1/properties/prop-1/readings", validBody, cookie)
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
			createReadingFn: func(context.Context, string, string, entity.Role, request.Reading) error {
				return errors.New("boom")
			},
		})

		cookie := loginAndGetCookie(t, app)

		resp := doRequest(t, app, http.MethodPost, "/v1/properties/prop-1/readings", validBody, cookie)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusInternalServerError {
			t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusInternalServerError)
		}
	})
}

func TestPay(t *testing.T) {
	validBody := request.Payment{Amount: 35000}

	t.Run("unauthenticated", func(t *testing.T) {
		app := newTestApp(&userUseCaseMock{}, &propertyUseCaseMock{})

		resp := doRequest(t, app, http.MethodPost, "/v1/properties/prop-1/pay", validBody)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnauthorized {
			t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusUnauthorized)
		}
	})

	t.Run("invalid body", func(t *testing.T) {
		app := newTestApp(&userUseCaseMock{
			loginFn: func(_ context.Context, email, _ string) (entity.User, error) {
				return entity.User{ID: "tenant-1", Email: email, Role: entity.RoleTenant}, nil
			},
		}, &propertyUseCaseMock{})

		cookie := loginAndGetCookie(t, app)

		resp := doRequest(t, app, http.MethodPost, "/v1/properties/prop-1/pay", map[string]any{"amount": 0}, cookie)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusBadRequest)
		}
	})

	t.Run("amount too large is rejected before reaching the usecase", func(t *testing.T) {
		app := newTestApp(&userUseCaseMock{
			loginFn: func(_ context.Context, email, _ string) (entity.User, error) {
				return entity.User{ID: "tenant-1", Email: email, Role: entity.RoleTenant}, nil
			},
		}, &propertyUseCaseMock{
			payFn: func(context.Context, string, string, entity.Role, request.Payment) error {
				t.Fatal("usecase should not be called: oversized amount must be rejected by validation")
				return nil
			},
		})

		cookie := loginAndGetCookie(t, app)

		resp := doRequest(t, app, http.MethodPost, "/v1/properties/prop-1/pay", request.Payment{Amount: 1e9}, cookie)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusBadRequest)
		}
	})

	t.Run("success top-up balance", func(t *testing.T) {
		app := newTestApp(&userUseCaseMock{
			loginFn: func(_ context.Context, email, _ string) (entity.User, error) {
				return entity.User{ID: "tenant-1", Email: email, Role: entity.RoleTenant}, nil
			},
		}, &propertyUseCaseMock{
			payFn: func(_ context.Context, propertyID, userID string, role entity.Role, body request.Payment) error {
				if body.BillID != nil {
					t.Fatalf("billId = %v, want nil", body.BillID)
				}
				return nil
			},
		})

		cookie := loginAndGetCookie(t, app)

		resp := doRequest(t, app, http.MethodPost, "/v1/properties/prop-1/pay", validBody, cookie)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusOK)
		}
	})

	t.Run("success pay specific bill", func(t *testing.T) {
		billID := "bill-999"
		app := newTestApp(&userUseCaseMock{
			loginFn: func(_ context.Context, email, _ string) (entity.User, error) {
				return entity.User{ID: "tenant-1", Email: email, Role: entity.RoleTenant}, nil
			},
		}, &propertyUseCaseMock{
			payFn: func(_ context.Context, propertyID, userID string, role entity.Role, body request.Payment) error {
				if body.BillID == nil || *body.BillID != billID {
					t.Fatalf("billId = %v, want %q", body.BillID, billID)
				}
				return nil
			},
		})

		cookie := loginAndGetCookie(t, app)

		resp := doRequest(t, app, http.MethodPost, "/v1/properties/prop-1/pay", request.Payment{Amount: 35000, BillID: &billID}, cookie)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusOK)
		}
	})

	t.Run("bill not found", func(t *testing.T) {
		app := newTestApp(&userUseCaseMock{
			loginFn: func(_ context.Context, email, _ string) (entity.User, error) {
				return entity.User{ID: "tenant-1", Email: email, Role: entity.RoleTenant}, nil
			},
		}, &propertyUseCaseMock{
			payFn: func(context.Context, string, string, entity.Role, request.Payment) error {
				return repo.ErrBillNotFound
			},
		})

		cookie := loginAndGetCookie(t, app)

		resp := doRequest(t, app, http.MethodPost, "/v1/properties/prop-1/pay", validBody, cookie)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusNotFound)
		}
	})

	t.Run("property not found", func(t *testing.T) {
		app := newTestApp(&userUseCaseMock{
			loginFn: func(_ context.Context, email, _ string) (entity.User, error) {
				return entity.User{ID: "tenant-1", Email: email, Role: entity.RoleTenant}, nil
			},
		}, &propertyUseCaseMock{
			payFn: func(context.Context, string, string, entity.Role, request.Payment) error {
				return repo.ErrPropertyNotFound
			},
		})

		cookie := loginAndGetCookie(t, app)

		resp := doRequest(t, app, http.MethodPost, "/v1/properties/prop-1/pay", validBody, cookie)
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
			payFn: func(context.Context, string, string, entity.Role, request.Payment) error {
				return errors.New("boom")
			},
		})

		cookie := loginAndGetCookie(t, app)

		resp := doRequest(t, app, http.MethodPost, "/v1/properties/prop-1/pay", validBody, cookie)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusInternalServerError {
			t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusInternalServerError)
		}
	})
}

func TestCreateCustomItem(t *testing.T) {
	validBody := request.CustomItem{Description: "Замена смесителя на кухне", Amount: 2500}

	t.Run("unauthenticated", func(t *testing.T) {
		app := newTestApp(&userUseCaseMock{}, &propertyUseCaseMock{})

		resp := doRequest(t, app, http.MethodPost, "/v1/properties/prop-1/custom-item", validBody)
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

		resp := doRequest(t, app, http.MethodPost, "/v1/properties/prop-1/custom-item", validBody, cookie)
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

		resp := doRequest(t, app, http.MethodPost, "/v1/properties/prop-1/custom-item", map[string]string{"description": ""}, cookie)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusBadRequest)
		}
	})

	t.Run("amount too large is rejected before reaching the usecase", func(t *testing.T) {
		app := newTestApp(&userUseCaseMock{
			loginFn: func(_ context.Context, email, _ string) (entity.User, error) {
				return entity.User{ID: "landlord-1", Email: email, Role: entity.RoleLandlord}, nil
			},
		}, &propertyUseCaseMock{
			createCustomItemFn: func(context.Context, string, string, request.CustomItem) error {
				t.Fatal("usecase should not be called: oversized amount must be rejected by validation")
				return nil
			},
		})

		cookie := loginAndGetCookie(t, app)

		oversized := validBody
		oversized.Amount = 1e9

		resp := doRequest(t, app, http.MethodPost, "/v1/properties/prop-1/custom-item", oversized, cookie)
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
			createCustomItemFn: func(context.Context, string, string, request.CustomItem) error {
				return nil
			},
		})

		cookie := loginAndGetCookie(t, app)

		resp := doRequest(t, app, http.MethodPost, "/v1/properties/prop-1/custom-item", validBody, cookie)
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
			createCustomItemFn: func(context.Context, string, string, request.CustomItem) error {
				return repo.ErrPropertyNotFound
			},
		})

		cookie := loginAndGetCookie(t, app)

		resp := doRequest(t, app, http.MethodPost, "/v1/properties/prop-1/custom-item", validBody, cookie)
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
			createCustomItemFn: func(context.Context, string, string, request.CustomItem) error {
				return errors.New("boom")
			},
		})

		cookie := loginAndGetCookie(t, app)

		resp := doRequest(t, app, http.MethodPost, "/v1/properties/prop-1/custom-item", validBody, cookie)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusInternalServerError {
			t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusInternalServerError)
		}
	})
}
