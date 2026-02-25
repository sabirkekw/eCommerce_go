package repository

import (
	"database/sql"

	sq "github.com/Masterminds/squirrel"
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

func (s *Repository) Get() *sql.DB {