package postgres

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	models "github.com/sabirkekw/ecommerce_go/cart-service/internal/models/product"
	"github.com/sabirkekw/ecommerce_go/pkg/apierrors"
	"go.uber.org/zap"
)

type Repository struct {
	db      *sql.DB
	builder sq.StatementBuilderType
	logger  *zap.SugaredLogger
}

func New(db *sql.DB, logger *zap.SugaredLogger) *Repository {
	return &Repository{
		db:      db,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
		logger:  logger,
	}
}

func (r *Repository) InsertIntoCart(ctx context.Context, userID int32, product *models.ProductData) error {
	const op = "Cart.Repository.Postgres.InsertIntoCart"
	r.logger.Debugw("Inserting cart product into database cart", "user_id", userID, "product_id", product.ID, "op", op)

	query := r.builder.Insert("cart").
		Columns("user_id", "product_id", "quantity", "description").
		Values(userID, product.ID, product.Quantity, product.Description).
		Suffix("ON CONFLICT (user_id, product_id) DO UPDATE SET quantity = EXCLUDED.quantity, updated_at = NOW()")

	strSql, args, err := query.ToSql()
	if err != nil {
		r.logger.Errorw("Failed to build SQL query", "error", err)
		return apierrors.ErrUnknown
	}

	if _, err := r.db.ExecContext(ctx, strSql, args...); err != nil {
		r.logger.Errorw("Failed to execute SQL query", "error", err)
		return apierrors.ErrUnknown
	}
	r.logger.Debugw("Successfully inserted product into database cart")
	return nil
}

func (r *Repository) DeleteFromCart(ctx context.Context, userID int32, productID int32) error {
	const op = "Cart.Repository.Postgres.DeleteFromCart"
	r.logger.Debugw("Deleting cart product from database cart", "user_id", userID, "product_id", productID, "op", op)

	query := r.builder.Delete("cart").
		Where(sq.Eq{"user_id": userID, "product_id": productID})

	strSql, args, err := query.ToSql()
	if err != nil {
		r.logger.Errorw("Failed to build SQL query", "error", err)
		return apierrors.ErrUnknown
	}

	if _, err := r.db.ExecContext(ctx, strSql, args...); err != nil {
		r.logger.Errorw("Failed to execute SQL query", "error", err)
		return apierrors.ErrUnknown
	}
	r.logger.Debugw("Successfully deleted product from database cart")
	return nil
}
func (r *Repository) GetCart(ctx context.Context, userID int32) ([]*models.ProductData, error) {
	query := r.builder.Select("user_id", "product_id", "quantity", "description").
		From("cart").
		Where(sq.Eq{"user_id": userID})
	strSql, args, err := query.ToSql()
	if err != nil {
		r.logger.Errorw("Failed to build SQL query", "error", err)
		return nil, apierrors.ErrUnknown
	}

	rows, err := r.db.QueryContext(ctx, strSql, args...)
	if err != nil {
		r.logger.Errorw("Failed to execute SQL query", "error", err)
		return nil, apierrors.ErrUnknown
	}
	defer rows.Close()

	var products []*models.ProductData
	for rows.Next() {
		var product models.ProductData
		if err := rows.Scan(&product.ID, &product.ID, &product.Quantity, &product.Description); err != nil {
			r.logger.Errorw("Failed to scan row", "error", err)
			return nil, apierrors.ErrUnknown
		}
		products = append(products, &product)
	}

	return products, nil
}

func (r *Repository) ClearCart(ctx context.Context, userID int32) error {
	panic("implement me")
}
