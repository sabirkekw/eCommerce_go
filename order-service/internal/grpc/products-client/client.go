package productsclient

import (
	"context"
	"fmt"

	productsProto "github.com/sabirkekw/ecommerce_go/pkg/api/products"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ProductsClient struct {
	Logger *zap.SugaredLogger
	Client productsProto.ProductsServiceClient
}

func New(logger *zap.SugaredLogger, port int) *ProductsClient {
	conn, err := grpc.NewClient(fmt.Sprintf("localhost:%d", port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Errorw("failed to start gRPC products client", "error", err)
		return nil
	}
	client := productsProto.NewProductsServiceClient(conn)
	return &ProductsClient{
		Logger: logger,
		Client: client,
	}
}

func (c *ProductsClient) GetProductByID(ctx context.Context, id string) (*productsProto.Product, error) {
	const op = "Order.ProductsClient.GetProductByID"
	c.Logger.Debugw("requesting product data from Products-service", "id", id, "op", op)

	req := &productsProto.GetProductRequest{
		Id: id,
	}

	resp, err := c.Client.GetProductByID(ctx, req)

	if err != nil {
		return nil, err
	}

	return resp.Product, nil
}
