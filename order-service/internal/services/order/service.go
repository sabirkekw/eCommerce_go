package orderservice

import (
	"context"
	"errors"

	order "github.com/sabirkekw/ecommerce_go/order-service/internal/models/order"
	"go.uber.org/zap"
)

type Repository interface {
	CreateOrder(ctx context.Context, order *order.OrderData) (string, error)
	GetOrder(ctx context.Context, id string) (*order.OrderData, error)
	UpdateOrder(ctx context.Context, order *order.OrderData) (*order.OrderData, error)
	DeleteOrder(ctx context.Context, id string) (bool, error)
	ListOrders(ctx context.Context) ([]*order.OrderData, error)
}

type Service struct {
	storage Repository
	logger  *zap.SugaredLogger
}

func NewService(storage Repository, logger *zap.SugaredLogger) *Service {
	return &Service{storage: storage, logger: logger}
}

func (s *Service) CreateOrder(ctx context.Context, order *order.OrderData) (string, error) {
	const op = "Service.CreateOrder"
	s.logger.Infow("Creating order", "item", order.Item, "quantity", order.Quantity, "op", op)

	if order.Item == "" || order.Quantity <= 0 {
		s.logger.Infow("Invalid order data", "item", order.Item, "quantity", order.Quantity, "op", op)
		return "", errors.New("invalid order data")
	}

	id, err := s.storage.CreateOrder(ctx, order)
	if err != nil {
		s.logger.Infow("Failed to create order", "error", err, "op", op)
		return "", err
	}

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

	orderData, err := s.storage.GetOrder(ctx, id)
	if err != nil {
		s.logger.Infow("Failed to get order", "id", id, "error", err, "op", op)
		return nil, err
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

	updatedOrder, err := s.storage.UpdateOrder(ctx, order)
	if err != nil {
		s.logger.Infow("Failed to update order", "id", order.ID, "error", err, "op", op)
		return nil, err
	}

	s.logger.Infow("Order updated", "id", updatedOrder.ID, "item", updatedOrder.Item, "quantity", updatedOrder.Quantity, "op", op)
	return updatedOrder, nil
}

func (s *Service) DeleteOrder(ctx context.Context, id string) (bool, error) {
	const op = "Service.DeleteOrder"
	s.logger.Infow("Deleting order", "id", id, "op", op)

	if id == "" {
		s.logger.Infow("Invalid order ID", "id", id, "op", op)
		return false, errors.New("invalid order ID")
	}

	success, err := s.storage.DeleteOrder(ctx, id)
	if err != nil {
		s.logger.Infow("Failed to delete order", "id", id, "error", err, "op", op)
		return false, err
	}

	s.logger.Infow("Order deleted", "id", id, "op", op)
	return success, nil
}

func (s *Service) ListOrders(ctx context.Context) ([]*order.OrderData, error) {
	const op = "Service.ListOrders"
	s.logger.Infow("Listing orders", "op", op)

	orders, err := s.storage.ListOrders(ctx)
	if err != nil {
		s.logger.Infow("Failed to list orders", "error", err, "op", op)
		return nil, err
	}

	s.logger.Infow("Orders listed", "op", op)
	return orders, nil
}
