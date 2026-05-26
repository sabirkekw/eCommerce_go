package apierrors

import "errors"

var (
	ErrProductNotFound     = errors.New("Product not found")
	ErrFailedToReadProduct = errors.New("Failed to read product")
	ErrIncorrectID         = errors.New("Incorrect ID")
	ErrUnknown             = errors.New("Unknown error")
	ErrFailedToUpdateProduct = errors.New("Failed to update product")
)
