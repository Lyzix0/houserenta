package v1_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/potom_pridumaem/internal/controller/v1/request"
	entity "github.com/potom_pridumaem/internal/entity/users"
	"github.com/potom_pridumaem/internal/repo"
	"github.com/potom_pridumaem/internal/usecase"
)

func TestRegister(t *testing.T) {
	validBody := request.Register{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "password123",
		Role:     "tenant",
		Document: "1234567890",
		Phone:    "+79991234567",
	}

	tests := []struct {
		name       string
		body       any
		registerFn func(ctx context.Context, name, email, password, role, document, phone string) (entity.User, error)
		wantStatus int
	}{
		{
			name: "success",
			body: validBody,
			registerFn: func(_ context.Context, name, email, _, role, document, phone string) (entity.User, error) {
				return entity.User{ID: "1", Name: name, Email: email, Role: entity.Role(role), Document: document, Phone: phone}, nil
			},
			wantStatus: http.StatusCreated,
		},
		{
			name:       "invalid body",
			body:       map[string]string{"email": "not-an-email"},
			registerFn: func(context.Context, string, string, string, string, string, string) (entity.User, error) { return entity.User{}, nil },
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "invalid role",
			body: validBody,
			registerFn: func(context.Context, string, string, string, string, string, string) (entity.User, error) {
				return entity.User{}, entity.ErrInvalidRole
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "email already taken",
			body: validBody,
			registerFn: func(context.Context, string, string, string, string, string, string) (entity.User, error) {
				return entity.User{}, repo.ErrEmailAlreadyTaken
			},
			wantStatus: http.StatusConflict,
		},
		{
			name: "internal error",
			body: validBody,
			registerFn: func(context.Context, string, string, string, string, string, string) (entity.User, error) {
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
		})
	}
}

func TestLogin(t *testing.T) {
	validBody := request.Login{Email: "john@example.com", Password: "password123"}

	tests := []struct {
		name       string
		body       any
		loginFn    func(ctx context.Context, email, password string) (entity.User, error)
		wantStatus int
	}{
		{
			name: "success",
			body: validBody,
			loginFn: func(_ context.Context, email, _ string) (entity.User, error) {
				return entity.User{ID: "1", Email: email, Role: entity.RoleTenant}, nil
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "invalid body",
			body:       map[string]string{"email": "not-an-email"},
			loginFn:    func(context.Context, string, string) (entity.User, error) { return entity.User{}, nil },
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "invalid credentials",
			body: validBody,
			loginFn: func(context.Context, string, string) (entity.User, error) {
				return entity.User{}, usecase.ErrInvalidCredentials
			},
			wantStatus: http.StatusUnauthorized,
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

	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusNoContent)
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

	t.Run("authenticated", func(t *testing.T) {
		app := newTestApp(&userUseCaseMock{
			loginFn: func(_ context.Context, email, _ string) (entity.User, error) {
				return entity.User{ID: "user-1", Email: email, Role: entity.RoleLandlord}, nil
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

		var body map[string]string
		decodeJSON(t, resp, &body)

		if body["id"] != "user-1" {
			t.Fatalf("id = %q, want %q", body["id"], "user-1")
		}
		if body["role"] != string(entity.RoleLandlord) {
			t.Fatalf("role = %q, want %q", body["role"], entity.RoleLandlord)
		}
	})
}
