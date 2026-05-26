package products

type ProductIDQuantity struct {
	ID       int32 `json:"id"`
	Quantity int32 `json:"quantity"`
}

type CheckoutMessage struct {
	UserID   int32               `json:"user_id"`
	Products []ProductIDQuantity `json:"products"`
}
