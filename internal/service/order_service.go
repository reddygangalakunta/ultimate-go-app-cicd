package service

import (
	"fmt"
	"sync"
	"time"

	"github.com/example/ultimate-ci-cd-pipeline/internal/model"
)

// OrderService interface defines business logic operations.
type OrderService interface {
	CreateOrder(req model.CreateOrderRequest) (*model.Order, error)
	GetOrder(id string) (*model.Order, error)
	ListOrders() ([]*model.Order, error)
}

type InMemoryOrderService struct {
	mu     sync.RWMutex
	orders map[string]*model.Order
	nextID int
}

func NewInMemoryOrderService() *InMemoryOrderService {
	s := &InMemoryOrderService{
		orders: make(map[string]*model.Order),
		nextID: 1,
	}

	// Seed sample order
	_, _ = s.CreateOrder(model.CreateOrderRequest{
		CustomerName: "Enterprise Client",
		Item:         "Cloud Microservice License",
		Quantity:     10,
		Price:        299.99,
	})

	return s
}

func (s *InMemoryOrderService) CreateOrder(req model.CreateOrderRequest) (*model.Order, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	id := fmt.Sprintf("ORD-%04d", s.nextID)
	s.nextID++

	now := time.Now().UTC()
	order := &model.Order{
		ID:           id,
		CustomerName: req.CustomerName,
		Item:         req.Item,
		Quantity:     req.Quantity,
		Price:        req.Price,
		Status:       "CREATED",
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	s.orders[id] = order
	return order, nil
}

func (s *InMemoryOrderService) GetOrder(id string) (*model.Order, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	order, exists := s.orders[id]
	if !exists {
		return nil, model.ErrOrderNotFound
	}

	return order, nil
}

func (s *InMemoryOrderService) ListOrders() ([]*model.Order, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	list := make([]*model.Order, 0, len(s.orders))
	for _, o := range s.orders {
		list = append(list, o)
	}

	return list, nil
}
