package handler_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/example/ultimate-ci-cd-pipeline/internal/handler"
)

func TestUIHandler_ServeHTTP(t *testing.T) {
	u := handler.NewUIHandler("1.0.0", "production", "order-service")

	t.Run("Root Path Serves HTML Dashboard", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rr := httptest.NewRecorder()
		u.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("expected 200 OK, got %d", rr.Code)
		}

		contentType := rr.Header().Get("Content-Type")
		if !strings.Contains(contentType, "text/html") {
			t.Errorf("expected Content-Type text/html, got %s", contentType)
		}

		body := rr.Body.String()
		if !strings.Contains(body, "Enterprise Microservice Control Center") {
			t.Errorf("expected dashboard title in HTML output")
		}
		if !strings.Contains(body, "1.0.0") {
			t.Errorf("expected version 1.0.0 in HTML output")
		}
	})

	t.Run("Invalid Subpath Returns 404 JSON", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/unknown-path", nil)
		rr := httptest.NewRecorder()
		u.ServeHTTP(rr, req)

		if rr.Code != http.StatusNotFound {
			t.Fatalf("expected 404 Not Found, got %d", rr.Code)
		}
	})
}
