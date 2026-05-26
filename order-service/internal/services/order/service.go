package orderservice

import (
	"context"
	"errors"

	order "github.com/sabirkekw/ecommerce_go/order-service/internal/models/order"
	productsProto "github.com/sabirkekw/ecommerce_go/pkg/api/products"
	"github.com/sabirkekw/ecommerce_go/pkg/apierrors"
	"go.uber.org/zap"
)

type Repository interface {
	CreateOrder(ctx context.Context, order *order.Order) (int32, error)
	GetOrderByID(ctx context.Context, orderID int32) (*order.Order, error)
	GetOrdersByUserID(ctx context.Context, userID int32) ([]*order.Order, error)
	DeleteOrder(ctx context.Context, orderID int32) error
}

type ProductClient interface {
	GetProductByID(ctx context.Context, id int32) (*productsProto.Product, error)
}

type Service struct {
	storage        Repository
	productsClient ProductClient // deprecated
	logger         *zap.SugaredLogger
}

func NewService(storage Repository, client ProductClient, logger *zap.SugaredLogger) *Service {
	return &Service{
		storage:        storage,
		productsClient: client,
		logger:         logger,
	}
}

func (s *Service) CreateOrder(ctx context.Context, order *order.Order) (int32, error) {
	const op = "Order.Service.CreateOrder"
	s.logger.Debugw("creating order", "user_id", order.UserID, "op", op)

	orderID, err := s.storage.CreateOrder(ctx, order)
	if errors.Is(err, apierrors.ErrIncorrectID) {
		s.logger.Debugw("incorrect ID", "error", err, "op", op)
		return 0, apierrors.ErrIncorrectID
	} else if err != nil {
		s.logger.Errorw("failed to create order in repository", "error", err, "op", op)
		return 0, err
	}
	s.logger.Debugw("order created successfully", "order_id", orderID, "op", op)
	return orderID, nil
}
func (s *Service) GetOrderByID(ctx context.Context, orderID int32) (*order.Order, error) {
	const op = "Order.Service.GetOrderByID"
	s.logger.Debugw("getting order by id", "order_id", orderID, "op", op)

	orderData, err := s.storage.GetOrderByID(ctx, orderID)
	if errors.Is(err, apierrors.ErrOrderNotFound) {
		s.logger.Debugw("order not found", "order_id", orderID, "op", op)
		return nil, apierrors.ErrOrderNotFound
	} else if err != nil {
		s.logger.Errorw("failed to get order from repository", "error", err, "order_id", orderID, "op", op)
		return nil, err
	}
	return orderData, nil
}
func (s *Service) GetOrderByUserID(ctx context.Context, userID int32) ([]*order.Order, error) {
	const op = "Order.Service.GetOrderByUserID"
	s.logger.Debugw("getting orders by user id", "user_id", userID, "op", op)

	ordersData, err := s.storage.GetOrdersByUserID(ctx, userID)
	if errors.Is(err, apierrors.ErrOrderNotFound) {
		s.logger.Debugw("orders not found", "user_id", userID, "op", op)
		return nil, apierrors.ErrOrderNotFound
	} else if err != nil {
		s.logger.Errorw("failed to get orders from repository", "error", err, "user_id", userID, "op", op)
		return nil, err
	}
	return ordersData, nil
}
func (s *Service) DeleteOrder(ctx context.Context, orderID int32) error {
	const op = "Order.Service.DeleteOrder"
	s.logger.Debugw("deleting order", "order_id", orderID, "op", op)

	err := s.storage.DeleteOrder(ctx, orderID)
	if errors.Is(err, apierrors.ErrOrderNotFound) {
		s.logger.Debugw("order not found", "order_id", orderID, "op", op)
		return apierrors.ErrOrderNotFound
	} else if err != nil {
		s.logger.Errorw("failed to delete order from repository", "error", err, "order_id", orderID, "op", op)
		return err
	}

	s.logger.Debugw("order deleted successfully", "order_id", orderID, "op", op)
	return nil
}
