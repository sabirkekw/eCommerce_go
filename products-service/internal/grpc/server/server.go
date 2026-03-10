package grpc

import (
	"context"
	"errors"

	proto "github.com/sabirkekw/ecommerce_go/pkg/api/products"
	"github.com/sabirkekw/ecommerce_go/pkg/apierrors"
	"github.com/sabirkekw/ecommerce_go/products-service/internal/models/product"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ProductsService interface {
	GetProductById(ctx context.Context, id string) (*product.ProductData, error)
	GetProducts(ctx context.Context) ([]*product.ProductData, error)
	UpdateProduct(ctx context.Context, id string, product *product.ProductData) (*product.ProductData, error)
}

type Server struct {
	Service ProductsService
	Logger  *zap.SugaredLogger
	proto.UnimplementedProductsServiceServer
}

func New(service ProductsService, logger *zap.SugaredLogger) *Server {
	return &Server{
		Service: service,
		Logger:  logger,
	}
}

func Register(grpc *grpc.Server, server *Server) {
	proto.RegisterProductsServiceServer(grpc, server)
}

func (s *Server) GetProductByID(ctx context.Context, req *proto.GetProductRequest) (*proto.GetProductResponse, error) {
	id := req.Id
	if id == "" {
		return nil, status.Errorf(codes.InvalidArgument, "incorrect ID")
	}

	product, err := s.Service.GetProductById(ctx, id)
	if errors.Is(err, apierrors.ErrProductNotFound) {
		return nil, status.Errorf(codes.NotFound, "product not found")
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "internal server error")
	}

	return &proto.GetProductResponse{
		Product: &proto.Product{
			Id:       id,
			Name:     product.Item,
			Quantity: product.Quantity,
		},
	}, nil
}
func (s *Server) ListProducts(ctx context.Context, req *proto.ListProductsRequest) (*proto.ListProductsResponse, error) {
	products, err := s.Service.GetProducts(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "internal server error")
	}

	response := make([]*proto.Product, 0, 200)
	for _, product := range products {
		response = append(response, &proto.Product{
			Id:       product.ID,
			Name:     product.Item,
			Quantity: product.Quantity,
		})
	}
	return &proto.ListProductsResponse{
		Products: response,
	}, nil
}
func (s *Server) UpdateProduct(ctx context.Context, req *proto.UpdateProductRequest) (*proto.UpdateProductResponse, error) {
	if req.Product == nil {
		return nil, status.Errorf(codes.InvalidArgument, "nil request")
	}
	id := req.Id

	reqProduct := product.ProductData{
		ID:       id,
		Item:     req.Product.Name,
		Quantity: req.Product.Quantity,
	}
	if id == "" {
		return nil, status.Errorf(codes.InvalidArgument, "incorrect ID")
	}

	respProduct, err := s.Service.UpdateProduct(ctx, id, &reqProduct)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "internal server error")
	}

	return &proto.UpdateProductResponse{
		UpdatedProduct: &proto.Product{
			Id:       respProduct.ID,
			Name:     respProduct.Item,
			Quantity: respProduct.Quantity,
		},
	}, nil
}
