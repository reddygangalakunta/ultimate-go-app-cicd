package service_test

import (
	"testing"

	"github.com/example/ultimate-ci-cd-pipeline/internal/model"
	"github.com/example/ultimate-ci-cd-pipeline/internal/service"
)

func TestOrderService_CreateAndGet(t *testing.T) {
	svc := service.NewInMemoryOrderService()

	// Verify initial seed order
	orders, err := svc.ListOrders()
	if err != nil {
		t.Fatalf("expected no error listing orders, got %v", err)
	}
	if len(orders) != 1 {
		t.Fatalf("expected 1 initial order, got %d", len(orders))
	}

	// Create new order
	req := model.CreateOrderRequest{
		CustomerName: "Acme Corp",
		Item:         "Kubernetes Cluster",
		Quantity:     2,
		Price:        4999.00,
	}

	created, err := svc.CreateOrder(req)
	if err != nil {
		t.Fatalf("expected order creation to succeed, got %v", err)
	}

	if created.ID == "" {
		t.Errorf("expected non-empty ID")
	}
	if created.CustomerName != req.CustomerName {
		t.Errorf("expected customer %s, got %s", req.CustomerName, created.CustomerName)
	}

	// Retrieve order
	fetched, err := svc.GetOrder(created.ID)
	if err != nil {
		t.Fatalf("expected to get order %s, got error: %v", created.ID, err)
	}
	if fetched.Item != req.Item {
		t.Errorf("expected item %s, got %s", req.Item, fetched.Item)
	}
}

func TestOrderService_CreateInvalid(t *testing.T) {
	svc := service.NewInMemoryOrderService()

	req := model.CreateOrderRequest{
		CustomerName: "",
		Item:         "Test Item",
		Quantity:     -1,
		Price:        0,
	}

	_, err := svc.CreateOrder(req)
	if err != model.ErrInvalidOrder {
		t.Fatalf("expected ErrInvalidOrder, got %v", err)
	}
}

func TestOrderService_GetNotFound(t *testing.T) {
	svc := service.NewInMemoryOrderService()

	_, err := svc.GetOrder("NON_EXISTENT_ID")
	if err != model.ErrOrderNotFound {
		t.Fatalf("expected ErrOrderNotFound, got %v", err)
	}
}
