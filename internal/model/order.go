package model

import (
	"errors"
	"time"
)

var (
	ErrOrderNotFound = errors.New("order not found")
	ErrInvalidOrder  = errors.New("invalid order payload")
	ErrInvalidStatus = errors.New("invalid order status")
)

// Order represents an order domain object.
type Order struct {
	ID           string    `json:"id"`
	CustomerName string    `json:"customer_name"`
	Item         string    `json:"item"`
	Quantity     int       `json:"quantity"`
	Price        float64   `json:"price"`
	Status       string    `json:"status"` // CREATED, PROCESSING, COMPLETED, CANCELLED
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// CreateOrderRequest payload for creating a new order.
type CreateOrderRequest struct {
	CustomerName string  `json:"customer_name"`
	Item         string  `json:"item"`
	Quantity     int     `json:"quantity"`
	Price        float64 `json:"price"`
}

func (r *CreateOrderRequest) Validate() error {
	if r.CustomerName == "" || r.Item == "" || r.Quantity <= 0 || r.Price <= 0 {
		return ErrInvalidOrder
	}
	return nil
}

// UpdateStatusRequest payload for updating order status.
type UpdateStatusRequest struct {
	Status string `json:"status"`
}

func (r *UpdateStatusRequest) Validate() error {
	s := r.Status
	if s != "CREATED" && s != "PROCESSING" && s != "COMPLETED" && s != "CANCELLED" {
		return ErrInvalidStatus
	}
	return nil
}

// SystemMetrics represents Go runtime performance and domain statistics.
type SystemMetrics struct {
	GoVersion     string  `json:"go_version"`
	NumGoroutine  int     `json:"num_goroutine"`
	AllocMemoryMB float64 `json:"alloc_memory_mb"`
	SysMemoryMB   float64 `json:"sys_memory_mb"`
	NumGC         uint32  `json:"num_gc"`
	UptimeSeconds int64   `json:"uptime_seconds"`
	TotalOrders   int     `json:"total_orders"`
	TotalRevenue  float64 `json:"total_revenue"`
}

// HealthResponse represents standard health check output.
type HealthResponse struct {
	Status    string            `json:"status"`
	Version   string            `json:"version"`
	Timestamp time.Time         `json:"timestamp"`
	Details   map[string]string `json:"details,omitempty"`
}

// ErrorResponse represents standard error output.
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
