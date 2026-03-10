package service

import (
	"context"

	"github.com/sabirkekw/ecommerce_go/products-service/internal/models/product"
	"go.uber.org/zap"
)

type Repository interface {
	ReadProduct(id string) (*product.ProductData, error)
	ReadManyProducts() ([]*product.ProductData, error)
	UpdateProduct(id string, product *product.ProductData) (*product.ProductData, error)
}

type Service struct {
	storage Repository
	logger  *zap.SugaredLogger
}

func New(storage Repository, logger *zap.SugaredLogger) *Service {
	return &Service{
		storage: storage,
		logger:  logger,
	}
}

func (s *Service) GetProductById(ctx context.Context, id string) (*product.ProductData, error) {
	panic("service: implement me")
}
func (s *Service) GetProducts(ctx context.Context) ([]*product.ProductData, error) {
	panic("service: implement me")
}
func (s *Service) UpdateProduct(ctx context.Context, id string, product *product.ProductData) (*product.ProductData, error) {
	panic("service: implement me")
}
