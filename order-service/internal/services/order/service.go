package orderservice

import (
	"context"
	"errors"

	order "github.com/sabirkekw/ecommerce_go/order-service/internal/models/order"
	productsProto "github.com/sabirkekw/ecommerce_go/pkg/api/products"
	"github.com/sabirkekw/ecommerce_go/pkg/apierrors"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Repository interface {
	CreateOrder(ctx context.Context, order *order.OrderData) (string, error)
	GetOrder(ctx context.Context, id string) (*order.OrderData, error)
	UpdateOrder(ctx context.Context, order *order.OrderData) (*order.OrderData, error)
	DeleteOrder(ctx context.Context, id string) (bool, error)
	ListOrders(ctx context.Context) ([]*order.OrderData, error)
}

type ProductClient interface {
	GetProductByID(ctx context.Context, id string) (*productsProto.Product, error)
}

type Service struct {
	storage        Repository
	productsClient ProductClient
	logger         *zap.SugaredLogger
}

func NewService(storage Repository, client ProductClient, logger *zap.SugaredLogger) *Service {
	return &Service{
		storage:        storage,
		productsClient: client,
		logger:         logger,
	}
}

func (s *Service) CreateOrder(ctx context.Context, order *order.OrderData) (string, error) {
	const op = "Service.CreateOrder"
	s.logger.Debugw("Creating order", "item", order.ItemID, "quantity", order.Quantity, "op", op)

	if order.ItemID == "" || order.Quantity <= 0 {
		s.logger.Debugw("Invalid order data", "item", order.ItemID, "quantity", order.Quantity, "op", op)
		return "", errors.New("invalid order data")
	}
	productData, err := s.productsClient.GetProductByID(ctx, order.ItemID)
	if errors.Is(err, status.Errorf(codes.NotFound, "product not found")) {
		s.logger.Debugw("product not found", "op", op)
		return "", apierrors.ErrProductNotFound
	} else if err != nil {
		s.logger.Debugw("failed to get product", "op", op, "error", err)
		return "", apierrors.ErrUnknown
	}
	if productData.Quantity < order.Quantity {
		return "", apierrors.ErrNotEnoughProduct
	}

	id, err := s.storage.CreateOrder(ctx, order)
	if err != nil {
		s.logger.Debugw("Failed to create order", "error", err, "op", op)
		return "", err
	}

	s.logger.Debugw("Order created", "id", id, "item", order.ItemID, "quantity", order.Quantity, "op", op)
	return id, nil
}

func (s *Service) GetOrder(ctx context.Context, id string) (*order.OrderData, error) {
	const op = "Service.GetOrder"
	s.logger.Debugw("Getting order", "id", id, "op", op)

	if id == "" {
		s.logger.Debugw("Invalid order ID", "id", id, "op", op)
		return nil, apierrors.ErrIncorrectID
	}

	orderData, err := s.storage.GetOrder(ctx, id)
	if errors.Is(err, apierrors.ErrOrderNotFound) {
		s.logger.Debugw("order not found", "id", "op", op)
		return nil, err
	} else if err != nil {
		s.logger.Debugw("Failed to get order", "id", id, "error", err, "op", op)
		return nil, err
	}

	s.logger.Debugw("Order found", "id", id, "item_id", orderData.ItemID, "quantity", orderData.Quantity, "op", op)
	return orderData, nil
}

func (s *Service) UpdateOrder(ctx context.Context, order *order.OrderData) (*order.OrderData, error) {
	const op = "Service.UpdateOrder"
	s.logger.Debugw("Updating order", "id", order.ID, "item_id", order.ItemID, "quantity", order.Quantity, "op", op)

	if order.ID == "" || order.ItemID == "" || order.Quantity <= 0 {
		s.logger.Debugw("Invalid order data", "id", order.ID, "item_id", order.ItemID, "quantity", order.Quantity, "op", op)
		return nil, apierrors.ErrInvalidOrderData
	}

	updatedOrder, err := s.storage.UpdateOrder(ctx, order)
	if errors.Is(err, apierrors.ErrOrderNotFound) {
		s.logger.Debugw("order not found", "op", op)
		return nil, err
	} else if err != nil {
		s.logger.Debugw("Failed to update order", "id", order.ID, "error", err, "op", op)
		return nil, err
	}

	s.logger.Debugw("Order updated", "id", updatedOrder.ID, "item_id", updatedOrder.ItemID, "quantity", updatedOrder.Quantity, "op", op)
	return updatedOrder, nil
}

func (s *Service) DeleteOrder(ctx context.Context, id string) (bool, error) {
	const op = "Service.DeleteOrder"
	s.logger.Debugw("Deleting order", "id", id, "op", op)

	if id == "" {
		s.logger.Debugw("Invalid order ID", "id", id, "op", op)
		return false, apierrors.ErrIncorrectID
	}

	success, err := s.storage.DeleteOrder(ctx, id)
	if err != nil {
		s.logger.Debugw("Failed to delete order", "id", id, "error", err, "op", op)
		return false, err
	}

	s.logger.Debugw("Order deleted", "id", id, "op", op)
	return success, nil
}

func (s *Service) ListOrders(ctx context.Context) ([]*order.OrderData, error) {
	const op = "Service.ListOrders"
	s.logger.Debugw("Listing orders", "op", op)

	orders, err := s.storage.ListOrders(ctx)
	if err != nil {
		s.logger.Debugw("Failed to list orders", "error", err, "op", op)
		return nil, err
	}

	s.logger.Debugw("Orders listed", "op", op)
	return orders, nil
}
