package repository

import (
	"database/sql"

	sq "github.com/Masterminds/squirrel"
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
		db: db,
		logger: logger,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *Repository) ReadProduct(id string) (*product.ProductData, error) {
	panic("repo: implement me")
}

func (r *Repository) ReadManyProducts() ([]*product.ProductData, error) {
	panic("repo: implement me")
}

func (r *Repository) UpdateProduct(id string, product *product.ProductData) (*product.ProductData, error) {
	panic("repo: implement me")
}