package app

import (
	"code-review-agent/backend/internal/config"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCorsAllowsLocalhostOrigins(t *testing.T) {
	app := New(config.Config{CORSOrigins: []string{"http://localhost:3002"}}, nil)

	req := httptest.NewRequest(http.MethodOptions, "/api/auth/login", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	res := httptest.NewRecorder()

	app.Routes().ServeHTTP(res, req)

	if got := res.Header().Get("Access-Control-Allow-Origin"); got != "http://localhost:3000" {
		t.Fatalf("expected localhost:3000 to be allowed, got %q", got)
	}
}
