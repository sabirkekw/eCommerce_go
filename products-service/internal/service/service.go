package service

import (
	"context"
	"errors"

	"github.com/sabirkekw/ecommerce_go/pkg/apierrors"
	"github.com/sabirkekw/ecommerce_go/products-service/internal/models/product"
	"go.uber.org/zap"
)

type Repository interface {
	ReadProduct(ctx context.Context, id string) (*product.ProductData, error)
	ReadManyProducts(ctx context.Context) ([]*product.ProductData, error)
	UpdateProduct(ctx context.Context, id string, product *product.ProductData) (*product.ProductData, error)
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
	const op = "Products.Service.GetProductByID"
	s.logger.Debugw("getting product by id", "op", op)

	product, err := s.storage.ReadProduct(ctx, id)
	if errors.Is(err, apierrors.ErrIncorrectID) {
		s.logger.Debugw("incorrect ID", "error", err, "op", op)
		return nil, apierrors.ErrIncorrectID
	} else if errors.Is(err, apierrors.ErrProductNotFound) {
		s.logger.Debugw("product not found", "op", op)
		return nil, apierrors.ErrProductNotFound
	} else if err != nil {
		s.logger.Errorw("failed to get product from repository", "error", err, "op", op)
		return nil, err
	}
	return product, nil
}
func (s *Service) GetProducts(ctx context.Context) ([]*product.ProductData, error) {
	const op = "Products.Service.GetProducts"
	s.logger.Debugw("getting all products", "op", op)

	products, err := s.storage.ReadManyProducts(ctx)
	if err != nil {
		s.logger.Debugw("failed to get all products", "error", err, "op", op)
		return nil, err
	}
	return products, nil
}
func (s *Service) UpdateProduct(ctx context.Context, id string, product *product.ProductData) (*product.ProductData, error) {
	const op = "Products.Service.UpdateProduct"
	s.logger.Debugw("updating product", "op", op)

	updatedProduct, err := s.storage.UpdateProduct(ctx, id, product)
	if errors.Is(err, apierrors.ErrFailedToUpdateProduct) {
		s.logger.Debugw("failed to update product", "error", err, "op", op)
		return nil, err
	} else if errors.Is(err, apierrors.ErrProductNotFound) {
		s.logger.Debugw("product not found", "op", op)
		return nil, err
	} else if err != nil {
		s.logger.Debugw("failed to update product", "error", err, "op", op)
		return nil, err
	}
	return updatedProduct, nil
}
