package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/example/ultimate-ci-cd-pipeline/internal/handler"
	"github.com/example/ultimate-ci-cd-pipeline/internal/model"
)

func TestHealthHandler_Healthz(t *testing.T) {
	h := handler.NewHealthHandler("1.0.0")

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rr := httptest.NewRecorder()

	h.Healthz(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	var resp model.HealthResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Status != "UP" || resp.Version != "1.0.0" {
		t.Errorf("unexpected response content: %+v", resp)
	}
}

func TestHealthHandler_Probes(t *testing.T) {
	h := handler.NewHealthHandler("1.0.0")

	t.Run("Livez", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/livez", nil)
		rr := httptest.NewRecorder()
		h.Livez(rr, req)
		if rr.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rr.Code)
		}
	})

	t.Run("Readyz", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
		rr := httptest.NewRecorder()
		h.Readyz(rr, req)
		if rr.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rr.Code)
		}
	})
}
