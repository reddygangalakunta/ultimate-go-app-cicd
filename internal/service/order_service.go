package service

import (
	"fmt"
	"runtime"
	"sort"
	"sync"
	"time"

	"github.com/example/ultimate-ci-cd-pipeline/internal/model"
)

// OrderService interface defines business logic operations.
type OrderService interface {
	CreateOrder(req model.CreateOrderRequest) (*model.Order, error)
	GetOrder(id string) (*model.Order, error)
	ListOrders() ([]*model.Order, error)
	UpdateOrderStatus(id string, status string) (*model.Order, error)
	DeleteOrder(id string) error
	GetMetrics() *model.SystemMetrics
}

type InMemoryOrderService struct {
	mu        sync.RWMutex
	orders    map[string]*model.Order
	nextID    int
	startTime time.Time
}

func NewInMemoryOrderService() *InMemoryOrderService {
	s := &InMemoryOrderService{
		orders:    make(map[string]*model.Order),
		nextID:    1,
		startTime: time.Now().UTC(),
	}

	// Seed initial enterprise orders
	o1, _ := s.CreateOrder(model.CreateOrderRequest{
		CustomerName: "Acme Enterprise Corp",
		Item:         "Cloud Kubernetes License",
		Quantity:     10,
		Price:        499.00,
	})
	s.UpdateOrderStatus(o1.ID, "COMPLETED")

	o2, _ := s.CreateOrder(model.CreateOrderRequest{
		CustomerName: "Global Tech Systems",
		Item:         "Managed Kafka Cluster",
		Quantity:     2,
		Price:        1250.00,
	})
	s.UpdateOrderStatus(o2.ID, "PROCESSING")

	_, _ = s.CreateOrder(model.CreateOrderRequest{
		CustomerName: "DevOps Solutions LLC",
		Item:         "CI/CD Runner Fleet Sub",
		Quantity:     5,
		Price:        199.50,
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

	// Sort orders by ID descending (newest first)
	sort.Slice(list, func(i, j int) bool {
		return list[i].ID > list[j].ID
	})

	return list, nil
}

func (s *InMemoryOrderService) UpdateOrderStatus(id string, status string) (*model.Order, error) {
	req := model.UpdateStatusRequest{Status: status}
	if err := req.Validate(); err != nil {
		return nil, err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	order, exists := s.orders[id]
	if !exists {
		return nil, model.ErrOrderNotFound
	}

	order.Status = status
	order.UpdatedAt = time.Now().UTC()
	return order, nil
}

func (s *InMemoryOrderService) DeleteOrder(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.orders[id]; !exists {
		return model.ErrOrderNotFound
	}

	delete(s.orders, id)
	return nil
}

func (s *InMemoryOrderService) GetMetrics() *model.SystemMetrics {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	var totalRevenue float64
	for _, o := range s.orders {
		if o.Status != "CANCELLED" {
			totalRevenue += (o.Price * float64(o.Quantity))
		}
	}

	return &model.SystemMetrics{
		GoVersion:     runtime.Version(),
		NumGoroutine:  runtime.NumGoroutine(),
		AllocMemoryMB: float64(memStats.Alloc) / 1024 / 1024,
		SysMemoryMB:   float64(memStats.Sys) / 1024 / 1024,
		NumGC:         memStats.NumGC,
		UptimeSeconds: int64(time.Since(s.startTime).Seconds()),
		TotalOrders:   len(s.orders),
		TotalRevenue:  totalRevenue,
	}
}
