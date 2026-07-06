package v1_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"

	v1 "github.com/potom_pridumaem/internal/controller/v1"
	"github.com/potom_pridumaem/internal/usecase"
)

func newTestApp(u usecase.User, p usecase.Property) *fiber.App {
	app := fiber.New()
	apiV1Group := app.Group("/v1")
	v1.NewRoutes(apiV1Group, u, p, zap.NewNop())

	return app
}

func doRequest(t *testing.T, app *fiber.App, method, path string, body any, cookies ...*http.Cookie) *http.Response {
	t.Helper()

	var reader *bytes.Reader
	if body != nil {
		raw, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("marshal body: %v", err)
		}
		reader = bytes.NewReader(raw)
	} else {
		reader = bytes.NewReader(nil)
	}

	req, err := http.NewRequest(method, path, reader)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	for _, c := range cookies {
		req.AddCookie(c)
	}

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}

	return resp
}

func decodeJSON(t *testing.T, resp *http.Response, v any) {
	t.Helper()

	if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
		t.Fatalf("decode json: %v", err)
	}
}

func sessionCookie(resp *http.Response) *http.Cookie {
	for _, c := range resp.Cookies() {
		if c.Name == "session_id" {
			return c
		}
	}
	return nil
}
