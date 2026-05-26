package apierrors

import "errors"

var (
	ErrFailedToCheckout = errors.New("Failed to checkout")
	ErrFailedToGetCart  = errors.New("Failed to get cart")
	ErrEmptyCart        = errors.New("Cart is empty")
)
