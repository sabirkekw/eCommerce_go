package repository

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/sabirkekw/ecommerce_go/order-service/internal/models/order"
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

func (s *Repository) CreateOrder(ctx context.Context, order *order.OrderData) (string, error) {
	const op = "Order.Repository.CreateOrder"
	s.log.Infow("Creating new order", "item", order.Item, "quantity", order.Quantity, "op", op)

	query := s.builder.Insert("orders").
		Columns("item", "quantity").
		Values(order.Item, order.Quantity).
		Suffix("RETURNING id")

	strSql, args, err := query.ToSql()
	if err != nil {
		s.log.Infow("Failed to build SQL query", "error", err, "op", op)
		return "", err
	}

	var id string
	err = s.db.QueryRowContext(ctx, strSql, args...).Scan(&id)
	if err != nil {
		s.log.Infow("Failed to execute SQL query", "error", err, "op", op)
		return "", err
	}

	s.log.Infow("Order created with ID", "id", id, "op", op)
	return id, nil
}

func (s *Repository) GetOrder(ctx context.Context, id string) (*order.OrderData, error) {
	const op = "Order.Repository.GetOrder"
	s.log.Infow("Getting order by ID", "id", id, "op", op)

	query := s.builder.Select("*").
		From("orders").
		Where(sq.Eq{"id": id})

	strSql, args, err := query.ToSql()
	if err != nil {
		s.log.Infow("Failed to build SQL query", "error", err, "op", op)
		return nil, err
	}

	var orderData order.OrderData
	err = s.db.QueryRowContext(ctx, strSql, args...).Scan(&orderData.ID, &orderData.Item, &orderData.Quantity)
	if err != nil {
		s.log.Infow("Failed to execute SQL query", "error", err, "op", op)
		return nil, err
	}

	s.log.Infow("Order retrieved", "id", id, "item", orderData.Item, "quantity", orderData.Quantity, "op", op)
	return &orderData, nil
}

func (s *Repository) UpdateOrder(ctx context.Context, orderdata *order.OrderData) (*order.OrderData, error) {
	const op = "Order.Repository.UpdateOrder"
	s.log.Infow("Updating order", "id", orderdata.ID, "item", orderdata.Item, "quantity", orderdata.Quantity, "op", op)

	query := s.builder.Update("orders").
		Set("item", orderdata.Item).
		Set("quantity", orderdata.Quantity).
		Where(sq.Eq{"id": orderdata.ID}).
		Suffix("RETURNING id, item, quantity")

	strSql, args, err := query.ToSql()
	if err != nil {
		s.log.Infow("Failed to build SQL query", "error", err, "op", op)
		return nil, err
	}

	var updatedOrder order.OrderData

	err = s.db.QueryRowContext(ctx, strSql, args...).Scan(&updatedOrder.ID, &updatedOrder.Item, &updatedOrder.Quantity)
	if err != nil {
		s.log.Infow("Failed to execute SQL query", "error", err, "op", op)
		return nil, err
	}
	s.log.Infow("Order updated")
	return &updatedOrder, nil
}

func (s *Repository) DeleteOrder(ctx context.Context, id string) (bool, error) {
	const op = "Order.Repository.DeleteOrder"

	s.log.Infow("Deleting order", "id", id, "op", op)

	query := s.builder.Delete("orders").
		Where(sq.Eq{"id": id})

	strSql, args, err := query.ToSql()
	if err != nil {
		s.log.Infow("Failed to build SQL query", "error", err, "op", op)
		return false, err
	}

	var success bool

	result, err := s.db.ExecContext(ctx, strSql, args...)
	if err != nil {
		s.log.Infow("Failed to execute SQL query", "error", err, "op", op)
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		s.log.Infow("Failed to get rows affected", "error", err, "op", op)
		return false, err
	}
	success = rowsAffected > 0

	s.log.Infow("Order deleted", "id", id, "success", success, "op", op)
	return success, nil
}

func (s *Repository) ListOrders(ctx context.Context) ([]*order.OrderData, error) {
	const op = "Order.Repository.ListOrders"
	s.log.Infow("Listing all orders", "op", op)

	query := s.builder.Select("*").From("orders")
	strSql, args, err := query.ToSql()
	if err != nil {
		s.log.Infow("Failed to build SQL query", "error", err, "op", op)
		return nil, err
	}

	rows, err := s.db.QueryContext(ctx, strSql, args...)
	if err != nil {
		s.log.Infow("Failed to execute SQL query", "error", err, "op", op)
		return nil, err
	}
	defer rows.Close()

	var orders []*order.OrderData
	for rows.Next() {
		var orderData order.OrderData
		if err := rows.Scan(&orderData.ID, &orderData.Item, &orderData.Quantity); err != nil {
			s.log.Infow("Failed to scan row", "error", err, "op", op)
			return nil, err
		}
		orders = append(orders, &orderData)
	}

	if err := rows.Err(); err != nil {
		s.log.Infow("Row iteration error", "error", err, "op", op)
		return nil, err
	}

	s.log.Infow("Orders listed", "count", len(orders), "op", op)
	return orders, nil
}
