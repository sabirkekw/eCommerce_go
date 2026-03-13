package repository

import (
	"context"
	"database/sql"
	"errors"
	"strconv"

	sq "github.com/Masterminds/squirrel"
	"github.com/sabirkekw/ecommerce_go/pkg/apierrors"
	"github.com/sabirkekw/ecommerce_go/products-service/internal/models/product"
	"go.uber.org/zap"
)

type Repository struct {
	db      *sql.DB
	logger  *zap.SugaredLogger
	builder sq.StatementBuilderType
}

func New(db *sql.DB, logger *zap.SugaredLogger) *Repository {
	return &Repository{
		db:      db,
		logger:  logger,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *Repository) ReadProduct(ctx context.Context, id string) (*product.ProductData, error) {
	const op = "Products.Repository.ReadProduct"
	r.logger.Debugw("reading product from database", "item_id", id, "op", op)

	id_num, err := strconv.Atoi(id)
	if err != nil {
		r.logger.Debugw("incorrect ID format", "id", id, "op", op)
		return nil, apierrors.ErrIncorrectID
	}
	query := r.builder.Select("*").From("products").Where(sq.Eq{"id": id_num})

	strSql, args, err := query.ToSql()
	if err != nil {
		r.logger.Errorw("failed to build sql query", "error", err, "op", op)
		return nil, apierrors.ErrUnknown
	}

	var product product.ProductData
	if err := r.db.QueryRowContext(ctx, strSql, args...).Scan(&product.ID, &product.Item, &product.Quantity); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			r.logger.Debugw("product not found", "op", op)
			return nil, apierrors.ErrProductNotFound
		} else {
			r.logger.Debugw("failed to fetch product")
			return nil, apierrors.ErrFailedToReadProduct
		}
	}
	r.logger.Debugw("found!", "data", product, "op", op)
	return &product, nil
}

func (r *Repository) ReadManyProducts(ctx context.Context) ([]*product.ProductData, error) {
	const op = "Products.Repository.ReadManyProducts"
	r.logger.Debugw("reading all products from database", "op", op)

	query := r.builder.Select("*").From("products")

	strSql, args, err := query.ToSql()
	if err != nil {
		r.logger.Errorw("failed to build sql query", "error", err, "op", op)
		return nil, apierrors.ErrUnknown
	}
	rows, err := r.db.QueryContext(ctx, strSql, args)
	if err != nil {
		r.logger.Errorw("failed to execute sql query", "error", err, "op", op)
		return nil, apierrors.ErrUnknown
	}
	defer rows.Close()
	var products []*product.ProductData
	for rows.Next() {
		var product product.ProductData
		if err := rows.Scan(&product.ID, &product.Item, &product.Quantity); err != nil {
			r.logger.Debugw("failed to read row", "error", err, "op", op)
			return nil, apierrors.ErrFailedToReadProduct
		}
		products = append(products, &product)
	}
	if err := rows.Err(); err != nil {
		r.logger.Debugw("Row iteration error", "error", err, "op", op)
		return nil, apierrors.ErrFailedToReadProduct
	}

	r.logger.Debugw("Products listed", "op", op)
	return products, nil
}

func (r *Repository) UpdateProduct(ctx context.Context, id string, oldProduct *product.ProductData) (*product.ProductData, error) {
	const op = "Products.Repository.ReadManyProducts"
	r.logger.Debugw("reading all products from database", "op", op)

	query := r.builder.Update("products").
		Set("item", oldProduct.Item).
		Set("quantity", oldProduct.Quantity).
		Where(sq.Eq{"id": oldProduct.ID})

	strSql, args, err := query.ToSql()
	if err != nil {
		r.logger.Errorw("failed to build sql query", "error", err, "op", op)
		return nil, apierrors.ErrUnknown
	}

	var newProduct product.ProductData
	if err := r.db.QueryRowContext(ctx, strSql, args).Scan(&newProduct.ID, &newProduct.Item, &newProduct.Quantity); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			r.logger.Debugw("product not found", "op", op)
			return nil, apierrors.ErrProductNotFound
		} else {
			r.logger.Errorw("failed to update row", "op", op)
			return nil, apierrors.ErrFailedToUpdateProduct
		}
	}
	return &newProduct, nil
}
