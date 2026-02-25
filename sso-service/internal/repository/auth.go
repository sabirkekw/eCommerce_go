package repository

import (
	"context"
	"database/sql"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/sabirkekw/ecommerce_go/sso-service/internal/models"
)

var ErrNoUser = errors.New("user not found")

type UserRepository struct {
	db      *sql.DB
	builder sq.StatementBuilderType
}

func New(db *sql.DB, builder sq.StatementBuilderType) *UserRepository {
	return &UserRepository{
		db:      db,
		builder: builder,
	}
}

func (s *UserRepository) CreateUser(ctx context.Context, firstName string, lastName string, email string, hash []byte) (int64, error) {
	query := s.builder.Insert("users").
		Columns("firstname", "lastname", "email", "hash").
		Values(firstName, lastName, email, hash).
		Suffix("RETURNING id")

	strSql, args, err := query.ToSql()
	if err != nil {
		return 0, err
	}

	var id int64
	err = s.db.QueryRowContext(ctx, strSql, args).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	query := s.builder.Select("*").
		From("users").
		Where(sq.Eq{"email": email})

	strSql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	var user models.User
	row := s.db.QueryRowContext(ctx, strSql, args)
	if err := row.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.PassHash); err != nil {
		return nil, err
	}
	return &user, nil
}
