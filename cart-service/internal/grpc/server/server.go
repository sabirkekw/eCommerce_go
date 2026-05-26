package server

import (
	"context"
	"errors"

	models "github.com/sabirkekw/ecommerce_go/cart-service/internal/models/product"
	proto "github.com/sabirkekw/ecommerce_go/pkg/api/cart"
	"github.com/sabirkekw/ecommerce_go/pkg/apierrors"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CartService interface {
	AddToCart(ctx context.Context, userID int32, productID int32, quantity int32) error
	RemoveFromCart(ctx context.Context, userID int32, productID int32) error
	GetCart(ctx context.Context, userID int32) ([]*models.ProductData, error)
	Checkout(ctx context.Context, userID int32) error
}

type Server struct {
	Service CartService
	Logger  *zap.SugaredLogger
	proto.UnimplementedCartServiceServer
}

func New(service CartService, logger *zap.SugaredLogger) *Server {
	return &Server{
		Service: service,
		Logger:  logger,
	}
}

func Register(grpc *grpc.Server, server *Server) {
	proto.RegisterCartServiceServer(grpc, server)
}

func (s *Server) AddToCart(ctx context.Context, req *proto.AddToCartRequest) (*proto.AddToCartResponse, error) {
	userID, ok := ctx.Value("user_id").(int32)
	if !ok {
		return nil, status.Errorf(codes.Internal, "Failed to fetch user ID")
	}

	productID := req.ProductId
	quantity := req.Quantity
	if productID <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid product ID")
	} else if quantity <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid quantity")
	}

	err := s.Service.AddToCart(ctx, userID, productID, quantity)
	if errors.Is(err, apierrors.ErrProductNotFound) {
		return nil, status.Errorf(codes.NotFound, "product not found")
	} else if errors.Is(err, apierrors.ErrNotEnoughProduct) {
		return nil, status.Errorf(codes.FailedPrecondition, "not enough product")
	} else if errors.Is(err, apierrors.ErrFailedToReadProduct) {
		return nil, status.Errorf(codes.Internal, "failed to read product")
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "internal server error")
	}
	return &proto.AddToCartResponse{
		Success: true,
	}, nil
}

func (s *Server) RemoveFromCart(ctx context.Context, req *proto.RemoveFromCartRequest) (*proto.RemoveFromCartResponse, error) {
	userID, ok := ctx.Value("user_id").(int32)
	if !ok {
		return nil, status.Errorf(codes.Internal, "Failed to fetch user ID")
	}

	productID := req.ProductId
	if productID <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid product ID")
	}
	err := s.Service.RemoveFromCart(ctx, userID, productID)
	if errors.Is(err, apierrors.ErrIncorrectID) {
		return nil, status.Errorf(codes.NotFound, "No product with id %v in cart", productID)
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "internal server error")
	}
	return &proto.RemoveFromCartResponse{
		Success: true,
	}, nil
}

func (s *Server) GetCart(ctx context.Context, req *proto.GetCartRequest) (*proto.GetCartResponse, error) {
	userID, ok := ctx.Value("user_id").(int32)
	if !ok {
		return nil, status.Errorf(codes.Internal, "Failed to fetch user ID")
	}

	products, err := s.Service.GetCart(ctx, userID)
	if errors.Is(err, apierrors.ErrFailedToGetCart) {
		return nil, status.Errorf(codes.Internal, "failed to get cart")
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "internal server error")
	}

	var response []*proto.CartProduct
	for _, product := range products {
		responseProduct := &proto.CartProduct{
			Id:          product.ID,
			ProductName: product.ProductName,
			Quantity:    product.Quantity,
			Description: product.Description,
		}
		response = append(response, responseProduct)
	}
	return &proto.GetCartResponse{
		Cart: &proto.Cart{
			Products: response,
		},
	}, nil
}

func (s *Server) Checkout(ctx context.Context, req *proto.CheckoutRequest) (*proto.CheckoutResponse, error) {
	userID, ok := ctx.Value("user_id").(int32)
	if !ok {
		return nil, status.Errorf(codes.Internal, "Failed to fetch user ID")
	}

	err := s.Service.Checkout(ctx, userID)
	if errors.Is(err, apierrors.ErrFailedToCheckout) {
		return nil, status.Errorf(codes.Internal, "failed to checkout cart")
	} else if errors.Is(err, apierrors.ErrFailedToGetCart) {
		return nil, status.Errorf(codes.Internal, "failed to get cart")
	} else if errors.Is(err, apierrors.ErrEmptyCart) {
		return nil, status.Errorf(codes.FailedPrecondition, "your cart is empty!")
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "internal server error")
	}
	return &proto.CheckoutResponse{
		Success: true,
	}, nil
}
