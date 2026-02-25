package orderservice

import (
	"context"
	"errors"
	"sync"

	uuid "github.com/google/uuid"
	order "github.com/sabirkekw/ecommerce_go/order-service/internal/models/order"
	"go.uber.org/zap"
)

type Service struct {
	storage map[string]*order.OrderData
	logger  *zap.SugaredLogger
	mx      sync.Mutex
}

func NewService(storage map[string]*order.OrderData, logger *zap.SugaredLogger) *Service {
	return &Service{storage: storage, logger: logger, mx: sync.Mutex{}}
}

func (s *Service) CreateOrder(ctx context.Context, order *order.OrderData) (string, error) {
	const op = "Service.CreateOrder"
	s.logger.Infow("Creating order", "item", order.Item, "quantity", order.Quantity, "op", op)

	if order.Item == "" || order.Quantity <= 0 {
		s.logger.Infow("Invalid order data", "item", order.Item, "quantity", order.Quantity, "op", op)
		return "", errors.New("invalid order data")
	}

	id := uuid.New().String()
	order.ID = id

	s.mx.Lock()
	s.storage[id] = order
	s.mx.Unlock()

	s.logger.Infow("Order created", "id", id, "item", order.Item, "quantity", order.Quantity, "op", op)
	return id, nil
}

func (s *Service) GetOrder(ctx context.Context, id string) (*order.OrderData, error) {
	const op = "Service.GetOrder"
	s.logger.Infow("Getting order", "id", id, "op", op)

	if id == "" {
		s.logger.Infow("Invalid order ID", "id", id, "op", op)
		return nil, errors.New("invalid order ID")
	}

	s.mx.Lock()
	orderData, exists := s.storage[id]
	s.mx.Unlock()

	if !exists {
		return nil, errors.New("order not found")
	}

	s.logger.Infow("Order found", "id", id, "item", orderData.Item, "quantity", orderData.Quantity, "op", op)
	return orderData, nil
}

func (s *Service) UpdateOrder(ctx context.Context, order *order.OrderData) (*order.OrderData, error) {
	const op = "Service.UpdateOrder"
	s.logger.Infow("Updating order", "id", order.ID, "item", order.Item, "quantity", order.Quantity, "op", op)

	if order.ID == "" || order.Item == "" || order.Quantity <= 0 {
		s.logger.Infow("Invalid order data", "id", order.ID, "item", order.Item, "quantity", order.Quantity, "op", op)
		return nil, errors.New("invalid order data")
	}

	s.mx.Lock()

	if _, exists := s.storage[order.ID]; !exists {
		return nil, errors.New("order not found")
	}

	s.storage[order.ID] = order

	s.mx.Unlock()

	s.logger.Infow("Order updated", "id", order.ID, "item", order.Item, "quantity", order.Quantity, "op", op)
	return order, nil
}

func (s *Service) DeleteOrder(ctx context.Context, id string) (bool, error) {
	const op = "Service.DeleteOrder"
	s.logger.Infow("Deleting order", "id", id, "op", op)

	if id == "" {
		s.logger.Infow("Invalid order ID", "id", id, "op", op)
		return false, errors.New("invalid order ID")
	}

	s.mx.Lock()
	if _, exists := s.storage[id]; !exists {
		return false, errors.New("order not found")
	}
	delete(s.storage, id)
	s.mx.Unlock()

	s.logger.Infow("Order deleted", "id", id, "op", op)
	return true, nil
}

func (s *Service) ListOrders(ctx context.Context) ([]*order.OrderData, error) {
	const op = "Service.ListOrders"
	s.logger.Infow("Listing orders", "op", op)

	s.mx.Lock()

	var orders []*order.OrderData
	for _, orderData := range s.storage {
		orders = append(orders, orderData)
	}

	s.mx.Unlock()

	s.logger.Infow("Orders listed", "op", op)
	return orders, nil
}
