package apierrors

import "errors"

var (
	ErrNotEnoughProduct = errors.New("failed to create order: product not found")
	ErrOrderNotFound    = errors.New("order not found")
	ErrInvalidOrderData = errors.New("invalid order data")
)
