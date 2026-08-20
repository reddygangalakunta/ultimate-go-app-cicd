package model

import (
	"errors"
	"time"
)

var (
	ErrOrderNotFound = errors.New("order not found")
	ErrInvalidOrder  = errors.New("invalid order payload")
)

// Order represents an order domain object.
type Order struct {
	ID          string    `json:"id"`
	CustomerName string   `json:"customer_name"`
	Item        string    `json:"item"`
	Quantity    int       `json:"quantity"`
	Price       float64   `json:"price"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateOrderRequest payload for creating a new order.
type CreateOrderRequest struct {
	CustomerName string  `json:"customer_name"`
	Item         string  `json:"item"`
	Quantity     int     `json:"quantity"`
	Price        float64 `json:"price"`
}

// Validate checks for required fields.
func (r *CreateOrderRequest) Validate() error {
	if r.CustomerName == "" || r.Item == "" || r.Quantity <= 0 || r.Price <= 0 {
		return ErrInvalidOrder
	}
	return nil
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
