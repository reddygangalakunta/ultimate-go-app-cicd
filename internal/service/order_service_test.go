package service_test

import (
	"strings"
	"testing"

	"github.com/example/ultimate-ci-cd-pipeline/internal/model"
	"github.com/example/ultimate-ci-cd-pipeline/internal/service"
)

func TestOrderService_AdvancedFeatures(t *testing.T) {
	svc := service.NewInMemoryOrderService()

	// Apply discount coupon
	count, pct, err := svc.ApplyDiscount("ENTERPRISE10")
	if err != nil {
		t.Fatalf("expected discount application to succeed, got %v", err)
	}
	if count <= 0 || pct != 10.0 {
		t.Errorf("unexpected discount results: count=%d, pct=%f", count, pct)
	}

	// Audit logs
	logs := svc.GetAuditLogs()
	if len(logs) == 0 {
		t.Errorf("expected audit logs to be recorded")
	}

	// Export CSV
	csvData := svc.ExportOrdersCSV()
	if !strings.Contains(string(csvData), "Order ID") {
		t.Errorf("expected CSV header in export output")
	}

	// Export JSON
	jsonData, err := svc.ExportOrdersJSON()
	if err != nil || len(jsonData) == 0 {
		t.Fatalf("failed JSON export: %v", err)
	}

	// Invalid coupon code
	_, _, err = svc.ApplyDiscount("INVALID_COUPON")
	if err != model.ErrInvalidCoupon {
		t.Errorf("expected ErrInvalidCoupon, got %v", err)
	}
}
