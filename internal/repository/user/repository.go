package user

import (
	db "auth/internal/client/db"
	repo "auth/internal/repository"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgconn"
)

type UserRepository struct {
	db db.Client
}

var (
	ErrNotFound  = errors.New("entity not found")
	ErrDuplicate = errors.New("entity already exists")
)

const (
	tableName       = "users"
	idColumn        = "id"
	columnName      = "name"
	columnEmail     = "email"
	columnPassword  = "password"
	columnRole      = "role"
	columnCreatedAt = "created_at"
	columnUpdatedAt = "updated_at"
)

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

func NewUserRepository(db db.Client) repo.UserRepository {
	return &UserRepository{
		db: db,
	}
}

func translateError(err error) error {

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505": // unique_violation
			return fmt.Errorf("%w: %s", ErrDuplicate, pgErr.ConstraintName)

		}
	}
	return err
}
