package grpcserver

import (
	"context"
	"errors"

	"github.com/sabirkekw/ecommerce_go/order-service/internal/models/order"
	proto "github.com/sabirkekw/ecommerce_go/pkg/api/order"
	"github.com/sabirkekw/ecommerce_go/pkg/apierrors"
	"go.uber.org/zap"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)
type OrderService interface {
	CreateOrder(ctx context.Context, order *order.Order) (int32, error)
	GetOrderByID(ctx context.Context, orderID int32) (*order.Order, error)
	GetOrderByUserID(ctx context.Context, userID int32) ([]*order.Order, error)
	DeleteOrder(ctx context.Context, orderID int32) error
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

func (s* Server) GetOrderByID(ctx context.Context, req *proto.GetOrderByIDRequest) (*proto.GetOrderByIDResponse, error) {
	orderID := req.OrderId
	if orderID <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "incorrect ID")
	}

	orderData, err := s.Service.GetOrderByID(ctx, orderID)
	if errors.Is(err, apierrors.ErrOrderNotFound) {
		return nil, status.Errorf(codes.NotFound, "order not found")
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "internal server error")
	}
	return &proto.GetOrderByIDResponse{Order: &proto.SingleOrder{
		Id: orderData.ID,
		UserId: orderData.UserID,
		Products: func() []*proto.ProductData {
			var products []*proto.ProductData
			for _, p := range orderData.Products {
				products = append(products, &proto.ProductData{
					ProductId: p.ID,
					Quantity: p.Quantity,
				})
			}
			return products
		}(),
	}}, nil
}

func (s *Server) GetOrdersByUserID(ctx context.Context, req *proto.GetOrdersByUserIDRequest) (*proto.GetOrdersByUserIDResponse, error) {
	userID, ok := ctx.Value("user_id").(int32)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "user ID not found")
	}
	orders, err := s.Service.GetOrderByUserID(ctx, userID)
	if errors.Is(err, apierrors.ErrOrderNotFound) {
		return nil, status.Errorf(codes.NotFound, "orders not found")
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "internal server error")
	}

	var protoOrders []*proto.SingleOrder
	for _, order := range orders {
		protoOrders = append(protoOrders, &proto.SingleOrder{
			Id: order.ID,
			UserId: order.UserID,
			Products: func() []*proto.ProductData {
				var products []*proto.ProductData
				for _, p := range order.Products {
					products = append(products, &proto.ProductData{
						ProductId: p.ID,
						Quantity: p.Quantity,
					})
				}
				return products
			}(),
		})
	}
	return &proto.GetOrdersByUserIDResponse{Orders: protoOrders}, nil
}

func (s *Server) DeleteOrder(ctx context.Context, req *proto.DeleteOrderRequest) (*proto.DeleteOrderResponse, error) {
	orderID := req.OrderId
	if orderID <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "incorrect ID")
	}
	
	err := s.Service.DeleteOrder(ctx, orderID)
	if errors.Is(err, apierrors.ErrOrderNotFound) {
		return nil, status.Errorf(codes.NotFound, "order not found")
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "internal server error")
	}
	return &proto.DeleteOrderResponse{}, nil
}