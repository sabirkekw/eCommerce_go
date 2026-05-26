package order

const (
	StatusCreated   = "created"
	StatusCancelled = "cancelled"
)

type ProductData struct {
	ID       int32
	Quantity int32
}

type Order struct {
	ID       int32
	UserID   int32
	Status   string
	Products []*ProductData
}
