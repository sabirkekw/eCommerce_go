package repository

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/sabirkekw/ecommerce_go/order-service/internal/models/order"
	"github.com/sabirkekw/ecommerce_go/pkg/apierrors"
	"go.uber.org/zap"
)

type Repository struct {
	log     *zap.SugaredLogger
	db      *sql.DB
	builder sq.StatementBuilderType
}

func New(db *sql.DB, log *zap.SugaredLogger) *Repository {
	return &Repository{
		log:     log,
		db:      db,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *Repository) CreateOrder(ctx context.Context, order *order.Order) (int32, error) {
	const op = "Order.Repository.CreateOrder"

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		r.log.Errorw("failed to begin transaction", "error", err, "op", op)
		return 0, apierrors.ErrUnknown
	}
	defer tx.Rollback()

	var orderID int32
	err = tx.QueryRowContext(ctx,
		`INSERT INTO orders (user_id, status) VALUES ($1, $2) RETURNING id`,
		order.UserID, order.Status,
	).Scan(&orderID)
	if err != nil {
		r.log.Errorw("failed to insert order", "error", err, "op", op)
		return 0, apierrors.ErrUnknown
	}

	for _, p := range order.Products {
		if _, err := tx.ExecContext(ctx,
			`INSERT INTO order_products (order_id, product_id, quantity) VALUES ($1, $2, $3)`,
			orderID, p.ID, p.Quantity,
		); err != nil {
			r.log.Errorw("failed to insert order product", "error", err, "op", op)
			return 0, apierrors.ErrUnknown
		}
	}

	if err := tx.Commit(); err != nil {
		r.log.Errorw("failed to commit transaction", "error", err, "op", op)
		return 0, apierrors.ErrUnknown
	}

	return orderID, nil
}

func (r *Repository) GetOrderByID(ctx context.Context, orderID int32) (*order.Order, error) {
	const op = "Order.Repository.GetOrderByID"

	query := r.builder.Select(
		"o.id",
		"o.user_id",
		"o.status",
		"oi.product_id",
		"oi.quantity",
	).
		From("orders o").
		LeftJoin("order_items oi ON oi.order_id = o.id").
		Where(sq.Eq{"o.id": orderID})

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return nil, apierrors.ErrUnknown
	}

	rows, err := r.db.QueryContext(ctx, sqlStr, args...)
	if err != nil {
		return nil, apierrors.ErrUnknown
	}
	defer rows.Close()

	var (
		orderData order.Order
		found     bool
	)

	for rows.Next() {
		found = true

		var (
			orderID   int32
			userID    int32
			status    string
			productID *int32
			quantity  *int32
		)

		if err := rows.Scan(&orderID, &userID, &status, &productID, &quantity); err != nil {
			return nil, apierrors.ErrUnknown
		}

		if orderData.ID == 0 {
			orderData.ID = orderID
			orderData.UserID = userID
			orderData.Status = status
		}

		if productID != nil {
			orderData.Products = append(orderData.Products, &order.ProductData{
				ID:       *productID,
				Quantity: *quantity,
			})
		}
	}

	if err := rows.Err(); err != nil {
		return nil, apierrors.ErrUnknown
	}

	if !found {
		return nil, apierrors.ErrOrderNotFound
	}

	return &orderData, nil
}

func (r *Repository) GetOrdersByUserID(ctx context.Context, userID int32) ([]*order.Order, error) {
	query := r.builder.Select(
		"o.id",
		"o.user_id",
		"o.status",
		"oi.product_id",
		"oi.quantity",
	).
		From("orders o").
		LeftJoin("order_items oi ON oi.order_id = o.id").
		Where(sq.Eq{"o.user_id": userID})

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.QueryContext(ctx, sqlStr, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ordersMap := make(map[int32]*order.Order)

	for rows.Next() {
		var (
			orderID   int32
			userID    int32
			status    string
			productID *int32
			quantity  *int32
		)

		if err := rows.Scan(&orderID, &userID, &status, &productID, &quantity); err != nil {
			return nil, err
		}

		o, exists := ordersMap[orderID]
		if !exists {
			o = &order.Order{
				ID:     orderID,
				UserID: userID,
				Status: status,
			}
			ordersMap[orderID] = o
		}

		if productID != nil {
			o.Products = append(o.Products, &order.ProductData{
				ID:       *productID,
				Quantity: *quantity,
			})
		}
	}

	var orders []*order.Order
	for _, o := range ordersMap {
		orders = append(orders, o)
	}

	return orders, nil
}

func (r *Repository) DeleteOrder(ctx context.Context, orderID int32) error {
	const op = "Order.Repository.DeleteOrder"

	query := r.builder.
		Update("orders").
		Set("status", order.StatusCancelled).
		Where(sq.Eq{"id": orderID})

	sqlStr, args, err := query.ToSql()
	if err != nil {
		r.log.Errorw("failed to build query", "error", err, "op", op)
		return apierrors.ErrUnknown
	}

	result, err := r.db.ExecContext(ctx, sqlStr, args...)
	if err != nil {
		r.log.Errorw("failed to execute query", "error", err, "op", op)
		return apierrors.ErrUnknown
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return apierrors.ErrUnknown
	}

	if rowsAffected == 0 {
		return apierrors.ErrOrderNotFound
	}

	return nil
}
