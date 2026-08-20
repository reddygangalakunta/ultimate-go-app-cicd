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

func TestOrderHandler_Integration(t *testing.T) {
	svc := service.NewInMemoryOrderService()
	h := handler.NewOrderHandler(svc)

	t.Run("List Initial Orders", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/orders", nil)
		rr := httptest.NewRecorder()
		h.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", rr.Code)
		}

		var list []model.Order
		if err := json.Unmarshal(rr.Body.Bytes(), &list); err != nil {
			t.Fatalf("failed unmarshaling list: %v", err)
		}
		if len(list) != 1 {
			t.Errorf("expected 1 order, got %d", len(list))
		}
	})

	t.Run("Create Order Success", func(t *testing.T) {
		payload := model.CreateOrderRequest{
			CustomerName: "Test Customer",
			Item:         "Cloud Container",
			Quantity:     5,
			Price:        150.00,
		}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/orders", bytes.NewBuffer(body))
		rr := httptest.NewRecorder()
		h.ServeHTTP(rr, req)

		if rr.Code != http.StatusCreated {
			t.Fatalf("expected 201 Created, got %d", rr.Code)
		}

		var created model.Order
		_ = json.Unmarshal(rr.Body.Bytes(), &created)
		if created.ID == "" || created.CustomerName != payload.CustomerName {
			t.Errorf("unexpected created order: %+v", created)
		}
	})

	t.Run("Create Order Invalid", func(t *testing.T) {
		payload := model.CreateOrderRequest{}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/orders", bytes.NewBuffer(body))
		rr := httptest.NewRecorder()
		h.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Fatalf("expected 400 Bad Request, got %d", rr.Code)
		}
	})
}
