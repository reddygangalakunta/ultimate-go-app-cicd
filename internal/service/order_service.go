package service

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"runtime"
	"sort"
	"strconv"
	"strings"
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
	ApplyDiscount(couponCode string) (int, float64, error)
	GetAuditLogs() []*model.AuditLog
	ExportOrdersCSV() []byte
	ExportOrdersJSON() ([]byte, error)
	GetMetrics() *model.SystemMetrics
}

type InMemoryOrderService struct {
	mu          sync.RWMutex
	orders      map[string]*model.Order
	auditLogs   []*model.AuditLog
	nextID      int
	nextAuditID int
	startTime   time.Time
}

func NewInMemoryOrderService() *InMemoryOrderService {
	s := &InMemoryOrderService{
		orders:      make(map[string]*model.Order),
		auditLogs:   make([]*model.AuditLog, 0),
		nextID:      1,
		nextAuditID: 1,
		startTime:   time.Now().UTC(),
	}

	// Seed initial enterprise orders & audit logs
	o1, _ := s.CreateOrder(model.CreateOrderRequest{
		CustomerName: "Acme Enterprise Corp",
		Item:         "Cloud Kubernetes License",
		Quantity:     10,
		Price:        499.00,
	})
	_, _ = s.UpdateOrderStatus(o1.ID, "COMPLETED")

	o2, _ := s.CreateOrder(model.CreateOrderRequest{
		CustomerName: "Global Tech Systems",
		Item:         "Managed Kafka Cluster",
		Quantity:     2,
		Price:        1250.00,
	})
	_, _ = s.UpdateOrderStatus(o2.ID, "PROCESSING")

	_, _ = s.CreateOrder(model.CreateOrderRequest{
		CustomerName: "DevOps Solutions LLC",
		Item:         "CI/CD Runner Fleet Sub",
		Quantity:     5,
		Price:        199.50,
	})

	return s
}

func (s *InMemoryOrderService) logAudit(action, orderID, details string) {
	audit := &model.AuditLog{
		ID:        fmt.Sprintf("AUD-%05d", s.nextAuditID),
		Action:    action,
		OrderID:   orderID,
		Details:   details,
		Timestamp: time.Now().UTC(),
	}
	s.nextAuditID++
	s.auditLogs = append(s.auditLogs, audit)
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
	s.logAudit("CREATE", id, fmt.Sprintf("Created order for %s (Item: %s, Qty: %d)", req.CustomerName, req.Item, req.Quantity))

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

	oldStatus := order.Status
	order.Status = status
	order.UpdatedAt = time.Now().UTC()

	s.logAudit("UPDATE_STATUS", id, fmt.Sprintf("Status changed from %s to %s", oldStatus, status))

	return order, nil
}

func (s *InMemoryOrderService) DeleteOrder(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	order, exists := s.orders[id]
	if !exists {
		return model.ErrOrderNotFound
	}

	customer := order.CustomerName
	delete(s.orders, id)

	s.logAudit("DELETE", id, fmt.Sprintf("Deleted order %s for customer %s", id, customer))

	return nil
}

func (s *InMemoryOrderService) ApplyDiscount(couponCode string) (int, float64, error) {
	code := strings.ToUpper(strings.TrimSpace(couponCode))
	var discountPercent float64

	switch code {
	case "ENTERPRISE10":
		discountPercent = 10.0
	case "CLOUD20":
		discountPercent = 20.0
	case "DEVOPS30":
		discountPercent = 30.0
	default:
		return 0, 0, model.ErrInvalidCoupon
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	count := 0
	for _, o := range s.orders {
		if o.Status != "CANCELLED" {
			o.Price = o.Price * (1.0 - (discountPercent / 100.0))
			o.UpdatedAt = time.Now().UTC()
			count++
		}
	}

	s.logAudit("DISCOUNT", "", fmt.Sprintf("Applied %.0f%% promo coupon '%s' to %d active orders", discountPercent, code, count))

	return count, discountPercent, nil
}

func (s *InMemoryOrderService) GetAuditLogs() []*model.AuditLog {
	s.mu.RLock()
	defer s.mu.RUnlock()

	logs := make([]*model.AuditLog, len(s.auditLogs))
	copy(logs, s.auditLogs)

	// Sort newest first
	sort.Slice(logs, func(i, j int) bool {
		return logs[i].Timestamp.After(logs[j].Timestamp)
	})

	return logs
}

func (s *InMemoryOrderService) ExportOrdersCSV() []byte {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var b bytes.Buffer
	writer := csv.NewWriter(&b)

	// CSV Header
	_ = writer.Write([]string{"Order ID", "Customer Name", "Item Description", "Quantity", "Unit Price", "Total Amount", "Status", "Created At"})

	for _, o := range s.orders {
		total := o.Price * float64(o.Quantity)
		record := []string{
			o.ID,
			o.CustomerName,
			o.Item,
			strconv.Itoa(o.Quantity),
			fmt.Sprintf("%.2f", o.Price),
			fmt.Sprintf("%.2f", total),
			o.Status,
			o.CreatedAt.Format(time.RFC3339),
		}
		_ = writer.Write(record)
	}

	writer.Flush()
	return b.Bytes()
}

func (s *InMemoryOrderService) ExportOrdersJSON() ([]byte, error) {
	orders, err := s.ListOrders()
	if err != nil {
		return nil, err
	}
	return json.MarshalIndent(orders, "", "  ")
}

func (s *InMemoryOrderService) GetMetrics() *model.SystemMetrics {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	var totalRevenue float64
	custRevenue := make(map[string]float64)
	activeCount := 0

	for _, o := range s.orders {
		if o.Status != "CANCELLED" {
			rev := o.Price * float64(o.Quantity)
			totalRevenue += rev
			custRevenue[o.CustomerName] += rev
			activeCount++
		}
	}

	avgOrderVal := 0.0
	if activeCount > 0 {
		avgOrderVal = totalRevenue / float64(activeCount)
	}

	topCust := "N/A"
	maxRev := 0.0
	for cust, rev := range custRevenue {
		if rev > maxRev {
			maxRev = rev
			topCust = cust
		}
	}

	return &model.SystemMetrics{
		GoVersion:       runtime.Version(),
		NumGoroutine:     runtime.NumGoroutine(),
		AllocMemoryMB:    float64(memStats.Alloc) / 1024 / 1024,
		SysMemoryMB:      float64(memStats.Sys) / 1024 / 1024,
		NumGC:            memStats.NumGC,
		UptimeSeconds:    int64(time.Since(s.startTime).Seconds()),
		TotalOrders:      len(s.orders),
		TotalRevenue:     totalRevenue,
		AverageOrderVal: avgOrderVal,
		TopCustomer:     topCust,
	}
}
