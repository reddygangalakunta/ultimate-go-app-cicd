package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/example/ultimate-ci-cd-pipeline/internal/handler"
	"github.com/example/ultimate-ci-cd-pipeline/internal/model"
	"github.com/example/ultimate-ci-cd-pipeline/internal/service"
)

func TestOrderHandler_InteractiveEndpoints(t *testing.T) {
	svc := service.NewInMemoryOrderService()
	h := handler.NewOrderHandler(svc)

	t.Run("Metrics Endpoint", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/metrics", nil)
		rr := httptest.NewRecorder()
		h.Metrics(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("expected 200 OK, got %d", rr.Code)
		}

		var m model.SystemMetrics
		if err := json.Unmarshal(rr.Body.Bytes(), &m); err != nil {
			t.Fatalf("failed unmarshaling metrics: %v", err)
		}
		if m.TotalOrders < 1 {
			t.Errorf("expected orders in metrics")
		}
	})

	t.Run("Update Order Status", func(t *testing.T) {
		payload := model.UpdateStatusRequest{Status: "COMPLETED"}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPut, "/api/v1/orders/ORD-0001/status", bytes.NewBuffer(body))
		rr := httptest.NewRecorder()
		h.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("expected 200 OK, got %d", rr.Code)
		}

		var updated model.Order
		_ = json.Unmarshal(rr.Body.Bytes(), &updated)
		if updated.Status != "COMPLETED" {
			t.Errorf("expected status COMPLETED, got %s", updated.Status)
		}
	})

	t.Run("Delete Order", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/orders/ORD-0001", nil)
		rr := httptest.NewRecorder()
		h.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("expected 200 OK, got %d", rr.Code)
		}

		// Verify 404 on subsequent get
		getReq := httptest.NewRequest(http.MethodGet, "/api/v1/orders/ORD-0001", nil)
		getRr := httptest.NewRecorder()
		h.ServeHTTP(getRr, getReq)
		if getRr.Code != http.StatusNotFound {
			t.Errorf("expected 404 after deletion, got %d", getRr.Code)
		}
	})
}
