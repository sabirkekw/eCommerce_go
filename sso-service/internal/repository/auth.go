package repository

import (
	"context"
	"database/sql"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/sabirkekw/ecommerce_go/sso-service/internal/models"
	"go.uber.org/zap"
)

var ErrNoUser = errors.New("user not found")

type UserRepository struct {
	logger  *zap.SugaredLogger
	db      *sql.DB
	builder sq.StatementBuilderType
}

func New(logger *zap.SugaredLogger, db *sql.DB, builder sq.StatementBuilderType) *UserRepository {
	return &UserRepository{
		logger:  logger,
		db:      db,
		builder: builder,
	}
}

func (s *UserRepository) CreateUser(ctx context.Context, firstName string, lastName string, email string, hash []byte) (int64, error) {
	const op = "sso.Auth.Repository.CreateUser"
	s.logger.Infow("Creating new user", "email", email, "op", op)

	query := s.builder.Insert("users").
		Columns("firstname", "lastname", "email", "hash").
		Values(firstName, lastName, email, hash).
		Suffix("RETURNING id")

	strSql, args, err := query.ToSql()
	if err != nil {
		s.logger.Infow("Failed to build SQL query", "error", err, "op", op)
		return 0, err
	}

	var id int64
	err = s.db.QueryRowContext(ctx, strSql, args...).Scan(&id)
	if err != nil {
		s.logger.Infow("Failed to execute SQL query", "error", err, "op", op)
		return 0, err
	}
	return id, nil
}

func (s *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	const op = "sso.Auth.Repository.GetByEmail"
	s.logger.Infow("Getting user by email", "email", email, "op", op)

	query := s.builder.Select("*").
		From("users").
		Where(sq.Eq{"email": email})

	strSql, args, err := query.ToSql()
	if err != nil {
		s.logger.Infow("Failed to build SQL query", "error", err)
		return nil, err
	}

	var user models.User
	row := s.db.QueryRowContext(ctx, strSql, args...)
	if err := row.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.PassHash); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.logger.Infow("User not found", "email", email, "op", "sso.Auth.Repository.GetByEmail")
			return nil, ErrNoUser
		}
		return nil, err
	}
	return &user, nil
}
