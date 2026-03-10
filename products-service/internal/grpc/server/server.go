package grpc

import (
	"context"

	proto "github.com/sabirkekw/ecommerce_go/pkg/api/products"
	"github.com/sabirkekw/ecommerce_go/products-service/internal/models/product"
	"go.uber.org/zap"
	"google.golang.org/grpc"
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
