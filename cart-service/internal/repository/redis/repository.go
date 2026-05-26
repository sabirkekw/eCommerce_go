package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/go-redis/redis/v8"
	models "github.com/sabirkekw/ecommerce_go/cart-service/internal/models/product"
	"github.com/sabirkekw/ecommerce_go/pkg/apierrors"
	"go.uber.org/zap"
)

type Repository struct {
	client *redis.Client
	logger *zap.SugaredLogger
}

func New(client *redis.Client, logger *zap.SugaredLogger) *Repository {
	return &Repository{
		client: client,
		logger: logger,
	}
}

func (r *Repository) InsertIntoCart(ctx context.Context, userID int32, product *models.ProductData) error {
	const op = "Cart.Repository.Redis.InsertIntoCart"
	r.logger.Debugw("Inserting cart product into Redis cart", "op", op)

	productJSON, err := json.Marshal(product)
	if err != nil {
		r.logger.Errorw("failed to marshal product to json", "error", err, "op", op)
		return apierrors.ErrUnknown
	}

	key := fmt.Sprintf("cart:%d", userID)
	field := strconv.Itoa(int(product.ID))
	err = r.client.HSet(ctx, key, field, productJSON).Err()
	if err != nil {
		r.logger.Errorw("failed to append product to Redis", "error", err, "op", op)
		return apierrors.ErrUnknown
	}
	r.logger.Debugw("Successfully inserted product into Redis cart", "op", op)
	return nil
}
func (r *Repository) DeleteFromCart(ctx context.Context, userID int32, productID int32) error {
	const op = "Cart.Repository.Redis.DeleteFromCart"
	r.logger.Debugw("Deleting cart product from Redis cart", "op", op)

	key := fmt.Sprintf("cart:%d", userID)
	field := strconv.Itoa(int(productID))
	err := r.client.HDel(ctx, key, field).Err()
	if err != nil {
		r.logger.Errorw("failed to delete product from Redis", "error", err, "op", op)
		return apierrors.ErrUnknown
	}
	r.logger.Debugw("Successfully deleted product from Redis cart", "op", op)
	return nil
}
func (r *Repository) GetCart(ctx context.Context, userID int32) ([]*models.ProductData, error) {
	const op = "Cart.Repository.Redis.GetCart"
	r.logger.Debugw("Getting cart from Redis", "op", op)

	key := fmt.Sprintf("cart:%d", userID)
	products, err := r.client.HGetAll(ctx, key).Result()
	if err != nil {
		r.logger.Errorw("failed to get cart products from Redis", "error", err, "op", op)
		return nil, apierrors.ErrUnknown
	}

	result := make([]*models.ProductData, 0, len(products))
	for _, v := range products {
		var p models.ProductData
		if err := json.Unmarshal([]byte(v), &p); err != nil {
			r.logger.Errorw("Failed to unmarshal product from JSON", "error", err, "op", op)
			return nil, apierrors.ErrUnknown
		}
		r.logger.Debugw("Added product:", "product", p, "op", op)
		result = append(result, &p)
	}

	if len(result) == 0 {
		r.logger.Debugw("Cart is empty", "op", op)
		return nil, apierrors.ErrEmptyCart
	}
	r.logger.Debugw("Successfully retrieved cart from Redis", "op", op)
	return result, nil
}
func (r *Repository) ClearCart(ctx context.Context, userID int32) error {
	panic("implement me")
}
