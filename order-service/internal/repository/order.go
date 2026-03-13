package repository

import (
	"context"
	"database/sql"
	"errors"

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

func (s *Repository) CreateOrder(ctx context.Context, order *order.OrderData) (string, error) {
	const op = "Order.Repository.CreateOrder"
	s.log.Debugw("Creating new order", "item_id", order.ItemID, "quantity", order.Quantity, "op", op)

	query := s.builder.Insert("orders").
		Columns("item_id", "quantity", "user_id").
		Values(order.ItemID, order.Quantity, order.UserID).
		Suffix("RETURNING id")

	strSql, args, err := query.ToSql()
	if err != nil {
		s.log.Debugw("Failed to build SQL query", "error", err, "op", op)
		return "", apierrors.ErrUnknown
	}

	var id string
	err = s.db.QueryRowContext(ctx, strSql, args...).Scan(&id)
	if err != nil {
		s.log.Debugw("Failed to execute SQL query", "error", err, "op", op)
		return "", apierrors.ErrUnknown
	}

	s.log.Debugw("Order created with ID", "id", id, "op", op)
	return id, nil
}

func (s *Repository) GetOrder(ctx context.Context, id string) (*order.OrderData, error) {
	const op = "Order.Repository.GetOrder"
	s.log.Debugw("Getting order by ID", "id", id, "op", op)

	query := s.builder.Select("*").
		From("orders").
		Where(sq.Eq{"id": id})

	strSql, args, err := query.ToSql()
	if err != nil {
		s.log.Debugw("Failed to build SQL query", "error", err, "op", op)
		return nil, apierrors.ErrUnknown
	}

	var orderData order.OrderData
	err = s.db.QueryRowContext(ctx, strSql, args...).Scan(&orderData.ID, &orderData.ItemID, &orderData.Quantity, &orderData.UserID)
	if errors.Is(err, sql.ErrNoRows) {
		s.log.Debugw("Order not found")
		return nil, apierrors.ErrOrderNotFound
	} else if err != nil {
		s.log.Debugw("Failed to execute SQL query", "error", err, "op", op)
		return nil, apierrors.ErrUnknown
	}
	s.log.Debugw("Order retrieved", "id", id, "item_id", orderData.ItemID, "quantity", orderData.Quantity, "op", op)
	return &orderData, nil
}

func (s *Repository) UpdateOrder(ctx context.Context, orderdata *order.OrderData) (*order.OrderData, error) {
	const op = "Order.Repository.UpdateOrder"
	s.log.Debugw("Updating order", "id", orderdata.ID, "item_id", orderdata.ItemID, "quantity", orderdata.Quantity, "op", op)

	query := s.builder.Update("orders").
		Set("item_id", orderdata.ItemID).
		Set("quantity", orderdata.Quantity).
		Where(sq.Eq{"id": orderdata.ID}).
		Suffix("RETURNING id, item_id, quantity")

	strSql, args, err := query.ToSql()
	if err != nil {
		s.log.Debugw("Failed to build SQL query", "error", err, "op", op)
		return nil, apierrors.ErrUnknown
	}

	var updatedOrder order.OrderData

	err = s.db.QueryRowContext(ctx, strSql, args...).Scan(&updatedOrder.ID, &updatedOrder.ItemID, &updatedOrder.Quantity)
	if errors.Is(err, sql.ErrNoRows) {
		s.log.Debugw("Order not found", "error", err, "op", op)
		return nil, apierrors.ErrOrderNotFound
	} else if err != nil {
		s.log.Debugw("Failed to execute SQL query", "error", err, "op", op)
		return nil, apierrors.ErrUnknown
	}
	s.log.Debugw("Order updated")
	return &updatedOrder, nil
}

func (s *Repository) DeleteOrder(ctx context.Context, id string) (bool, error) {
	const op = "Order.Repository.DeleteOrder"

	s.log.Debugw("Deleting order", "id", id, "op", op)

	query := s.builder.Delete("orders").
		Where(sq.Eq{"id": id})

	strSql, args, err := query.ToSql()
	if err != nil {
		s.log.Debugw("Failed to build SQL query", "error", err, "op", op)
		return false, apierrors.ErrUnknown
	}

	var success bool

	result, err := s.db.ExecContext(ctx, strSql, args...)
	if err != nil {
		s.log.Debugw("Failed to execute SQL query", "error", err, "op", op)
		return false, apierrors.ErrUnknown
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		s.log.Debugw("Failed to get rows affected", "error", err, "op", op)
		return false, apierrors.ErrUnknown
	}
	success = rowsAffected > 0

	s.log.Debugw("Order deleted", "id", id, "success", success, "op", op)
	return success, nil
}

func (s *Repository) ListOrders(ctx context.Context) ([]*order.OrderData, error) {
	const op = "Order.Repository.ListOrders"
	s.log.Debugw("Listing all orders", "op", op)

	query := s.builder.Select("*").From("orders")
	strSql, args, err := query.ToSql()
	if err != nil {
		s.log.Debugw("Failed to build SQL query", "error", err, "op", op)
		return nil, apierrors.ErrUnknown
	}

	rows, err := s.db.QueryContext(ctx, strSql, args...)
	if err != nil {
		s.log.Debugw("Failed to execute SQL query", "error", err, "op", op)
		return nil, apierrors.ErrUnknown
	}
	defer rows.Close()

	var orders []*order.OrderData
	for rows.Next() {
		var orderData order.OrderData
		if err := rows.Scan(&orderData.ID, &orderData.ItemID, &orderData.Quantity, &orderData.UserID); err != nil {
			s.log.Debugw("Failed to scan row", "error", err, "op", op)
			return nil, apierrors.ErrUnknown
		}
		orders = append(orders, &orderData)
	}

	if err := rows.Err(); err != nil {
		s.log.Debugw("Row iteration error", "error", err, "op", op)
		return nil, apierrors.ErrUnknown
	}

	s.log.Debugw("Orders listed", "count", len(orders), "op", op)
	return orders, nil
}
