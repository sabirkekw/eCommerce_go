package service

import (
	"context"
	"errors"

	models "github.com/sabirkekw/ecommerce_go/cart-service/internal/models/product"
	protoProducts "github.com/sabirkekw/ecommerce_go/pkg/api/products"
	"github.com/sabirkekw/ecommerce_go/pkg/apierrors"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type MessageSender interface {
	SendCheckoutMessage(ctx context.Context, userID int32, products []*models.ProductData) error
}

type ProductsProvider interface {
	GetProductByID(ctx context.Context, productID int32) (*protoProducts.Product, error)
}

type Repository interface {
	InsertIntoCart(ctx context.Context, userID int32, product *models.ProductData) error
	DeleteFromCart(ctx context.Context, userID int32, productID int32) error
	GetCart(ctx context.Context, userID int32) ([]*models.ProductData, error)
	ClearCart(ctx context.Context, userID int32) error
}

type Service struct {
	storage          Repository
	cache            Repository
	productsProvider ProductsProvider
	messageSender    MessageSender
	logger           *zap.SugaredLogger
}

func New(storage Repository, cache Repository, productsProvider ProductsProvider, messageSender MessageSender, logger *zap.SugaredLogger) *Service {
	return &Service{
		storage:          storage,
		cache:            cache,
		productsProvider: productsProvider,
		messageSender:    messageSender,
		logger:           logger,
	}
}

func (s *Service) AddToCart(ctx context.Context, userID int32, productID int32, quantity int32) error {
	const op = "Cart.Service.AddToCart"
	s.logger.Debugw("Adding product to cart", "product_id", productID, "op", op)

	providedProduct, err := s.productsProvider.GetProductByID(ctx, productID)
	if errors.Is(err, status.Errorf(codes.NotFound, "product not found")) {
		s.logger.Debugw("Product not found", "op", op)
		return apierrors.ErrProductNotFound
	} else if err != nil {
		s.logger.Errorw("Failed to get product", "error", err, "op", op)
		return apierrors.ErrFailedToReadProduct
	}
	if providedProduct.Quantity <= quantity {
		s.logger.Debugw("Not enought product", "op", op)
		return apierrors.ErrNotEnoughProduct
	}

	product := &models.ProductData{
		ID:          productID,
		ProductName: providedProduct.ProductName,
		Quantity:    quantity,
		Description: providedProduct.Description,
	}
	if err := s.storage.InsertIntoCart(ctx, userID, product); err != nil {
		s.logger.Errorw("Failed to save product into cart: storage", "error", err, "op", op)
		return apierrors.ErrUnknown
	}
	if err := s.cache.InsertIntoCart(ctx, userID, product); err != nil {
		s.logger.Errorw("Failed to save product into cart: cache", "error", err, "op", op)
		return apierrors.ErrUnknown
	}
	s.logger.Debugw("Successfully added product to cart", "op", op)
	return nil
}
func (s *Service) RemoveFromCart(ctx context.Context, userID int32, productID int32) error {
	const op = "Cart.Service.RemoveFromCart"
	s.logger.Debugw("Removing product from cart", "id", productID, "op", op)

	err := s.cache.DeleteFromCart(ctx, userID, productID)
	if errors.Is(err, apierrors.ErrProductNotFound) {
		s.logger.Debugw("No product in cached cart", "op", op)
	} else if err != nil {
		s.logger.Errorw("Failed to remove product from cached cart", "error", err, "op", op)
		return apierrors.ErrUnknown
	}

	err = s.storage.DeleteFromCart(ctx, userID, productID)
	if errors.Is(err, apierrors.ErrProductNotFound) {
		s.logger.Debugw("No product in storage cart", "error", err, "op", op)
		return err
	} else if err != nil {
		s.logger.Errorw("Failed to remove product from cached cart", "error", err, "op", op)
		return apierrors.ErrUnknown
	}
	return nil
}

func (s *Service) GetCart(ctx context.Context, userID int32) ([]*models.ProductData, error) {
	const op = "Cart.Service.GetCart"
	s.logger.Debugw("Getting all products from cart", "User ID", userID, "op", op)

	products, err := s.cache.GetCart(ctx, userID)
	if errors.Is(err, apierrors.ErrEmptyCart) {
		s.logger.Debugw("Cache cart is empty", "op", op)
	} else if err != nil {
		s.logger.Errorw("Failed to get cached cart", "error", err, "op", op)
	} else {
		s.logger.Debugw("Retrieved cart from cache", "op", op)
		return products, nil
	}

	products, err = s.storage.GetCart(ctx, userID)
	if errors.Is(err, apierrors.ErrEmptyCart) {
		s.logger.Errorw("Storage cart is empty", "op", op)
		return nil, err
	} else if err != nil {
		s.logger.Errorw("Failed to get storage cart", "error", err, "op", op)
		return nil, err
	} else {
		s.logger.Debugw("Retrieved cart from storage", "op", op)
		return products, nil
	}
}
func (s *Service) Checkout(ctx context.Context, userID int32) error {
	const op = "Cart.Service.Checkout"
	s.logger.Debugw("Checking out cart", "User ID", userID, "op", op)

	// getting cart products
	products, err := s.GetCart(ctx, userID)
	if err != nil {
		s.logger.Errorw("Failed to get cart for checkout", "error", err, "op", op)
		return err
	}
	// checking if products are still available and have enough quantity
	for _, item := range products {
		p, err := s.productsProvider.GetProductByID(ctx, item.ID)

		if err != nil {
			return apierrors.ErrFailedToReadProduct
		}

		if p.Quantity < item.Quantity {
			return apierrors.ErrNotEnoughProduct
		}
	}

	// sending checkout message
	err = s.messageSender.SendCheckoutMessage(ctx, userID, products)
	if err != nil {
		s.logger.Errorw("Failed to send checkout message", "error", err, "op", op)
		return err
	}

	// clearing cart, TODO: outbox pattern
	err = s.cache.ClearCart(ctx, userID)
	if err != nil {
		s.logger.Errorw("Failed to clear cache cart after checkout", "error", err, "op", op)
		return apierrors.ErrUnknown
	}
	err = s.storage.ClearCart(ctx, userID)
	if err != nil {
		s.logger.Errorw("Failed to clear storage cart after checkout", "error", err, "op", op)
		return apierrors.ErrUnknown
	}

	s.logger.Debugw("Successfully checked out cart", "op", op)
	return nil
}
