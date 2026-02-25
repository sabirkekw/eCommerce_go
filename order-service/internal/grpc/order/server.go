package grpcserver

import (
	"context"

	"github.com/sabirkekw/ecommerce_go/order-service/internal/models/order"
	proto "github.com/sabirkekw/ecommerce_go/pkg/api/order"
	"go.uber.org/zap"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OrderService interface {
	CreateOrder(ctx context.Context, order *order.OrderData) (string, error)
	GetOrder(ctx context.Context, id string) (*order.OrderData, error)
	UpdateOrder(ctx context.Context, order *order.OrderData) (*order.OrderData, error)
	DeleteOrder(ctx context.Context, id string) (bool, error)
	ListOrders(ctx context.Context) ([]*order.OrderData, error)
}

type Server struct {
	Service OrderService
	Logger  *zap.SugaredLogger
	proto.UnimplementedOrderServiceServer
}

func New(service OrderService, logger *zap.SugaredLogger) *Server {
	return &Server{
		Service: service,
		Logger:  logger,
	}
}

func Register(grpc *grpc.Server, server *Server) {
	proto.RegisterOrderServiceServer(grpc, server)
}

func (s *Server) CreateOrder(ctx context.Context, req *proto.CreateOrderRequest) (*proto.CreateOrderResponse, error) {
	const op = "Server.CreateOrder"
	s.Logger.Infow("Received CreateOrder request", "item", req.GetItem(), "quantity", req.GetQuantity(), "op", op)
	order := &order.OrderData{
		Item:     req.GetItem(),
		Quantity: req.GetQuantity(),
	}
	id, err := s.Service.CreateOrder(ctx, order)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid order data")
	}

	return &proto.CreateOrderResponse{Id: id}, nil
}

func (s *Server) GetOrder(ctx context.Context, req *proto.GetOrderRequest) (*proto.GetOrderResponse, error) {
	const op = "Server.GetOrder"
	s.Logger.Infow("Received GetOrder request", "id", req.GetId(), "op", op)
	order, err := s.Service.GetOrder(ctx, req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid order ID")
	}
	return &proto.GetOrderResponse{
		Order: &proto.Order{
			Id:       req.GetId(),
			Item:     order.Item,
			Quantity: order.Quantity,
		},
	}, nil
}

func (s *Server) UpdateOrder(ctx context.Context, req *proto.UpdateOrderRequest) (*proto.UpdateOrderResponse, error) {
	const op = "Server.UpdateOrder"
	s.Logger.Infow("Received UpdateOrder request", "id", req.GetId(), "item", req.GetItem(), "quantity", req.GetQuantity(), "op", op)
	order := &order.OrderData{
		ID:       req.GetId(),
		Item:     req.GetItem(),
		Quantity: req.GetQuantity(),
	}

	updatedOrder, err := s.Service.UpdateOrder(ctx, order)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid order data")
	}
	return &proto.UpdateOrderResponse{
		Order: &proto.Order{
			Id:       updatedOrder.ID,
			Item:     updatedOrder.Item,
			Quantity: updatedOrder.Quantity,
		},
	}, nil
}

func (s *Server) DeleteOrder(ctx context.Context, req *proto.DeleteOrderRequest) (*proto.DeleteOrderResponse, error) {
	const op = "Server.DeleteOrder"
	s.Logger.Infow("Received DeleteOrder request", "id", req.GetId(), "op", op)
	success, err := s.Service.DeleteOrder(ctx, req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid order ID")
	}
	return &proto.DeleteOrderResponse{Success: success}, nil
}

func (s *Server) ListOrders(ctx context.Context, req *proto.ListOrdersRequest) (*proto.ListOrdersResponse, error) {
	const op = "Server.ListOrders"
	s.Logger.Infow("Received ListOrders request", "op", op)
	orders, err := s.Service.ListOrders(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list orders")
	}
	var protoOrders []*proto.Order
	for _, order := range orders {
		protoOrders = append(protoOrders, &proto.Order{
			Id:       order.ID,
			Item:     order.Item,
			Quantity: order.Quantity,
		})
	}
	return &proto.ListOrdersResponse{Orders: protoOrders}, nil
}
