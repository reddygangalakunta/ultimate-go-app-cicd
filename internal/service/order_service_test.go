package service_test

import (
	"testing"

	"github.com/example/ultimate-ci-cd-pipeline/internal/model"
	"github.com/example/ultimate-ci-cd-pipeline/internal/service"
)

func TestOrderService_CRUDAndMetrics(t *testing.T) {
	svc := service.NewInMemoryOrderService()

	// Initial seeds
	orders, err := svc.ListOrders()
	if err != nil {
		t.Fatalf("expected no error listing orders, got %v", err)
	}
	if len(orders) != 3 {
		t.Fatalf("expected 3 seed orders, got %d", len(orders))
	}

	// Create order
	req := model.CreateOrderRequest{
		CustomerName: "Test Enterprise Client",
		Item:         "Monitoring Dashboard",
		Quantity:     1,
		Price:        1000.00,
	}

	created, err := svc.CreateOrder(req)
	if err != nil {
		t.Fatalf("expected create order to succeed, got %v", err)
	}

	// Update status
	updated, err := svc.UpdateOrderStatus(created.ID, "COMPLETED")
	if err != nil {
		t.Fatalf("expected status update to succeed, got %v", err)
	}
	if updated.Status != "COMPLETED" {
		t.Errorf("expected status COMPLETED, got %s", updated.Status)
	}

	// Metrics
	metrics := svc.GetMetrics()
	if metrics.TotalOrders < 4 {
		t.Errorf("expected at least 4 orders in metrics, got %d", metrics.TotalOrders)
	}
	if metrics.TotalRevenue <= 0 {
		t.Errorf("expected positive total revenue in metrics")
	}

	// Delete order
	if err := svc.DeleteOrder(created.ID); err != nil {
		t.Fatalf("expected delete to succeed, got %v", err)
	}

	// Verify deleted
	_, err = svc.GetOrder(created.ID)
	if err != model.ErrOrderNotFound {
		t.Errorf("expected ErrOrderNotFound after deletion, got %v", err)
	}
}

func TestOrderService_InvalidStatus(t *testing.T) {
	svc := service.NewInMemoryOrderService()

	_, err := svc.UpdateOrderStatus("ORD-0001", "INVALID_STATUS")
	if err != model.ErrInvalidStatus {
		t.Fatalf("expected ErrInvalidStatus, got %v", err)
	}
}
