package v1_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/potom_pridumaem/internal/controller/v1/request"
	"github.com/potom_pridumaem/internal/controller/v1/response"
	entity "github.com/potom_pridumaem/internal/entity/users"
	"github.com/potom_pridumaem/internal/repo"
	"github.com/potom_pridumaem/internal/usecase"
)

func TestRegister(t *testing.T) {
	validBody := request.Register{
		Name:        "John Doe",
		Email:       "john@example.com",
		Password:    "password123",
		InitialRole: "tenant",
		Document:    "1234567890",
		Phone:       "+79991234567",
	}

	tests := []struct {
		name       string
		body       any
		registerFn func(ctx context.Context, name, email, password, role, document, phone string, paymentCard *string) (entity.User, error)
		wantStatus int
	}{
		{
			name: "success",
			body: validBody,
			registerFn: func(_ context.Context, name, email, _, role, document, phone string, _ *string) (entity.User, error) {
				return entity.User{ID: "1", Name: name, Email: email, Role: entity.Role(role), Document: document, Phone: phone}, nil
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "invalid body",
			body: map[string]string{"email": "not-an-email"},
			registerFn: func(context.Context, string, string, string, string, string, string, *string) (entity.User, error) {
				return entity.User{}, nil
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "invalid role",
			body: validBody,
			registerFn: func(context.Context, string, string, string, string, string, string, *string) (entity.User, error) {
				return entity.User{}, entity.ErrInvalidRole
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "email already taken",
			body: validBody,
			registerFn: func(context.Context, string, string, string, string, string, string, *string) (entity.User, error) {
				return entity.User{}, repo.ErrEmailAlreadyTaken
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "internal error",
			body: validBody,
			registerFn: func(context.Context, string, string, string, string, string, string, *string) (entity.User, error) {
				return entity.User{}, errors.New("boom")
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := newTestApp(&userUseCaseMock{registerFn: tt.registerFn}, &propertyUseCaseMock{})

			resp := doRequest(t, app, http.MethodPost, "/v1/auth/register", tt.body)
			defer resp.Body.Close()

			if resp.StatusCode != tt.wantStatus {
				t.Fatalf("status = %d, want %d", resp.StatusCode, tt.wantStatus)
			}

			if tt.wantStatus == http.StatusOK && sessionCookie(resp) == nil {
				t.Fatal("expected session_id cookie to be set on successful register")
			}
		})
	}
}

func TestLogin(t *testing.T) {
	validBody := request.Login{Email: "john@example.com", Password: "password123"}

	tests := []struct {
		name       string
		body       any
		loginFn    func(ctx context.Context, identifier, password string) (entity.User, error)
		wantStatus int
	}{
		{
			name: "success",
			body: validBody,
			loginFn: func(_ context.Context, identifier, _ string) (entity.User, error) {
				return entity.User{ID: "1", Email: identifier, Role: entity.RoleTenant}, nil
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "invalid body",
			body:       map[string]string{"email": "john@example.com"},
			loginFn:    func(context.Context, string, string) (entity.User, error) { return entity.User{}, nil },
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "invalid credentials",
			body: validBody,
			loginFn: func(context.Context, string, string) (entity.User, error) {
				return entity.User{}, usecase.ErrInvalidCredentials
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "internal error",
			body: validBody,
			loginFn: func(context.Context, string, string) (entity.User, error) {
				return entity.User{}, errors.New("boom")
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := newTestApp(&userUseCaseMock{loginFn: tt.loginFn}, &propertyUseCaseMock{})

			resp := doRequest(t, app, http.MethodPost, "/v1/auth/login", tt.body)
			defer resp.Body.Close()

			if resp.StatusCode != tt.wantStatus {
				t.Fatalf("status = %d, want %d", resp.StatusCode, tt.wantStatus)
			}

			if tt.wantStatus == http.StatusOK && sessionCookie(resp) == nil {
				t.Fatal("expected session_id cookie to be set on successful login")
			}
		})
	}
}

func TestLogout(t *testing.T) {
	app := newTestApp(&userUseCaseMock{}, &propertyUseCaseMock{})

	resp := doRequest(t, app, http.MethodPost, "/v1/auth/logout", nil)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusOK)
	}

	var body map[string]bool
	decodeJSON(t, resp, &body)

	if !body["ok"] {
		t.Fatalf("body = %+v, want ok=true", body)
	}
}

func TestMe(t *testing.T) {
	t.Run("unauthenticated", func(t *testing.T) {
		app := newTestApp(&userUseCaseMock{}, &propertyUseCaseMock{})

		resp := doRequest(t, app, http.MethodGet, "/v1/auth/me", nil)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnauthorized {
			t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusUnauthorized)
		}
	})

	t.Run("authenticated landlord", func(t *testing.T) {
		app := newTestApp(&userUseCaseMock{
			loginFn: func(_ context.Context, email, _ string) (entity.User, error) {
				return entity.User{ID: "user-1", Email: email, Role: entity.RoleLandlord}, nil
			},
			meFn: func(_ context.Context, userID string) (usecase.UserProfile, error) {
				return usecase.UserProfile{
					User: entity.User{
						ID:       userID,
						Name:     "Ivan Petrov",
						Email:    "john@example.com",
						Role:     entity.RoleLandlord,
						Document: "1234567890",
						Phone:    "+79991234567",
					},
				}, nil
			},
		}, &propertyUseCaseMock{})

		loginResp := doRequest(t, app, http.MethodPost, "/v1/auth/login", request.Login{Email: "john@example.com", Password: "password123"})
		defer loginResp.Body.Close()

		cookie := sessionCookie(loginResp)
		if cookie == nil {
			t.Fatal("expected session_id cookie from login")
		}

		resp := doRequest(t, app, http.MethodGet, "/v1/auth/me", nil, cookie)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusOK)
		}

		var body response.Me
		decodeJSON(t, resp, &body)

		if body.ID != "user-1" {
			t.Fatalf("id = %q, want %q", body.ID, "user-1")
		}
		if body.Role != string(entity.RoleLandlord) {
			t.Fatalf("role = %q, want %q", body.Role, entity.RoleLandlord)
		}
		if body.TenantPropertyID != nil {
			t.Fatalf("tenantPropertyId = %v, want nil", *body.TenantPropertyID)
		}
	})

	t.Run("authenticated tenant with lease", func(t *testing.T) {
		app := newTestApp(&userUseCaseMock{
			loginFn: func(_ context.Context, email, _ string) (entity.User, error) {
				return entity.User{ID: "tenant-1", Email: email, Role: entity.RoleTenant}, nil
			},
			meFn: func(_ context.Context, userID string) (usecase.UserProfile, error) {
				propertyID := "prop-1"
				return usecase.UserProfile{
					User:             entity.User{ID: userID, Role: entity.RoleTenant},
					TenantPropertyID: &propertyID,
				}, nil
			},
		}, &propertyUseCaseMock{})

		loginResp := doRequest(t, app, http.MethodPost, "/v1/auth/login", request.Login{Email: "tenant@example.com", Password: "password123"})
		defer loginResp.Body.Close()

		cookie := sessionCookie(loginResp)
		if cookie == nil {
			t.Fatal("expected session_id cookie from login")
		}

		resp := doRequest(t, app, http.MethodGet, "/v1/auth/me", nil, cookie)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusOK)
		}

		var body response.Me
		decodeJSON(t, resp, &body)

		if body.TenantPropertyID == nil || *body.TenantPropertyID != "prop-1" {
			t.Fatalf("tenantPropertyId = %v, want %q", body.TenantPropertyID, "prop-1")
		}
	})

	t.Run("internal error", func(t *testing.T) {
		app := newTestApp(&userUseCaseMock{
			loginFn: func(_ context.Context, email, _ string) (entity.User, error) {
				return entity.User{ID: "user-1", Email: email, Role: entity.RoleLandlord}, nil
			},
			meFn: func(context.Context, string) (usecase.UserProfile, error) {
				return usecase.UserProfile{}, errors.New("boom")
			},
		}, &propertyUseCaseMock{})

		loginResp := doRequest(t, app, http.MethodPost, "/v1/auth/login", request.Login{Email: "john@example.com", Password: "password123"})
		defer loginResp.Body.Close()

		cookie := sessionCookie(loginResp)
		if cookie == nil {
			t.Fatal("expected session_id cookie from login")
		}

		resp := doRequest(t, app, http.MethodGet, "/v1/auth/me", nil, cookie)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusInternalServerError {
			t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusInternalServerError)
		}
	})

	t.Run("triggers billing check but survives its failure", func(t *testing.T) {
		var billingRan bool

		u := &userUseCaseMock{
			loginFn: func(_ context.Context, email, _ string) (entity.User, error) {
				return entity.User{ID: "user-1", Email: email, Role: entity.RoleLandlord}, nil
			},
			meFn: func(_ context.Context, userID string) (usecase.UserProfile, error) {
				return usecase.UserProfile{User: entity.User{ID: userID, Role: entity.RoleLandlord}}, nil
			},
		}
		app := newTestAppWithBilling(u, &propertyUseCaseMock{}, &billingUseCaseMock{
			runFn: func(context.Context) error {
				billingRan = true
				return errors.New("billing boom")
			},
		})

		loginResp := doRequest(t, app, http.MethodPost, "/v1/auth/login", request.Login{Email: "john@example.com", Password: "password123"})
		defer loginResp.Body.Close()

		cookie := sessionCookie(loginResp)
		if cookie == nil {
			t.Fatal("expected session_id cookie from login")
		}

		resp := doRequest(t, app, http.MethodGet, "/v1/auth/me", nil, cookie)
		defer resp.Body.Close()

		if !billingRan {
			t.Fatal("expected billing.Run to be called")
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("status = %d, want %d (billing failure must not break the response)", resp.StatusCode, http.StatusOK)
		}
	})
}

func TestProfile(t *testing.T) {
	name := "Ivan Kolesnikov"
	validBody := request.Profile{Name: &name}

	t.Run("unauthenticated", func(t *testing.T) {
		app := newTestApp(&userUseCaseMock{}, &propertyUseCaseMock{})

		resp := doRequest(t, app, http.MethodPost, "/v1/auth/profile", validBody)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnauthorized {
			t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusUnauthorized)
		}
	})

	t.Run("success", func(t *testing.T) {
		app := newTestApp(&userUseCaseMock{
			loginFn: func(_ context.Context, identifier, _ string) (entity.User, error) {
				return entity.User{ID: "user-1", Email: identifier, Role: entity.RoleLandlord}, nil
			},
			updateProfileFn: func(context.Context, string, request.Profile) error {
				return nil
			},
		}, &propertyUseCaseMock{})

		cookie := loginAndGetCookie(t, app)

		resp := doRequest(t, app, http.MethodPost, "/v1/auth/profile", validBody, cookie)
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

	t.Run("invalid body", func(t *testing.T) {
		app := newTestApp(&userUseCaseMock{
			loginFn: func(_ context.Context, identifier, _ string) (entity.User, error) {
				return entity.User{ID: "user-1", Email: identifier, Role: entity.RoleLandlord}, nil
			},
		}, &propertyUseCaseMock{})

		cookie := loginAndGetCookie(t, app)

		badEmail := "not-an-email"
		resp := doRequest(t, app, http.MethodPost, "/v1/auth/profile", request.Profile{Email: &badEmail}, cookie)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusBadRequest)
		}
	})

	t.Run("email already taken", func(t *testing.T) {
		app := newTestApp(&userUseCaseMock{
			loginFn: func(_ context.Context, identifier, _ string) (entity.User, error) {
				return entity.User{ID: "user-1", Email: identifier, Role: entity.RoleLandlord}, nil
			},
			updateProfileFn: func(context.Context, string, request.Profile) error {
				return repo.ErrEmailAlreadyTaken
			},
		}, &propertyUseCaseMock{})

		cookie := loginAndGetCookie(t, app)

		resp := doRequest(t, app, http.MethodPost, "/v1/auth/profile", validBody, cookie)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusBadRequest)
		}
	})

	t.Run("internal error", func(t *testing.T) {
		app := newTestApp(&userUseCaseMock{
			loginFn: func(_ context.Context, identifier, _ string) (entity.User, error) {
				return entity.User{ID: "user-1", Email: identifier, Role: entity.RoleLandlord}, nil
			},
			updateProfileFn: func(context.Context, string, request.Profile) error {
				return errors.New("boom")
			},
		}, &propertyUseCaseMock{})

		cookie := loginAndGetCookie(t, app)

		resp := doRequest(t, app, http.MethodPost, "/v1/auth/profile", validBody, cookie)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusInternalServerError {
			t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusInternalServerError)
		}
	})
}
