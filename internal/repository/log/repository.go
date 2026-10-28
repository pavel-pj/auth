package repository

import (
	"auth/internal/client/db"
	repo "auth/internal/repository"
	model "auth/internal/repository/log/model"
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
)

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

const (
	tableName      = "logs"
	idColumn       = "id"
	columnLevel    = "level"
	columnCreateAt = "created_at"
	columnMessage  = "message"
)

type LogRepository struct {
	db db.Client
}

func NewLogRepository(db db.Client) repo.LogRepository {
	return &LogRepository{
		db: db,
	}
}

func (r *LogRepository) Store(ctx context.Context, log model.Log) error {
	strSql, args, err := psql.Insert(tableName).
		Columns(columnCreateAt, columnLevel, columnMessage).
		Values(time.Now(), log.Level, log.Message).
		ToSql()

	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "logs_repository.Store",
		QueryRow: strSql,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return err
	}

	return nil

}
