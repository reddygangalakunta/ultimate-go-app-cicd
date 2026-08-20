package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/example/ultimate-ci-cd-pipeline/internal/model"
	"github.com/example/ultimate-ci-cd-pipeline/internal/service"
)

type OrderHandler struct {
	svc service.OrderService
}

func NewOrderHandler(svc service.OrderService) *OrderHandler {
	return &OrderHandler{svc: svc}
}

// ServeHTTP routes order REST API requests cleanly.
func (h *OrderHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/orders")

	switch {
	case path == "" || path == "/":
		if r.Method == http.MethodGet {
			h.ListOrders(w, r)
			return
		}
		if r.Method == http.MethodPost {
			h.CreateOrder(w, r)
			return
		}
		writeJSON(w, http.StatusMethodNotAllowed, model.ErrorResponse{Code: 405, Message: "Method not allowed"})

	case path == "/export" || path == "/export/":
		if r.Method == http.MethodGet {
			h.ExportOrders(w, r)
			return
		}
		writeJSON(w, http.StatusMethodNotAllowed, model.ErrorResponse{Code: 405, Message: "Method not allowed"})

	case path == "/discount" || path == "/discount/":
		if r.Method == http.MethodPost {
			h.ApplyDiscount(w, r)
			return
		}
		writeJSON(w, http.StatusMethodNotAllowed, model.ErrorResponse{Code: 405, Message: "Method not allowed"})

	default:
		parts := strings.Split(strings.Trim(path, "/"), "/")
		id := parts[0]

		if len(parts) == 1 {
			if r.Method == http.MethodGet {
				h.GetOrder(w, r, id)
				return
			}
			if r.Method == http.MethodDelete {
				h.DeleteOrder(w, r, id)
				return
			}
		}

		if len(parts) == 2 && parts[1] == "status" {
			if r.Method == http.MethodPut || r.Method == http.MethodPatch {
				h.UpdateOrderStatus(w, r, id)
				return
			}
		}

		writeJSON(w, http.StatusNotFound, model.ErrorResponse{Code: 404, Message: "Endpoint not found"})
	}
}

func (h *OrderHandler) Metrics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, model.ErrorResponse{Code: 405, Message: "Method not allowed"})
		return
	}
	metrics := h.svc.GetMetrics()
	writeJSON(w, http.StatusOK, metrics)
}

func (h *OrderHandler) AuditLogs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, model.ErrorResponse{Code: 405, Message: "Method not allowed"})
		return
	}
	logs := h.svc.GetAuditLogs()
	writeJSON(w, http.StatusOK, logs)
}

func (h *OrderHandler) ExportOrders(w http.ResponseWriter, r *http.Request) {
	format := strings.ToLower(r.URL.Query().Get("format"))

	if format == "csv" {
		csvData := h.svc.ExportOrdersCSV()
		w.Header().Set("Content-Type", "text/csv")
		w.Header().Set("Content-Disposition", "attachment; filename=orders_export.csv")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(csvData)
		return
	}

	jsonData, err := h.svc.ExportOrdersJSON()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, model.ErrorResponse{Code: 500, Message: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Disposition", "attachment; filename=orders_export.json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(jsonData)
}

func (h *OrderHandler) ApplyDiscount(w http.ResponseWriter, r *http.Request) {
	var req model.ApplyDiscountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, model.ErrorResponse{Code: 400, Message: "Invalid JSON payload"})
		return
	}

	count, pct, err := h.svc.ApplyDiscount(req.CouponCode)
	if err != nil {
		if err == model.ErrInvalidCoupon {
			writeJSON(w, http.StatusBadRequest, model.ErrorResponse{Code: 400, Message: err.Error()})
			return
		}
		writeJSON(w, http.StatusInternalServerError, model.ErrorResponse{Code: 500, Message: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"message":          "Promo discount applied successfully",
		"coupon_code":      req.CouponCode,
		"discount_percent": pct,
		"affected_orders":  count,
	})
}

func (h *OrderHandler) ListOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := h.svc.ListOrders()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, model.ErrorResponse{Code: 500, Message: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, orders)
}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req model.CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, model.ErrorResponse{Code: 400, Message: "Invalid request payload"})
		return
	}

	order, err := h.svc.CreateOrder(req)
	if err != nil {
		if err == model.ErrInvalidOrder {
			writeJSON(w, http.StatusBadRequest, model.ErrorResponse{Code: 400, Message: err.Error()})
			return
		}
		writeJSON(w, http.StatusInternalServerError, model.ErrorResponse{Code: 500, Message: err.Error()})
		return
	}

	writeJSON(w, http.StatusCreated, order)
}

func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request, id string) {
	order, err := h.svc.GetOrder(id)
	if err != nil {
		if err == model.ErrOrderNotFound {
			writeJSON(w, http.StatusNotFound, model.ErrorResponse{Code: 404, Message: err.Error()})
			return
		}
		writeJSON(w, http.StatusInternalServerError, model.ErrorResponse{Code: 500, Message: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, order)
}

func (h *OrderHandler) UpdateOrderStatus(w http.ResponseWriter, r *http.Request, id string) {
	var req model.UpdateStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, model.ErrorResponse{Code: 400, Message: "Invalid JSON payload"})
		return
	}

	order, err := h.svc.UpdateOrderStatus(id, req.Status)
	if err != nil {
		if err == model.ErrOrderNotFound {
			writeJSON(w, http.StatusNotFound, model.ErrorResponse{Code: 404, Message: err.Error()})
			return
		}
		if err == model.ErrInvalidStatus {
			writeJSON(w, http.StatusBadRequest, model.ErrorResponse{Code: 400, Message: err.Error()})
			return
		}
		writeJSON(w, http.StatusInternalServerError, model.ErrorResponse{Code: 500, Message: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, order)
}

func (h *OrderHandler) DeleteOrder(w http.ResponseWriter, r *http.Request, id string) {
	err := h.svc.DeleteOrder(id)
	if err != nil {
		if err == model.ErrOrderNotFound {
			writeJSON(w, http.StatusNotFound, model.ErrorResponse{Code: 404, Message: err.Error()})
			return
		}
		writeJSON(w, http.StatusInternalServerError, model.ErrorResponse{Code: 500, Message: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"message": "order deleted successfully",
		"id":      id,
	})
}
